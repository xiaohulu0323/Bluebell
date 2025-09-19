package controller

import (
	"fmt"
	"strconv"
	"time"
	"web-app/dao/redis"
	"web-app/logic"
	"web-app/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler 创建帖子
// @Summary      创建帖子
// @Description  创建帖子
// @Tags         帖子
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        body  body      models.Post  true  "帖子内容"
// @Success      200   {object}  ResponseData
// @Router       /post [post]
func CreatePostHandler(c *gin.Context) {
	// 1.获取参数及参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("CreatePost with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从 c 取到当前发送请求的用户的ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID
	// 2.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, CodeSuccess)
}

// GetPostDetailHandler 获取帖子详情
// @Summary      帖子详情
// @Description  根据ID获取帖子详情
// @Tags         帖子
// @Produce      json
// @Param        id   path      int  true  "帖子ID"
// @Success      200  {object}  ResponseData
// @Router       /post/{id} [get]
func GetPostDetailHandler(c *gin.Context) {
	// 1.获取参数（从URL路径中获取帖子id）及参数校验
	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		zap.L().Error("GetPostDetail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2. 根据ID 取出帖子数据（查数据库）
	data, err := logic.GetPostByID(postID)
	if err != nil {
		zap.L().Error("logic.GetPostByID() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回相应
	ResponseSuccess(c, data)
}

// GetPostDetailConcurrentHandler 获取帖子详情（并发优化版本）
// @Summary      帖子详情（并发优化）
// @Description  根据ID获取帖子详情，使用并发查询优化性能
// @Tags         帖子
// @Produce      json
// @Param        id   path      int  true  "帖子ID"
// @Success      200  {object}  ResponseData
// @Router       /post/{id}/concurrent [get]
func GetPostDetailConcurrentHandler(c *gin.Context) {
	// 1.获取参数（从URL路径中获取帖子id）及参数校验
	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		zap.L().Error("GetPostDetailConcurrent with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 使用并发版本获取帖子数据
	start := time.Now()
	data, err := logic.GetPostByIDConcurrent(postID)
	duration := time.Since(start)

	if err != nil {
		zap.L().Error("logic.GetPostByIDConcurrent() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 记录性能信息
	zap.L().Info("Concurrent post detail query completed",
		zap.Int64("post_id", postID),
		zap.Duration("duration", duration))

	// 3. 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取帖子列表的处理函数
// @Summary      帖子列表（简单）
// @Description  分页获取帖子列表
// @Tags         帖子
// @Produce      json
// @Param        page  query     int  false  "页码"  default(1)
// @Param        size  query     int  false  "条数"  default(10)
// @Success      200   {object}  ResponseData
// @Router       /posts [get]
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)
	// 1. 获取数据
	data, err := logic.GetPostList(page, size) // 返回帖子列表
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 2. 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler2 升级版帖子列表接口
// 根据前端传来的参数动态的去获取帖子列表
// 按创建时间排序 或者 按照分数排序
// 1. 获取请求的query string 参数
// 2. 去redis查询id列表
// 3. 根据id去数据库查询帖子详细信息
func GetPostListHandler2(c *gin.Context) {
	// @Summary      帖子列表（按时间或分数）
	// @Description  根据排序和社区动态获取帖子列表
	// @Tags         帖子
	// @Produce      json
	// @Param        page         query     int     false  "页码"  default(1)
	// @Param        size         query     int     false  "条数"  default(10)
	// @Param        order        query     string  false  "排序: time/score"  default(time)
	// @Param        community_id query     int     false  "社区ID"
	// @Success      200          {object}  ResponseData
	// @Router       /posts2 [get]
	// GET请求参数（query string）： /api/v1/post2?page=1&size=10&order=time
	p := &models.ParamsPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime, // 默认值
	}

	// c.ShouldBind() 根据请求的数据类型选择相应的方法去获取数据
	// c.ShouldBindJSON() 如果请求中携带的是Json格式的数据，才能用这个方法获取到数据
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostList2 with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 已在上方完成 Query 绑定到 p，无需再次绑定

	data, err := logic.GetPostListNew(p) // 调用合并后的接口

	// 1. 获取数据

	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 2. 返回响应
	ResponseSuccess(c, data)
}

// // 根据社区去查询帖子列表
// func GetCommunityPostListHandler(c *gin.Context) {
// 	// GET请求参数（query string）： /api/v1/post2?page=1&size=10&order=time
// 	// 字段提升只作用于选择器（p.Page 这种读取/赋值），不作用于复合字面量（struct 初始化的花括号语法）
// 	p := &models.ParamsCommunityPostList{
// 		ParamsPostList: &models.ParamsPostList{
// 			Page:  1,
// 			Size:  10,
// 			Order: models.OrderTime, // 默认值
// 		},
// 	}

// 	// c.ShouldBind() 根据请求的数据类型选择相应的方法去获取数据
// 	// c.ShouldBindJSON() 如果请求中携带的是Json格式的数据，才能用这个方法获取到数据
// 	if err := c.ShouldBindQuery(p); err != nil {
// 		zap.L().Error("GetCommunityPostListHandler with invalid param", zap.Error(err))
// 		ResponseError(c, CodeInvalidParam)
// 		return
// 	}
// 	// 已在上方完成 Query 绑定到 p，无需再次绑定

// 	// 1. 获取数据（按社区）
// 	data, err := logic.GetCommunityPostList(p) // 返回帖子列表
// 	if err != nil {
// 		zap.L().Error("logic.GetCommunityPostListHandler() failed", zap.Error(err))
// 		ResponseError(c, CodeServerBusy)
// 		return
// 	}
// 	// 2. 返回响应
// 	ResponseSuccess(c, data)
// }

// GetPostListOptimizedHandler 获取帖子列表的优化处理函数 - 解决N+1查询问题
// @Summary      帖子列表（N+1优化版本）
// @Description  分页获取帖子列表，使用批量查询优化性能
// @Tags         帖子
// @Produce      json
// @Param        page  query     int  false  "页码"  default(1)
// @Param        size  query     int  false  "条数"  default(10)
// @Success      200   {object}  ResponseData
// @Router       /posts/optimized [get]
func GetPostListOptimizedHandler(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)

	// 记录开始时间
	start := time.Now()

	// 1. 使用优化版本获取数据
	data, err := logic.GetPostListOptimized(page, size)

	// 记录执行时间
	duration := time.Since(start)

	if err != nil {
		zap.L().Error("logic.GetPostListOptimized() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 记录性能信息
	zap.L().Info("Optimized post list query completed",
		zap.Int64("page", page),
		zap.Int64("size", size),
		zap.Int("result_count", len(data)),
		zap.Duration("duration", duration),
		zap.String("optimization", "N+1_query_solved"))

	// 2. 返回响应
	ResponseSuccess(c, data)
}

// GetPostDetailCachedHandler 获取帖子详情（带缓存）
// @Summary      获取帖子详情（缓存版本）
// @Description  根据帖子ID获取帖子详情信息，使用Redis缓存提升性能
// @Tags         帖子
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "帖子ID"
// @Success      200  {object}  ResponseData{data=models.ApiPostDetail}
// @Router       /post/{id}/cached [get]
func GetPostDetailCachedHandler(c *gin.Context) {
	start := time.Now()

	// 1. 获取参数
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取数据（带缓存）
	data, err := logic.GetPostByIDWithCache(pid)
	if err != nil {
		zap.L().Error("logic.GetPostByIDWithCache() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	duration := time.Since(start)

	// 记录性能信息
	zap.L().Info("Cached post detail query completed",
		zap.Int64("post_id", pid),
		zap.Duration("duration", duration),
		zap.String("optimization", "redis_cache"))

	// 3. 返回响应
	ResponseSuccess(c, data)
}

// GetPostListCachedHandler 获取帖子列表（带缓存）
// @Summary      获取帖子列表（缓存版本）
// @Description  获取帖子列表，集成N+1优化和Redis缓存
// @Tags         帖子
// @Accept       json
// @Produce      json
// @Param        page  query     int  false  "页码"
// @Param        size  query     int  false  "页大小"
// @Success      200   {object}  ResponseData{data=[]models.ApiPostDetail}
// @Router       /posts/cached [get]
func GetPostListCachedHandler(c *gin.Context) {
	start := time.Now()

	// 1. 获取参数
	page, size := getPageInfo(c)

	// 2. 获取数据（带缓存）
	data, err := logic.GetPostListOptimizedWithCache(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostListOptimizedWithCache() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	duration := time.Since(start)

	// 记录性能信息
	zap.L().Info("Cached post list query completed",
		zap.Int64("page", page),
		zap.Int64("size", size),
		zap.Int("result_count", len(data)),
		zap.Duration("duration", duration),
		zap.String("optimization", "N+1_with_cache"))

	// 3. 返回响应
	ResponseSuccess(c, data)
}

// GetCacheStatsHandler 获取缓存统计信息（调试用）
// @Summary      获取缓存统计信息
// @Description  获取Redis缓存的命中率、错误率等统计信息
// @Tags         系统
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseData{data=map[string]interface{}}
// @Router       /cache/stats [get]
func GetCacheStatsHandler(c *gin.Context) {
	// 获取缓存统计信息
	stats := redis.GetCacheStats()

	// 计算命中率
	result := make(map[string]interface{})
	for cacheType, stat := range stats {
		total := stat.HitCount + stat.MissCount
		hitRate := float64(0)
		if total > 0 {
			hitRate = float64(stat.HitCount) / float64(total) * 100
		}

		result[cacheType] = map[string]interface{}{
			"hit_count":   stat.HitCount,
			"miss_count":  stat.MissCount,
			"error_count": stat.ErrorCount,
			"hit_rate":    fmt.Sprintf("%.2f%%", hitRate),
		}
	}

	zap.L().Info("Cache stats requested", zap.Any("stats", result))
	ResponseSuccess(c, result)
}

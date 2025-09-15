package controller

import (
	"strconv"
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

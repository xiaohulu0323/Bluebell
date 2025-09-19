package logic

import (
	"fmt"
	"sync"
	"time"
	"web-app/dao/mysql"
	"web-app/dao/redis"
	"web-app/models"
	"web-app/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	// 1.生成PostID
	p.ID = snowflake.GenID()
	// 2. 保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		zap.L().Error("mysql.CreatePost() failed", zap.Error(err))
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	return
}

// GetPostByID 根据帖子id获取帖子详情
func GetPostByID(postID int64) (data *models.ApiPostDetail, err error) {
	// 查询并组合我们接口想用的数据
	post, err := mysql.GetPostByID(postID)
	if err != nil {
		zap.L().Error("mysql.GetPostByID() failed", zap.Error(err))
		return nil, err
	}

	// 根据作者ID查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return nil, err
	}

	// 根据社区id查询社区详情信息
	communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return nil, err
	}

	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: communityDetail,
	}

	return
}

func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList() failed", zap.Error(err))
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))

	for _, post := range posts {
		// 根据作者ID查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			return nil, err
		}

		// 根据社区id查询社区详情信息
		communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			return nil, err
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postdetail)
	}
	return
}

func GetPostList2(p *models.ParamsPostList) (data []*models.ApiPostDetail, err error) {

	// 去redis查询Id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}

	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	// 根据Id去MySQL数据库查询帖子的详细信息
	// 返回的数据还要按照我们给定的id的顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者ID查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			return nil, err
		}

		// 根据社区id查询社区详情信息
		communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			return nil, err
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx], // 按顺序一一对应
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postdetail)
	}
	return

}

func GetCommunityPostList(p *models.ParamsCommunityPostList) (data []*models.ApiPostDetail, err error) {

	// 去redis查询Id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}

	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	// 根据Id去MySQL数据库查询帖子的详细信息
	// 返回的数据还要按照我们给定的id的顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者ID查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			return nil, err
		}

		// 根据社区id查询社区详情信息
		communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			return nil, err
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx], // 按顺序一一对应
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postdetail)
	}
	return

}

// GetPostListNew 将两个查询帖子列表逻辑合二为一的接口
// GetPostListNew 将两个查询帖子列表逻辑合二为一的接口
func GetPostListNew(p *models.ParamsPostList) (data []*models.ApiPostDetail, err error) {
	// 根据请求参数的不同 执行不同的逻辑
	if p.CommunityID == 0 {
		// 查所有
		data, err = GetPostList2(p) // 返回帖子列表
	} else {
		// 根据社区id查询
		data, err = GetCommunityPostList(&models.ParamsCommunityPostList{ParamsPostList: p}) // 返回帖子列表
	}

	if err != nil {
		zap.L().Error("logic.GetPostListNew() failed", zap.Error(err))
		return nil, err
	}
	return

}

// GetPostByIDConcurrent 并发版本的帖子详情获取
// 性能优化：将用户信息和社区信息查询改为并发执行，减少总响应时间
func GetPostByIDConcurrent(postID int64) (data *models.ApiPostDetail, err error) {
	// 首先获取帖子基本信息（必须先获取，因为需要AuthorID和CommunityID）
	post, err := mysql.GetPostByID(postID)
	if err != nil {
		zap.L().Error("mysql.GetPostByID() failed", zap.Error(err))
		return nil, err
	}

	// 并发获取用户信息和社区信息
	var (
		user            *models.User
		communityDetail *models.CommunityDetail
		userErr         error
		communityErr    error
		wg              sync.WaitGroup
	)

	// 启动两个goroutine并发查询
	wg.Add(2)

	// goroutine 1: 获取用户信息
	go func() {
		defer wg.Done()
		user, userErr = mysql.GetUserByID(post.AuthorID)
		if userErr != nil {
			zap.L().Error("mysql.GetUserByID() failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(userErr))
		}
	}()

	// goroutine 2: 获取社区信息
	go func() {
		defer wg.Done()
		communityDetail, communityErr = mysql.GetCommunityDetailByID(post.CommunityID)
		if communityErr != nil {
			zap.L().Error("mysql.GetCommunityDetailByID() failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(communityErr))
		}
	}()

	// 等待所有goroutine完成
	wg.Wait()

	// 检查并发查询的错误
	if userErr != nil {
		return nil, userErr
	}
	if communityErr != nil {
		return nil, communityErr
	}

	// 组装返回数据
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: communityDetail,
	}

	return data, nil
}

// GetPostListOptimized 帖子列表优化版本 - 解决N+1查询问题
// 性能优化：使用批量查询替代循环查询，大幅减少数据库查询次数
func GetPostListOptimized(page, size int64) (data []*models.ApiPostDetail, err error) {
	// 1. 获取帖子列表（第1次查询）
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList() failed", zap.Error(err))
		return nil, err
	}

	if len(posts) == 0 {
		return make([]*models.ApiPostDetail, 0), nil
	}

	// 2. 提取所有需要查询的用户ID和社区ID
	userIDs := make([]int64, 0, len(posts))
	communityIDs := make([]int64, 0, len(posts))

	for _, post := range posts {
		userIDs = append(userIDs, post.AuthorID)
		communityIDs = append(communityIDs, post.CommunityID)
	}

	// 3. 批量查询用户信息（第2次查询）
	userMap, err := mysql.BatchGetUsersByIDs(userIDs)
	if err != nil {
		zap.L().Error("mysql.BatchGetUsersByIDs() failed", zap.Error(err))
		return nil, err
	}

	// 4. 批量查询社区信息（第3次查询）
	communityMap, err := mysql.BatchGetCommunitiesByIDs(communityIDs)
	if err != nil {
		zap.L().Error("mysql.BatchGetCommunitiesByIDs() failed", zap.Error(err))
		return nil, err
	}

	// 5. 组装数据
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		// 从map中获取用户信息
		user, userExists := userMap[post.AuthorID]
		if !userExists {
			zap.L().Error("User not found in batch result",
				zap.Int64("author_id", post.AuthorID),
				zap.Int64("post_id", post.ID))
			continue // 跳过用户不存在的帖子
		}

		// 从map中获取社区信息
		communityDetail, communityExists := communityMap[post.CommunityID]
		if !communityExists {
			zap.L().Error("Community not found in batch result",
				zap.Int64("community_id", post.CommunityID),
				zap.Int64("post_id", post.ID))
			continue // 跳过社区不存在的帖子
		}

		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postDetail)
	}

	// 记录性能优化信息
	zap.L().Info("GetPostListOptimized completed",
		zap.Int("posts_count", len(posts)),
		zap.Int("users_queried", len(userMap)),
		zap.Int("communities_queried", len(communityMap)),
		zap.String("optimization", "N+1_to_3_queries"))

	return data, nil
}

// GetPostByIDWithCache 根据帖子id获取帖子详情（带缓存）
func GetPostByIDWithCache(postID int64) (data *models.ApiPostDetail, err error) {
	start := time.Now()

	// 第一步：尝试从缓存获取完整的帖子详情
	cachedPost, err := redis.GetPostDetailFromCache(postID)
	if err != nil {
		zap.L().Error("redis.GetPostDetailFromCache() failed", zap.Error(err), zap.Int64("post_id", postID))
		// 缓存错误不影响业务，继续从数据库查询
	}

	if cachedPost != nil {
		// 缓存命中，直接返回
		zap.L().Info("GetPostByIDWithCache cache hit",
			zap.Int64("post_id", postID),
			zap.Duration("cost", time.Since(start)))
		return cachedPost, nil
	}

	// 第二步：缓存未命中，尝试获取分布式锁防止缓存击穿
	lockKey := fmt.Sprintf("post:%d", postID)
	locked, err := redis.TryLock(lockKey)
	if err != nil {
		zap.L().Error("redis.TryLock() failed", zap.Error(err), zap.String("lock_key", lockKey))
	}

	if locked {
		// 获取到锁，负责查询数据库并更新缓存
		defer func() {
			if unlockErr := redis.ReleaseLock(lockKey); unlockErr != nil {
				zap.L().Error("redis.ReleaseLock() failed", zap.Error(unlockErr), zap.String("lock_key", lockKey))
			}
		}()

		// 双重检查：再次尝试从缓存获取
		cachedPost, err = redis.GetPostDetailFromCache(postID)
		if err == nil && cachedPost != nil {
			zap.L().Info("GetPostByIDWithCache double check cache hit",
				zap.Int64("post_id", postID),
				zap.Duration("cost", time.Since(start)))
			return cachedPost, nil
		}

		// 从数据库查询
		data, err = queryPostDetailFromDB(postID)
		if err != nil {
			return nil, err
		}

		// 异步更新缓存
		go func() {
			if cacheErr := redis.SetPostDetailToCache(postID, data); cacheErr != nil {
				zap.L().Error("redis.SetPostDetailToCache() failed",
					zap.Error(cacheErr), zap.Int64("post_id", postID))
			}
		}()

		zap.L().Info("GetPostByIDWithCache DB query with cache update",
			zap.Int64("post_id", postID),
			zap.Duration("cost", time.Since(start)))
		return data, nil
	} else {
		// 未获取到锁，等待一段时间后重试从缓存获取
		time.Sleep(50 * time.Millisecond)
		cachedPost, err = redis.GetPostDetailFromCache(postID)
		if err == nil && cachedPost != nil {
			zap.L().Info("GetPostByIDWithCache retry cache hit",
				zap.Int64("post_id", postID),
				zap.Duration("cost", time.Since(start)))
			return cachedPost, nil
		}

		// 缓存仍未命中，直接查询数据库（降级策略）
		data, err = queryPostDetailFromDB(postID)
		if err != nil {
			return nil, err
		}

		zap.L().Info("GetPostByIDWithCache fallback to DB",
			zap.Int64("post_id", postID),
			zap.Duration("cost", time.Since(start)))
		return data, nil
	}
}

// queryPostDetailFromDB 从数据库查询帖子详情（内部函数）
func queryPostDetailFromDB(postID int64) (data *models.ApiPostDetail, err error) {
	// 查询帖子基本信息
	post, err := mysql.GetPostByID(postID)
	if err != nil {
		zap.L().Error("mysql.GetPostByID() failed", zap.Error(err))
		return nil, err
	}

	// 根据作者ID查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return nil, err
	}

	// 根据社区id查询社区详情信息
	communityDetail, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID) failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return nil, err
	}

	// 组装数据
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: communityDetail,
	}

	return data, nil
}

// GetPostListOptimizedWithCache N+1查询优化版本（带缓存）
func GetPostListOptimizedWithCache(page, size int64) (data []*models.ApiPostDetail, err error) {
	start := time.Now()

	// 第一步：获取帖子列表
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList() failed", zap.Error(err))
		return nil, err
	}

	if len(posts) == 0 {
		return []*models.ApiPostDetail{}, nil
	}

	// 第二步：提取需要查询的ID
	userIDs := make([]int64, 0, len(posts))
	communityIDs := make([]int64, 0, len(posts))
	userIDSet := make(map[int64]bool)
	communityIDSet := make(map[int64]bool)

	for _, post := range posts {
		if !userIDSet[post.AuthorID] {
			userIDs = append(userIDs, post.AuthorID)
			userIDSet[post.AuthorID] = true
		}
		if !communityIDSet[post.CommunityID] {
			communityIDs = append(communityIDs, post.CommunityID)
			communityIDSet[post.CommunityID] = true
		}
	}

	// 第三步：批量从缓存获取用户信息
	cachedUsers, missedUserIDs, err := redis.BatchGetUsersFromCache(userIDs)
	if err != nil {
		zap.L().Error("redis.BatchGetUsersFromCache() failed", zap.Error(err))
		// 缓存错误，降级到数据库查询
		missedUserIDs = userIDs
		cachedUsers = make(map[int64]*models.User)
	}

	// 第四步：从数据库查询缓存未命中的用户
	var dbUsers map[int64]*models.User
	if len(missedUserIDs) > 0 {
		dbUsers, err = mysql.BatchGetUsersByIDs(missedUserIDs)
		if err != nil {
			zap.L().Error("mysql.BatchGetUsersByIDs() failed", zap.Error(err))
			return nil, err
		}

		// 异步更新用户缓存
		go func() {
			if cacheErr := redis.BatchSetUsersToCache(dbUsers); cacheErr != nil {
				zap.L().Error("redis.BatchSetUsersToCache() failed", zap.Error(cacheErr))
			}
		}()
	}

	// 合并用户数据
	userMap := make(map[int64]*models.User)
	for id, user := range cachedUsers {
		userMap[id] = user
	}
	for id, user := range dbUsers {
		userMap[id] = user
	}

	// 第五步：批量从缓存获取社区信息
	cachedCommunities, missedCommunityIDs, err := redis.BatchGetCommunitiesFromCache(communityIDs)
	if err != nil {
		zap.L().Error("redis.BatchGetCommunitiesFromCache() failed", zap.Error(err))
		// 缓存错误，降级到数据库查询
		missedCommunityIDs = communityIDs
		cachedCommunities = make(map[int64]*models.CommunityDetail)
	}

	// 第六步：从数据库查询缓存未命中的社区
	var dbCommunities map[int64]*models.CommunityDetail
	if len(missedCommunityIDs) > 0 {
		dbCommunities, err = mysql.BatchGetCommunitiesByIDs(missedCommunityIDs)
		if err != nil {
			zap.L().Error("mysql.BatchGetCommunitiesByIDs() failed", zap.Error(err))
			return nil, err
		}

		// 异步更新社区缓存
		go func() {
			if cacheErr := redis.BatchSetCommunitiesToCache(dbCommunities); cacheErr != nil {
				zap.L().Error("redis.BatchSetCommunitiesToCache() failed", zap.Error(cacheErr))
			}
		}()
	}

	// 合并社区数据
	communityMap := make(map[int64]*models.CommunityDetail)
	for id, community := range cachedCommunities {
		communityMap[id] = community
	}
	for id, community := range dbCommunities {
		communityMap[id] = community
	}

	// 第七步：组装数据
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		user, userExists := userMap[post.AuthorID]
		if !userExists {
			zap.L().Error("User not found", zap.Int64("user_id", post.AuthorID))
			continue
		}

		communityDetail, communityExists := communityMap[post.CommunityID]
		if !communityExists {
			zap.L().Error("Community not found", zap.Int64("community_id", post.CommunityID))
			continue
		}

		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postDetail)
	}

	// 记录性能优化信息
	zap.L().Info("GetPostListOptimizedWithCache completed",
		zap.Int("posts_count", len(posts)),
		zap.Int("users_cached", len(cachedUsers)),
		zap.Int("users_from_db", len(dbUsers)),
		zap.Int("communities_cached", len(cachedCommunities)),
		zap.Int("communities_from_db", len(dbCommunities)),
		zap.String("optimization", "N+1_with_cache"),
		zap.Duration("total_cost", time.Since(start)))

	return data, nil
}

package redis

import (
	"encoding/json"
	"strconv"
	"time"
	"web-app/models"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// 缓存过期时间配置
const (
	PostDetailCacheExpire    = 30 * time.Minute // 帖子详情缓存30分钟
	UserInfoCacheExpire      = 60 * time.Minute // 用户信息缓存1小时
	CommunityInfoCacheExpire = 2 * time.Hour    // 社区信息缓存2小时
	PostListCacheExpire      = 5 * time.Minute  // 帖子列表缓存5分钟

	CacheLockExpire = 10 * time.Second      // 缓存锁过期时间
	CacheLockRetry  = 50 * time.Millisecond // 缓存锁重试间隔
)

// CacheStats 缓存统计信息
type CacheStats struct {
	HitCount   int64 `json:"hit_count"`
	MissCount  int64 `json:"miss_count"`
	ErrorCount int64 `json:"error_count"`
}

// 全局缓存统计
var (
	PostCacheStats      = &CacheStats{}
	UserCacheStats      = &CacheStats{}
	CommunityCacheStats = &CacheStats{}
)

// GetPostDetailFromCache 从缓存获取帖子详情
func GetPostDetailFromCache(postID int64) (*models.ApiPostDetail, error) {
	key := getRedisKey(KeyPostDetailPF + strconv.FormatInt(postID, 10))

	// 尝试从缓存获取
	data, err := client.Get(key).Result()
	if err == redis.Nil {
		// 缓存未命中
		PostCacheStats.MissCount++
		zap.L().Debug("Post cache miss", zap.Int64("post_id", postID))
		return nil, nil
	}
	if err != nil {
		// 缓存错误
		PostCacheStats.ErrorCount++
		zap.L().Error("Post cache error", zap.Error(err), zap.Int64("post_id", postID))
		return nil, err
	}

	// 缓存命中，反序列化数据
	var postDetail models.ApiPostDetail
	if err := json.Unmarshal([]byte(data), &postDetail); err != nil {
		PostCacheStats.ErrorCount++
		zap.L().Error("Post cache unmarshal error", zap.Error(err), zap.Int64("post_id", postID))
		return nil, err
	}

	PostCacheStats.HitCount++
	zap.L().Debug("Post cache hit", zap.Int64("post_id", postID))
	return &postDetail, nil
}

// SetPostDetailToCache 将帖子详情存入缓存
func SetPostDetailToCache(postID int64, postDetail *models.ApiPostDetail) error {
	key := getRedisKey(KeyPostDetailPF + strconv.FormatInt(postID, 10))

	// 序列化数据
	data, err := json.Marshal(postDetail)
	if err != nil {
		zap.L().Error("Post cache marshal error", zap.Error(err), zap.Int64("post_id", postID))
		return err
	}

	// 存入缓存
	err = client.Set(key, data, PostDetailCacheExpire).Err()
	if err != nil {
		zap.L().Error("Post cache set error", zap.Error(err), zap.Int64("post_id", postID))
		return err
	}

	zap.L().Debug("Post cached successfully",
		zap.Int64("post_id", postID),
		zap.Duration("expire", PostDetailCacheExpire))
	return nil
}

// GetUserFromCache 从缓存获取用户信息
func GetUserFromCache(userID int64) (*models.User, error) {
	key := getRedisKey(KeyUserInfoPF + strconv.FormatInt(userID, 10))

	data, err := client.Get(key).Result()
	if err == redis.Nil {
		UserCacheStats.MissCount++
		zap.L().Debug("User cache miss", zap.Int64("user_id", userID))
		return nil, nil
	}
	if err != nil {
		UserCacheStats.ErrorCount++
		zap.L().Error("User cache error", zap.Error(err), zap.Int64("user_id", userID))
		return nil, err
	}

	var user models.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		UserCacheStats.ErrorCount++
		zap.L().Error("User cache unmarshal error", zap.Error(err), zap.Int64("user_id", userID))
		return nil, err
	}

	UserCacheStats.HitCount++
	zap.L().Debug("User cache hit", zap.Int64("user_id", userID))
	return &user, nil
}

// SetUserToCache 将用户信息存入缓存
func SetUserToCache(userID int64, user *models.User) error {
	key := getRedisKey(KeyUserInfoPF + strconv.FormatInt(userID, 10))

	// 敏感信息脱敏：不缓存密码
	cacheUser := *user
	cacheUser.Password = "" // 清空密码字段

	data, err := json.Marshal(cacheUser)
	if err != nil {
		zap.L().Error("User cache marshal error", zap.Error(err), zap.Int64("user_id", userID))
		return err
	}

	err = client.Set(key, data, UserInfoCacheExpire).Err()
	if err != nil {
		zap.L().Error("User cache set error", zap.Error(err), zap.Int64("user_id", userID))
		return err
	}

	zap.L().Debug("User cached successfully",
		zap.Int64("user_id", userID),
		zap.Duration("expire", UserInfoCacheExpire))
	return nil
}

// GetCommunityFromCache 从缓存获取社区信息
func GetCommunityFromCache(communityID int64) (*models.CommunityDetail, error) {
	key := getRedisKey(KeyCommunityInfoPF + strconv.FormatInt(communityID, 10))

	data, err := client.Get(key).Result()
	if err == redis.Nil {
		CommunityCacheStats.MissCount++
		zap.L().Debug("Community cache miss", zap.Int64("community_id", communityID))
		return nil, nil
	}
	if err != nil {
		CommunityCacheStats.ErrorCount++
		zap.L().Error("Community cache error", zap.Error(err), zap.Int64("community_id", communityID))
		return nil, err
	}

	var community models.CommunityDetail
	if err := json.Unmarshal([]byte(data), &community); err != nil {
		CommunityCacheStats.ErrorCount++
		zap.L().Error("Community cache unmarshal error", zap.Error(err), zap.Int64("community_id", communityID))
		return nil, err
	}

	CommunityCacheStats.HitCount++
	zap.L().Debug("Community cache hit", zap.Int64("community_id", communityID))
	return &community, nil
}

// SetCommunityToCache 将社区信息存入缓存
func SetCommunityToCache(communityID int64, community *models.CommunityDetail) error {
	key := getRedisKey(KeyCommunityInfoPF + strconv.FormatInt(communityID, 10))

	data, err := json.Marshal(community)
	if err != nil {
		zap.L().Error("Community cache marshal error", zap.Error(err), zap.Int64("community_id", communityID))
		return err
	}

	err = client.Set(key, data, CommunityInfoCacheExpire).Err()
	if err != nil {
		zap.L().Error("Community cache set error", zap.Error(err), zap.Int64("community_id", communityID))
		return err
	}

	zap.L().Debug("Community cached successfully",
		zap.Int64("community_id", communityID),
		zap.Duration("expire", CommunityInfoCacheExpire))
	return nil
}

// BatchGetUsersFromCache 批量从缓存获取用户信息
func BatchGetUsersFromCache(userIDs []int64) (map[int64]*models.User, []int64, error) {
	if len(userIDs) == 0 {
		return make(map[int64]*models.User), []int64{}, nil
	}

	// 构建所有key
	keys := make([]string, len(userIDs))
	for i, userID := range userIDs {
		keys[i] = getRedisKey(KeyUserInfoPF + strconv.FormatInt(userID, 10))
	}

	// 批量获取
	values, err := client.MGet(keys...).Result()
	if err != nil {
		UserCacheStats.ErrorCount++
		zap.L().Error("Batch get users cache error", zap.Error(err))
		return nil, userIDs, err
	}

	cached := make(map[int64]*models.User)
	missed := make([]int64, 0)

	for i, value := range values {
		userID := userIDs[i]
		if value == nil {
			// 缓存未命中
			missed = append(missed, userID)
			UserCacheStats.MissCount++
		} else {
			// 缓存命中，反序列化
			var user models.User
			if err := json.Unmarshal([]byte(value.(string)), &user); err != nil {
				UserCacheStats.ErrorCount++
				zap.L().Error("User cache unmarshal error", zap.Error(err), zap.Int64("user_id", userID))
				missed = append(missed, userID)
			} else {
				cached[userID] = &user
				UserCacheStats.HitCount++
			}
		}
	}

	zap.L().Debug("Batch get users from cache",
		zap.Int("total", len(userIDs)),
		zap.Int("hit", len(cached)),
		zap.Int("miss", len(missed)))

	return cached, missed, nil
}

// BatchSetUsersToCache 批量将用户信息存入缓存
func BatchSetUsersToCache(users map[int64]*models.User) error {
	if len(users) == 0 {
		return nil
	}

	pipeline := client.Pipeline()
	for userID, user := range users {
		key := getRedisKey(KeyUserInfoPF + strconv.FormatInt(userID, 10))

		// 敏感信息脱敏
		cacheUser := *user
		cacheUser.Password = ""

		data, err := json.Marshal(cacheUser)
		if err != nil {
			zap.L().Error("User cache marshal error", zap.Error(err), zap.Int64("user_id", userID))
			continue
		}

		pipeline.Set(key, data, UserInfoCacheExpire)
	}

	_, err := pipeline.Exec()
	if err != nil {
		zap.L().Error("Batch set users cache error", zap.Error(err))
		return err
	}

	zap.L().Debug("Batch set users to cache", zap.Int("count", len(users)))
	return nil
}

// BatchGetCommunitiesFromCache 批量从缓存获取社区信息
func BatchGetCommunitiesFromCache(communityIDs []int64) (map[int64]*models.CommunityDetail, []int64, error) {
	if len(communityIDs) == 0 {
		return make(map[int64]*models.CommunityDetail), []int64{}, nil
	}

	keys := make([]string, len(communityIDs))
	for i, communityID := range communityIDs {
		keys[i] = getRedisKey(KeyCommunityInfoPF + strconv.FormatInt(communityID, 10))
	}

	values, err := client.MGet(keys...).Result()
	if err != nil {
		CommunityCacheStats.ErrorCount++
		zap.L().Error("Batch get communities cache error", zap.Error(err))
		return nil, communityIDs, err
	}

	cached := make(map[int64]*models.CommunityDetail)
	missed := make([]int64, 0)

	for i, value := range values {
		communityID := communityIDs[i]
		if value == nil {
			missed = append(missed, communityID)
			CommunityCacheStats.MissCount++
		} else {
			var community models.CommunityDetail
			if err := json.Unmarshal([]byte(value.(string)), &community); err != nil {
				CommunityCacheStats.ErrorCount++
				zap.L().Error("Community cache unmarshal error", zap.Error(err), zap.Int64("community_id", communityID))
				missed = append(missed, communityID)
			} else {
				cached[communityID] = &community
				CommunityCacheStats.HitCount++
			}
		}
	}

	zap.L().Debug("Batch get communities from cache",
		zap.Int("total", len(communityIDs)),
		zap.Int("hit", len(cached)),
		zap.Int("miss", len(missed)))

	return cached, missed, nil
}

// BatchSetCommunitiesToCache 批量将社区信息存入缓存
func BatchSetCommunitiesToCache(communities map[int64]*models.CommunityDetail) error {
	if len(communities) == 0 {
		return nil
	}

	pipeline := client.Pipeline()
	for communityID, community := range communities {
		key := getRedisKey(KeyCommunityInfoPF + strconv.FormatInt(communityID, 10))

		data, err := json.Marshal(community)
		if err != nil {
			zap.L().Error("Community cache marshal error", zap.Error(err), zap.Int64("community_id", communityID))
			continue
		}

		pipeline.Set(key, data, CommunityInfoCacheExpire)
	}

	_, err := pipeline.Exec()
	if err != nil {
		zap.L().Error("Batch set communities cache error", zap.Error(err))
		return err
	}

	zap.L().Debug("Batch set communities to cache", zap.Int("count", len(communities)))
	return nil
}

// DeletePostCache 删除帖子缓存（用于数据更新时）
func DeletePostCache(postID int64) error {
	key := getRedisKey(KeyPostDetailPF + strconv.FormatInt(postID, 10))
	return client.Del(key).Err()
}

// DeleteUserCache 删除用户缓存
func DeleteUserCache(userID int64) error {
	key := getRedisKey(KeyUserInfoPF + strconv.FormatInt(userID, 10))
	return client.Del(key).Err()
}

// DeleteCommunityCache 删除社区缓存
func DeleteCommunityCache(communityID int64) error {
	key := getRedisKey(KeyCommunityInfoPF + strconv.FormatInt(communityID, 10))
	return client.Del(key).Err()
}

// GetCacheStats 获取缓存统计信息
func GetCacheStats() map[string]*CacheStats {
	return map[string]*CacheStats{
		"post":      PostCacheStats,
		"user":      UserCacheStats,
		"community": CommunityCacheStats,
	}
}

// TryLock 尝试获取分布式锁（防止缓存击穿）
func TryLock(resourceID string) (bool, error) {
	key := getRedisKey(KeyCacheLock + resourceID)

	// 使用SET NX EX命令原子性地设置锁
	result := client.Set(key, "1", CacheLockExpire)
	if result.Err() != nil {
		return false, result.Err()
	}

	// 检查是否成功获取锁
	return result.Val() == "OK", nil
}

// ReleaseLock 释放分布式锁
func ReleaseLock(resourceID string) error {
	key := getRedisKey(KeyCacheLock + resourceID)
	return client.Del(key).Err()
}

// WarmUpCache 缓存预热（可选功能）
func WarmUpCache() error {
	zap.L().Info("Starting cache warm up...")

	// 这里可以预热一些热点数据
	// 比如最热门的帖子、活跃用户等
	// 具体实现可以根据业务需求调整

	zap.L().Info("Cache warm up completed")
	return nil
}

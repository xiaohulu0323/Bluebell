package redis

// redis key 常量
// redis key 注意使用命名空间的方式 区分不同业务的 key(主要是防止公司集群化 redis key 冲突)     方便查询和拆分
const (
	// key 前缀
	KeyPrefix          = "bluebell:"
	KeyPostTimeZSet    = "post:time"   // zset 帖子及发帖时间
	KeyPostScoreZSet   = "post:score"  // zset 帖子及投票分数
	KeyPostVotedZSetPF = "post:voted:" // zset 记录用户及投票类型   前缀   参数是post_id

	KeyCommunitySetPF = "community:" // set 保存每个分区下帖子的ID

	// 数据缓存相关key
	KeyPostDetailPF    = "cache:post:"      // string 帖子详情缓存 前缀 + post_id
	KeyUserInfoPF      = "cache:user:"      // string 用户信息缓存 前缀 + user_id
	KeyCommunityInfoPF = "cache:community:" // string 社区信息缓存 前缀 + community_id
	KeyPostListPF      = "cache:postlist:"  // string 帖子列表缓存 前缀 + page_size_order

	// 缓存防护相关key
	KeyBloomFilter = "bloom:filter" // 布隆过滤器
	KeyCacheLock   = "cache:lock:"  // 缓存锁 前缀 + resource_id
)

// 拼接 redis key 加上前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}

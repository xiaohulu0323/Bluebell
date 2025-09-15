package redis

// redis key 常量
// redis key 注意使用命名空间的方式 区分不同业务的 key(主要是防止公司集群化 redis key 冲突)     方便查询和拆分
const (
	// key 前缀
	KeyPrefix = "bluebell:"
	KeyPostTimeZSet = "post:time" // zset 帖子及发帖时间	
	KeyPostScoreZSet = "post:score" // zset 帖子及投票分数
	KeyPostVotedZSetPF = "post:voted:" // zset 记录用户及投票类型   前缀   参数是post_id 
	
	KeyCommunitySetPF = "community:" // set 保存每个分区下帖子的ID
)



// 拼接 redis key 加上前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}

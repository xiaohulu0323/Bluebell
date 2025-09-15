package redis

import (
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// 本项目使用简化版的投票分数
// 投一票就加432分   86400/200 = 432    200张赞成票可以给你的帖子续一天

/* 投票的几种情况
direction = 1 时 有两种情况：
1. 之前没有投过票，现在投赞成票 +432分      --> 更新分数和投票记录
2. 之前投反对票，现在改投赞成票--> 更新分数和投票记录
direction = 0 时 有两种情况：
1. 之前投过赞成票，现在要取消投票 -432分--> 更新分数和投票记录
2. 之前投过反对票，现在要取消投票 +432分--> 更新分数和投票记录
direction = -1 时 有两种情况：
1. 之前没有投过票，现在投反对票 -432分--> 更新分数和投票记录
2. 之前投过赞成票，现在改投反对票--> 更新分数和投票记录

投票的限制：
每个帖子自发表之日起一个星期之内允许投票，超过一个星期就不允许投票了
1. 到期之后将redis中保存的赞成票和反对票数取出来，更新到MySQL中
2. 到期之后删除那个 KeyPostVotedZSetPF
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每票的分数
)

var (
	ErrorVoteTimeExpire = errors.New("投票时间已过")
	ErrorVoteRepeated    = errors.New("不允许重复投票")
)

func CreatePost(postID, communityID int64) error {
	
	pipeline := client.TxPipeline()      // 使用事务 要么一起成功 要么一起失败
	// 帖子发帖时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score: float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score: float64(time.Now().Unix()),
		Member: postID,
	})
	// 把帖子id加到社区的set中
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(postID)))
	pipeline.SAdd(cKey, postID)
	
	_, err := pipeline.Exec()
	return err
}


func VoteForPost(userID, postID string, value float64) error {
	// 1. 判断投票限制
	// 去redis取帖子发帖时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if time.Now().Unix()-int64(postTime) > oneWeekInSeconds {
		return ErrorVoteTimeExpire
	}

	// 2. 更新帖子分数
	// 先查当前用户给当前帖子的投票记录
	oldValue := client.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()

	// 如果和之前的投票一样，则不需要更新
	if value == oldValue {
		return ErrorVoteRepeated
	}

	// 计算分数变化值
	var diff float64
	if value == 0 {
		// 取消投票
		diff = -oldValue * scorePerVote
	} else if oldValue == 0 {
		// 新投票
		diff = value * scorePerVote
	} else {
		// 改投票
		diff = (value - oldValue) * scorePerVote
	}

	// 3. 使用Pipeline确保原子性操作
	pipeline := client.TxPipeline()

	// 更新帖子分数
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), diff, postID)

	// 4. 记录用户为该帖子投票的数据
	if value == 0 {
		// 取消投票，删除投票记录
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPF+postID), userID)
	} else {
		// 添加或更新投票记录
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value,
			Member: userID,
		})
	}

	// 执行所有操作
	_, err := pipeline.Exec()
	return err
}

package logic

import (
	"strconv"
	"web-app/dao/redis"
	"web-app/models"

	"go.uber.org/zap"
)

// VoteForPost 为帖子投票
func VoteForPost(userID int64, p *models.ParamsVote) error {
	zap.L().Debug("VoteForPost", 
		zap.Int64("userID", userID), 
		zap.String("postID", p.PostID), 
		zap.String("postID", p.PostID), 
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
} 
package logic

import (
"web-app/dao/mysql"
"web-app/models"
"web-app/pkg/snowflake"

"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
// 1.生成PostID
p.ID = snowflake.GenID()
// 2. 保存到数据库
return mysql.CreatePost(p)

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

return data, nil
}

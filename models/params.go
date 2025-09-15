package models

// 定义请求参数的结构体

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// ParamsSignUp 注册请求参数
type ParamsSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamsLogin 登录请求参数
type ParamsLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamsVote 投票参数
type ParamsVote struct {
	PostID    string `json:"post_id" binding:"required"`              // 帖子id（前端以字符串传递）
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1"` // 赞成票(1) 反对票(-1) 取消投票(0)
}

// ParamsPostList 获取帖子列表的query string参数
type ParamsPostList struct {
	CommunityID int64  `json:"community_id" form:"community_id"`  // 可以为空
	Page       int64  `json:"page" form:"page"`
	Size       int64  `json:"size" form:"size"`
	Order      string `json:"order" form:"order"`
}

// ParamsCommunityPostList 按社区获取帖子列表的query string参数
type ParamsCommunityPostList struct {
	*ParamsPostList
	
}

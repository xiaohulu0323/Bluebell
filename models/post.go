package models

import "time"

// 内存对齐概念
type Post struct {
	ID          int64     `db:"post_id" json:"id"`
	AuthorID    int64     `db:"author_id" json:"author_id"`
	CommunityID int64     `db:"community_id" json:"community_id" binding:"required"`
	Status      int32     `db:"status" json:"status"`
	Title       string    `db:"title" json:"title" binding:"required"`
	Content     string    `db:"content" json:"content" binding:"required"`
	CreateTime  time.Time `db:"create_time" json:"create_time"`
}

// ApiPostDetail 帖子详情接口结构体
type ApiPostDetail struct {
	AuthorName       string              `json:"author_name"` // 作者用户名
	*Post                                // 嵌入帖子结构体
	*CommunityDetail `json:"community"` // 嵌入社区信息
}

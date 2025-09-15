package mysql

import (
	"database/sql"
	"strings"
	"web-app/models"

	"github.com/jmoiron/sqlx"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(
post_id, title, content, author_id, community_id)
value(?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)

	return
}

// GetPostByID 根据帖子id获取单个帖子详情
func GetPostByID(postID int64) (post *models.Post, err error) {
	sqlStr := `select
post_id, title, content, author_id, community_id, status, create_time
from post
where post_id = ?`
	post = new(models.Post)
	err = db.Get(post, sqlStr, postID)
	if err == sql.ErrNoRows {
		return nil, ErrorInvalidID
	}
	if err != nil {
		return nil, err
	}
	return post, nil
}

func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select
post_id, title, content, author_id, community_id, status, create_time
from post
order by create_time desc
limit ?, ?`
	posts = make([]*models.Post, 0, 2) // 预先分配好容量，避免多次切片扩容 不要写成make([]*models.Post, 2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size) // page=1,size=10 => limit 0,10

	return
}

// GetPostListByIDs根据给定的id列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where post_id in (?)
	order by FIND_IN_SET(post_id, ?)
	`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	
	err = db.Select(&postList, query, args...)  // !!!

	return

}
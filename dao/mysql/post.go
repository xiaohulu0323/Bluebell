package mysql

import (
"database/sql"
"web-app/models"
)

func CreatePost(p *models.Post) (err error){
sqlStr := `insert into post(
post_id, title, content, author_id, community_id)
value(?, ?, ?, ?, ?)`
_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)

return
}

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

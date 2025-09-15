package logic

import (
	"web-app/dao/mysql"
	"web-app/dao/redis"
	"web-app/models"
	"web-app/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	// 1.生成PostID
	p.ID = snowflake.GenID()
	// 2. 保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		zap.L().Error("mysql.CreatePost() failed", zap.Error(err))
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	return
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

	return
}

func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList() failed", zap.Error(err))
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))

	for _, post := range posts {
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
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postdetail)
	}
	return
}

func GetPostList2(p *models.ParamsPostList) (data []*models.ApiPostDetail, err error) {

	// 去redis查询Id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}

	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	// 根据Id去MySQL数据库查询帖子的详细信息
	// 返回的数据还要按照我们给定的id的顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
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
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx], // 按顺序一一对应
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postdetail)
	}
	return

}

func GetCommunityPostList(p *models.ParamsCommunityPostList) (data []*models.ApiPostDetail, err error) {

	// 去redis查询Id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}

	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	// 根据Id去MySQL数据库查询帖子的详细信息
	// 返回的数据还要按照我们给定的id的顺序返回
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
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
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx], // 按顺序一一对应
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postdetail)
	}
	return

}

// GetPostListNew 将两个查询帖子列表逻辑合二为一的接口
// GetPostListNew 将两个查询帖子列表逻辑合二为一的接口
func GetPostListNew(p *models.ParamsPostList) (data []*models.ApiPostDetail, err error) {
	// 根据请求参数的不同 执行不同的逻辑
	if p.CommunityID == 0 {
		// 查所有
		data, err = GetPostList2(p) // 返回帖子列表
	} else {
		// 根据社区id查询
	data, err = GetCommunityPostList(&models.ParamsCommunityPostList{ParamsPostList: p}) // 返回帖子列表
	}

	if err != nil {
		zap.L().Error("logic.GetPostListNew() failed", zap.Error(err))
		return nil, err
	}
	return

}

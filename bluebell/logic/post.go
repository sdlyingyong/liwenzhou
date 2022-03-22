package logic

import (
	"lwz/bluebell/dao/mysql"
	"lwz/bluebell/models"
	"lwz/bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(post *models.Post) (err error) {
	//post id 生成器
	post.ID, err = snowflake.GenId()
	if err != nil {
		return
	}
	//数据库存储器
	//结果返回
	return mysql.CreatePost(post)
}

func GetPostById(pid int64) (date *models.ApiPostDetail, err error) {
	//查询并组合接口需要的数据
	//根据id获取帖子详情
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById(pid) failed",
			zap.Int64("pid", pid),
			zap.Error(err),
		)
		return
	}
	//根据作者id查询作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
			zap.Int64("post.AuthorID", post.AuthorID),
			zap.Error(err),
		)
		return
	}
	//根据id获取社区数据
	community, err := mysql.GetCommunityDetailById(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailById(post.CommunityID) failed",
			zap.Int64("post.CommunityID", post.CommunityID),
			zap.Error(err),
		)
		return
	}
	//拼装数据
	date = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

func GetPostList(page, size int64) (date []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return
	}
	//处理每个帖子详情
	date = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		//获取用户信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("post.AuthorID", post.AuthorID),
				zap.Error(err),
			)
			continue
		}
		//获取社区信息
		community, err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error("",
				zap.Int64("post.CommunityID", post.CommunityID),
				zap.Error(err),
			)
			continue
		}
		//格式处理
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}

		date = append(date, postDetail)
	}

	return
}

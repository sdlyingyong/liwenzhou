package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(post *models.Post) (err error) {
	//post id 生成器
	post.ID, err = snowflake.GenId()
	if err != nil {
		return
	}
	//数据库存储器
	err = mysql.CreatePost(post)
	if err != nil {
		zap.L().Error("mysql.CreatePost(post) failed", zap.Error(err))
		return
	}
	err = redis.CreatePost(post.ID, post.CommunityID)
	if err != nil {
		zap.L().Error("redis.CreatePost(post) failed", zap.Error(err))
		return
	}
	//结果返回
	return
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

//按照时间/分数获取帖子列表
//1.获取参数
//2,去redis查询id列表
//3.根据id去数据库查询帖子详细信息
func GetPostListOrder(p *models.ParamPostList) (date []*models.ApiPostDetail, err error) {
	//1.从redis获取按照时间/分数排序的ids
	ids, err := redis.GetPostList2(p)
	if err != nil {
		zap.L().Error("redis.GetPostList2(p) failed",
			zap.Error(err))
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostList2(p) return 0 len data")
		return
	}
	//2.根据ids去mysql获取详细数据
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		zap.L().Error("mysql.GetPostListByIDs(ids) failed",
			zap.Error(err))
		return
	}
	//提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		zap.L().Error("redis.GetPostVoteData(ids) failed",
			zap.Error(err))
		return
	}

	//将帖子作者和分区信息填充到列表中
	//处理每个帖子详情
	date = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
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
			Votes:           voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}

		date = append(date, postDetail)
	}
	return
}

//按照分区获取排序号的帖子IDs
func GetCommunityPostIDsInOrder(p *models.ParamPostList) (date []*models.ApiPostDetail, err error) {
	//1.从redis获取分区排序好的postIDs
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		zap.L().Error("redis.GetCommunityPostIDsInOrder(p) failed",
			zap.Error(err))
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostIDsInOrder(p) return 0 len data")
		return
	}
	//2.根据ids去mysql获取详细数据
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		zap.L().Error("mysql.GetPostListByIDs(ids) failed",
			zap.Error(err))
		return
	}
	//提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		zap.L().Error("redis.GetPostVoteData(ids) failed",
			zap.Error(err))
		return
	}
	//将帖子作者和分区信息填充到列表中
	date = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
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
			Votes:           voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		date = append(date, postDetail)
	}
	return
}

//可选社区id的排序post帖子列表
func GetPostListHandler2(p *models.ParamPostList) (date []*models.ApiPostDetail, err error) {
	//根据请求参数是否传递,执行不同的逻辑
	if p.CommunityID == 0 {
		//排序 无社区
		date, err = GetPostListOrder(p)
	} else {
		//排序 指定社区
		date, err = GetCommunityPostIDsInOrder(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed",
			zap.Error(err))
		return
	}
	return
}

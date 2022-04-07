package mysql

import (
	"bluebell/models"
	"strings"

	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"
)

//创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post 
			(post_id,title,content,author_id,community_id) 
			values (?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

//根据id获取帖子详情
func GetPostById(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id, title, content, author_id, community_id ,create_time 
			from post 
			where post_id = ?`
	err = db.Get(post, sqlStr, pid)
	return
}

//查询帖子列表
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `SELECT
		post_id, title, content, author_id, community_id, create_time
		FROM post
		ORDER BY  create_time DESC
		limit ?,?`
	posts = make([]*models.Post, 0, 2) //make(切片,现在长度,容量)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}

//根据给定id列查询帖子数据
//返回时按照排序号的id号返回列表
func GetPostListByIDs(ids []string) (posts []*models.Post, err error) {
	sqlStr := `SELECT post_id, title, content, author_id, community_id, create_time
		FROM post
		WHERE post_id in (?)
		ORDER BY FIND_IN_SET(post_id,?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		zap.L().Error(`sqlx.In(sqlStr, ids, strings.Join(ids, ",")) failed`,
			zap.Error(err))
		return
	}
	query = db.Rebind(query)
	err = db.Select(&posts, query, args...) //!!!
	if err != nil {
		zap.L().Error("db.Select(&posts, query, args) failed",
			zap.Error(err))
		return
	}
	return
}

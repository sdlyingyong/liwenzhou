package models

import "time"

//内存对齐概念
type Post struct {
	ID          int64     `json:"post_id" db:"post_id"`
	AuthorID    int64     `json:"author_id" db:"author_id" `
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

//帖子详情接口
type ApiPostDetail struct {
	AuthorName       string `json:"author_name"`
	Votes            int64  `json:"votes"`
	*Post                   //嵌入帖子结构体
	*CommunityDetail        //嵌入社区信息
}

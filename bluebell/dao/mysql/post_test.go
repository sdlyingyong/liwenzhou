package mysql

import (
	"bluebell/models"
	"bluebell/settings"
	"testing"
	"time"
)

//初始化mysql连接
func init() {
	dbCfg := settings.MysqlConfig{
		Host:         "127.0.0.1",
		Port:         "3306",
		User:         "root",
		Password:     "root",
		Dbname:       "bluebell",
		MaxOpenConns: "20",
		MaxIdleConns: "10",
	}
	err := Init(&dbCfg)
	if err != nil {
		panic(err)
	}
}

func reset() {
	db.Close()
}

//测试创建帖子
func TestCreatePost(t *testing.T) {
	post := models.Post{
		ID:          1000,
		AuthorID:    123,
		CommunityID: 1,
		Status:      0,
		Title:       "test",
		Content:     "just a test",
		CreateTime:  time.Time{},
	}
	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("CreatePost() failed, err: %v\n", err)
	}
	t.Logf("CreatePost() success")
}

package logic

import (
	mysql "bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunityList() (data interface{}, err error) {
	//查数据库
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailById(id)
}

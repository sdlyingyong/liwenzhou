package mysql

import (
	"bluebell/models"
	"database/sql"

	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := `select community_id, community_name from community`
	if err = db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

func GetCommunityDetailById(id int64) (community *models.CommunityDetail, err error) {
	sqlStr := `
		select community_id, community_name, introduction, create_time 
		from community
		where community_id = ?`
	community = new(models.CommunityDetail) //初始化对象,并返回一个引用地址
	if err = db.Get(community, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = ErrorInvalidId
		}
	}
	return
}

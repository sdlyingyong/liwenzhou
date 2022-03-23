package logic

import (
	"lwz/bluebell/dao/redis"
	"lwz/bluebell/models"
	"strconv"

	"go.uber.org/zap"
)

//投票功能
func VoteForPost(userID int64, p *models.ParamVoteDate) (err error) {
	direction, err := strconv.ParseFloat(p.Direction, 10)
	if err != nil {
		zap.L().Debug("strconv.ParseFloat(p.Direction, 10) failed", zap.Error(err))
		return
	}
	return redis.VoteForPost(strconv.FormatInt(userID, 10), p.PostID, direction)
}

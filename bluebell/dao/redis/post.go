package redis

import (
	"time"

	"go.uber.org/zap"

	"github.com/go-redis/redis"
)

func CreatePost(pId int64) (err error) {
	//发帖时间存储
	now := float64(time.Now().Unix())
	_, err = client.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  now,
		Member: pId,
	}).Result()
	if err != nil {
		zap.L().Debug("client.ZAdd(getRedisKey(KeyPostTimeZSet) failed",
			zap.Error(err))
		return
	}
	//帖子初始score存储
	_, err = client.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  now,
		Member: pId,
	}).Result()
	return
}

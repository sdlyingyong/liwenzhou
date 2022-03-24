package redis

import (
	"lwz/bluebell/models"
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

func GetPostList2(p *models.ParamPostList) (postIDs []string, err error) {
	//排序
	key := getRedisKey(KeyPostTimeZSet)
	//1.根据用户排序参数决定post id数据来源
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	} else {
		key = getRedisKey(KeyPostTimeZSet)
	}
	zap.L().Debug("GetPostList2", zap.Any("key", key))
	//2.确定查询的起止位置
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	return client.ZRevRange(key, start, end).Result()
}

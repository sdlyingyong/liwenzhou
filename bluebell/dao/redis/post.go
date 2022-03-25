package redis

import (
	"lwz/bluebell/models"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/go-redis/redis"
)

//按照排序获取postIDs
func getIDsFormKey(key string, page, size int64) (postIDs []string, err error) {
	start := (page - 1) * size
	end := start + size - 1
	postIDs, err = client.ZRevRange(key, start, end).Result()
	if err != nil {
		zap.L().Error("client.ZRevRange(key, start, end).Result() failed",
			zap.Error(err),
		)
		return
	}
	return
}

func CreatePost(pId, communityID int64) (err error) {
	//使用事务方式执行
	pipeline := client.Pipeline()
	//发帖时间存储
	now := float64(time.Now().Unix())
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  now,
		Member: pId,
	})
	//帖子初始score存储
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  now,
		Member: pId,
	})
	//帖子id存储到对应社区的set中
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, pId)
	//事务执行
	_, err = pipeline.Exec()
	if err != nil {
		zap.L().Debug("pipeline.Exec() failed",
			zap.Error(err))
		return
	}
	return
}

func GetPostList2(p *models.ParamPostList) (postIDs []string, err error) {
	//排序
	key := getRedisKey(KeyPostTimeZSet)
	//1.根据用户排序参数决定post id数据来源
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	//2.确定查询的起止位置
	return getIDsFormKey(key, p.Page, p.Size)
}

//获取分社区的排序postIDs
//使用 zInterStore 把分区的帖子set与帖子分数的zset生成新的zset
func GetCommunityPostIDsInOrder(p *models.ParamPostList) (ids []string, err error) {
	//1.根据用户排序参数决定post id数据来源
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	//社区的key
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	//利用缓存key减少zInterStore执行的次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(key).Val() < 1 {
		//不存在,需要计算
		zap.L().Debug("不存在,需要计算",
			zap.Any("key", key),
			zap.Any("cKey", cKey),
			zap.Any("orderKey", orderKey),
		)
		pipeline := client.Pipeline()
		pipeline.ZInterStore(
			key,
			redis.ZStore{Aggregate: "MAX"},
			cKey,
			orderKey,
		) //zInterStore 计算
		pipeline.Expire(key, 5*time.Minute) //设置超时时间
		_, err = pipeline.Exec()
		if err != nil {
			zap.L().Error("pipeline.Exec() failed",
				zap.Error(err))
			return
		}
	}
	//针对新的zSet按照之前的逻辑取数据
	//存在的话直接获取
	return getIDsFormKey(key, p.Page, p.Size)
}

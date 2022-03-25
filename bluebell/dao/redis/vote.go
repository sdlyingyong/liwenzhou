package redis

import (
	"errors"
	"math"
	"time"

	"go.uber.org/zap"

	"github.com/go-redis/redis"
)

/**
投票的几种情况
direction =1,有两种情况
	1.之前没投票,现在投赞成票		->更新分数和投票记录	差值的绝对值: 1  +432
	2.之前投反对票,现在投赞成票	->更新分数和记录投票 差值的绝对值: 2  +432*2

direction =0时,有两种情况
	1.之前投过赞成票,现在取消投票		->更新分数和投票记录 差值的绝对值: 1  -432
	2.之前投过反对票,现在要取消投票	->更新分数和记录投票 差值的绝对值: 1  +432

direction =-1,有两种情况
	1.之前没有投过票,现在投反对票		->更新分数和投票记录 差值的绝对值: 1  -432
	2.之前投赞成票,现在投反对票		->更新分数和投票记录 差值的绝对值: 2  -432*2

投票的限制:
每个帖子自发表之日第一个星期之内允许用户投票,超过一个星期就不允许再投票了
	1.到期之后将redis中保存的赞成票数和反对票数存储到mysql中
	2.到期之后删除那个 KeyPostVotedZSetPF

*/

const (
	oneWeekInSeconds = float64(7 * 24 * 3600)
	scorePerVote     = 432 //每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

//投票功能
//简化版的投票分数算法
//投一票+432分
func VoteForPost(userID, postID string, value float64) (err error) {
	//1.判断投票限制(一周内可以投票)
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	now := float64(time.Now().Unix())
	if now-postTime > oneWeekInSeconds {
		zap.L().Debug("now - postTime > oneWeekInSeconds failed",
			zap.Float64("now", now),
			zap.Float64("postTime", postTime))
		err = ErrVoteTimeExpire
		return
	}
	//2.更新帖子的分数
	//查询用户目前的投票记录
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	if ov == value {
		err = ErrVoteRepeated
		return
	}
	//操作产生的分数增加或减少到帖子的score上
	pipeline := client.TxPipeline()
	diff := math.Abs(ov - value)
	var op float64
	if op > ov {
		op = 1
	} else {
		op = -1
	}
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)
	//3.记录用户为该帖子投票的数据
	if value == 0 {
		//未投票状态 移除帖子投票数中的用户
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPF+postID), userID)
	} else {
		//已投票状态 将投票结果存在帖子
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value, //赞成还是反对票
			Member: userID,
		})
	}
	_, err = pipeline.Exec()
	return
}

//查询每篇帖子的投票数
func GetPostVoteData(ids []string) (date []int64, err error) {
	////遍历获取投票数
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPF + id)
	//	//查找key中分数是1的元素数量
	//	v := client.ZCount(key, "1", "1").Val() //返回分数范围内的成员数量
	//	date = append(date, v)
	//}

	//一次请求拿回所有数据
	date = make([]int64, 0, len(ids))
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		zap.L().Error("pipeline.Exec() failed", zap.Error(err))
		return
	}
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		date = append(date, v)
	}

	return
}

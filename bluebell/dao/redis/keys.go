package redis

const (

	//redis key 注意使用命名空间的方式,方便查询和拆分

	//1.key开头容易找
	Prefix             = "bluebell:"
	KeyPostTimeZSet    = "post:time"   //zset:帖子及发帖时间
	KeyPostScoreZSet   = "post:score"  //zset:帖子及投票的分数
	KeyPostVotedZSetPF = "post:voted:" //zset:记录用户级投票类型;参数是post id

	KeyCommunitySetPF = "community:" //set:保存每个分区下帖子的id
)

func getRedisKey(key string) string {
	return Prefix + key
}

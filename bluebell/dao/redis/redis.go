package redis

import (
	"fmt"
	"lwz/bluebell/settings"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	rdb *redis.Client
)

func Init(cfg *settings.RedisConfig) (err error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port")),
		Password: cfg.Password,
		DB:       cfg.Db,
		PoolSize: cfg.PoolSize,
	})
	_, err = rdb.Ping().Result()
	if err != nil {
		zap.L().Error("conn redis fail, err :", zap.Error(err))
		return
	}
	return nil
}

func Close() {
	err := rdb.Close()
	if err != nil {
		zap.L().Error("close redis fail, err :", zap.Error(err))
	}

}

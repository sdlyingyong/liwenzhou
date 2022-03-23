package redis

import (
	"fmt"
	"lwz/bluebell/settings"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	client *redis.Client
)

func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port")),
		Password: cfg.Password,
		DB:       cfg.Db,
		PoolSize: cfg.PoolSize,
	})
	_, err = client.Ping().Result()
	if err != nil {
		zap.L().Error("conn redis fail, err :", zap.Error(err))
		return
	}
	return nil
}

func Close() {
	err := client.Close()
	if err != nil {
		zap.L().Error("close redis fail, err :", zap.Error(err))
	}

}

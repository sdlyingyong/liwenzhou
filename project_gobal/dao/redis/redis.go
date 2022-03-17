package redis

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	rdb *redis.Client
)

func Init() (err error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.pool_size"),
	})
	_, err = rdb.Ping().Result()
	if err != nil {
		err = errors.New("conn redis fail")
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

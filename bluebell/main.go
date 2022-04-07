package main

import (
	"bluebell/controller"
	mysql "bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/routes"
	"bluebell/settings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	//1.配置加载器
	if err := settings.Init(); err != nil {
		zap.L().Error("Init settings failed, err :", zap.Error(err))
		return
	}

	//2.日志收集器
	if err := logger.Init(settings.Conf, settings.Conf.Mode); err != nil {
		zap.L().Error("Init logger failed, err :", zap.Error(err))
		return
	}
	defer zap.L().Sync()
	//zap.L().Debug("logger init success...")

	//3.数据库连接器
	if err := mysql.Init(settings.Conf.MysqlConfig); err != nil {
		zap.L().Error("Init mysql failed, err :", zap.Error(err))
		return
	}
	defer mysql.Close()

	//4.缓存数据库连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		zap.L().Error("Init redis failed, err :", zap.Error(err))
		return
	}
	defer redis.Close()

	//5.snow id 生成器
	if err := snowflake.Init(viper.GetString("start_time"), viper.GetInt64("machine_id")); err != nil {
		zap.L().Error("Init snowflake failed, err :", zap.Error(err))
		return
	}

	//6.翻译器
	if err := controller.InitTrans("zh"); err != nil {
		zap.L().Error("Init trans failed, err :", zap.Error(err))
		return
	}

	//6.路由转发器
	r := routes.Setup(settings.Conf.Mode)

	//7.服务器
	r.Run(":8888")

}

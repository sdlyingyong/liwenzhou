package main

import (
	"go.uber.org/zap"
	"lwz/project_gobal/dao/msyql"
	"lwz/project_gobal/dao/redis"
	"lwz/project_gobal/logger"
	"lwz/project_gobal/routes"
	"lwz/project_gobal/settings"
)

func main() {
	//1.配置加载器
	if err := settings.Init(); err != nil {
		zap.L().Error("Init settings failed, err :", zap.Error(err))
		return
	}

	//2.日志收集器
	if err := logger.Init(); err != nil {
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
	if err := redis.Init(); err != nil {
		zap.L().Error("Init redis failed, err :", zap.Error(err))
		return
	}
	defer redis.Close()

	//5.路由转发器
	r := routes.Setup()

	//6.服务器
	r.Run(":8080")

}

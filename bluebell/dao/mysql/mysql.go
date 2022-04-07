package mysql

import (
	"bluebell/settings"
	"fmt"

	_ "github.com/go-sql-driver/mysql" //用来驱动注册到全局map
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	db *sqlx.DB
)

func Init(cfg *settings.MysqlConfig) (err error) {
	//时间time.time 转为时间戳 按照当地时区 => parseTime=True
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Dbname,
	)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		//panic(err)
		zap.L().Error("connect DB failed, err : ", zap.Error(err))
	}
	//必须连接 出问题就直接panic
	db = sqlx.MustConnect("mysql", dsn)

	//最大连接数
	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
	//最大处理数
	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))
	return
}

func Close() {
	err := db.Close()
	if err != nil {
		zap.L().Error("close mysql fail, err :", zap.Error(err))
	}
}

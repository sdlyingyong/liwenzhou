package settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	Port         string `mapstructure:"port"`
	*LogConfig   `mapstructure:"log"`
	*MysqlConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    string `mapstructure:"max_size"`
	MaxAge     string `mapstructure::"max_age"`
	MaxBackups string `mapstructure:"max_backups"`
}

type MysqlConfig struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Dbname       string `mapstructure:"dbname"`
	MaxOpenConns string `mapstructure:"max_open_conns"`
	MaxIdleConns string `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:host`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Db       string `mapstructure:"db""`
	PoolSize string `mapstructure:"pool_size""`
}

var (
	Conf = new(AppConfig)
)

func Init() (err error) {

	//方法1:读取指定的配置文件
	//viper.SetConfigFile("config.yaml")	//指定读这个文件

	//方法2:在目录下搜索指定名字配置文件
	viper.AddConfigPath(".")      //查找路径(相对路径)
	viper.SetConfigName("config") //读取配置文件名
	viper.SetConfigType("yaml")   //配置文件类型,用来远程etcd获取配置信息格式

	//配置读取并加载
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			//文件未找到错误提示
			fmt.Println("viper.ReadInConfig failed. err :", err)
		} else {
			//文件找到,加载错误提示
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}

	//配置文件信息赋值给Conf结构变量中使用
	err = viper.Unmarshal(Conf)
	if err != nil {
		fmt.Println("load config with struct fail, err :", err)
		return
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了")
		err := viper.Unmarshal(Conf)
		if err != nil {
			fmt.Println("热加载配置文件时出错: ", err)
			return
		}
		fmt.Println("新配置文件加载成功", Conf)
	})
	return
}

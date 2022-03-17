package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	//showLoadFile()
	//showFindSort()
	//showWrite()
	showWatchConf()
}

func showWatchConf() {
	//配置加载
	viper.SetConfigFile("./conf/watch.yaml")
	viper.ReadInConfig()
	//更新检查
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("conf is change.reload...", in.String())
	})
	//服务器
	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "version: "+viper.GetString("version"))
	})
	r.Run()
}

func showWrite() {
	//把当前配置写入文件
	//文件选择
	viper.AddConfigPath("./conf")
	viper.SetConfigFile("tmpConf")

	viper.WriteConfig()     //覆盖写入
	viper.SafeWriteConfig() //不覆盖

	viper.WriteConfigAs("./conf/wca.conf") //覆盖
	//viper.SafeWriteConfig("./conf/swc.conf")
	viper.SafeWriteConfigAs("./conf/swca.conf") //不覆盖
}

//问题1:如果同时存在./conf/config.json和./conf/config.yaml两个配置文件的话，viper会从哪个配置文件加载配置呢？
//json生效,因为方法searchInPath中找到,顺序是 var SupportedExts = []string{"json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl", "tfvars", "dotenv", "env", "ini"}
//因为在每个目录下搜索时,是遍历这个数组查找对应格式,json在第一个,所以读取json

//问题2.在上面两个语句下搭配使用viper.SetConfigType("yaml")指定配置文件类型可以实现预期的效果吗？
//解答:不会生效,因为方法searchInPath种,因为遍历ext查找在判断v.type之前,找打josn会直接读取文件
func showFindSort() {
	viper.Debug()
	viper.SetConfigName("config")
	viper.AddConfigPath("./conf")
	viper.SetConfigType("yaml")
	//配置读取并加载
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("not found config file")
			return
			//文件未找到错误提示
		} else {
			//文件找到,加载错误提示
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}
	fmt.Println("viper find config file success.")
}

func showLoadFile() {
	//配置文件设置
	viper.SetConfigFile("./config.yaml")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	//查找设定
	viper.AddConfigPath("/etc/appname")
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath(".")
	//配置读取并加载
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			//文件未找到错误提示
		} else {
			//文件找到,加载错误提示
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}

	fmt.Println("load config success.")
}

func showDefaultSet() {
	viper.SetDefault("ContentDir", "content")
	viper.SetDefault("LayoutDir", "layouts")
	viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})

}

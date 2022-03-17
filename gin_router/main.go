package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	shopGroup := r.Group("/shop", func1, func2)
	shopGroup.Use(func3)
	{
		shopGroup.GET("/index", func4, func5)
	}

	r.Run()
}

func func1(c *gin.Context) {
	fmt.Println("func1")
}
func func2(c *gin.Context) {
	fmt.Println("func2 start")
	c.Next() //当前调用链index++,直接从后一个调用链开始执行方法
	fmt.Println("func2 end")
}
func func3(c *gin.Context) {
	fmt.Println("func3 start ")
	c.Next()
	c.Abort() //调用链index 直接指到最大值,调用链退出
	fmt.Println("func3 end ")
}
func func4(c *gin.Context) {
	fmt.Println("func4")
	c.Set("token_user", "李文周")
}
func func5(c *gin.Context) {
	fmt.Println("func5")
	v, ok := c.Get("name")
	if ok {
		fmt.Println("name is ", v.(string))
	}
}

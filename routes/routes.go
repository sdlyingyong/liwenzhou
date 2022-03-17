package routes

import (
	"fmt"
	"lwz/bluebell/controllers"
	"lwz/bluebell/logger"
	"lwz/bluebell/pkg/jwt"
	sf "lwz/bluebell/pkg/snowflake"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//业务路由注册器
func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(zap.L()), gin.Recovery())

	r.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "hello gin")
	})
	r.GET("/ping", JWTAuthMiddleware(), func(context *gin.Context) {
		context.String(http.StatusOK, "pong")
	})

	r.GET("/gen_id", func(context *gin.Context) {
		id, _ := sf.GenId()
		context.String(http.StatusOK, fmt.Sprintf("%d", id))
	})

	r.POST("/signup", controllers.SignUpHandle)
	r.POST("/login", controllers.LoginHandle)

	return r
}

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("user_id", mc.UserID)
		c.Next() // 后续的处理函数可以用过c.Get("user_id")来获取当前请求的用户信息
	}
}

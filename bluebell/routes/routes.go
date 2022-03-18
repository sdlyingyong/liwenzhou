package routes

import (
	"fmt"
	"lwz/bluebell/controllers"
	"lwz/bluebell/logger"
	"lwz/bluebell/middlewares"
	sf "lwz/bluebell/pkg/snowflake"
	"net/http"

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
	r.POST("/ping", middlewares.JWTAuthMiddleware(), func(context *gin.Context) {
		userID, err := controllers.GetCurrentUser(context)
		if err != nil {
			controllers.ResponseError(context, controllers.CodeNeedAuth)
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"msg":  "hello",
			"data": gin.H{"user_id": userID},
		})
	})

	r.GET("/gen_id", func(context *gin.Context) {
		id, _ := sf.GenId()
		context.String(http.StatusOK, fmt.Sprintf("%d", id))
	})

	r.POST("/signup", controllers.SignUpHandle)
	r.POST("/login", controllers.LoginHandle)
	r.POST("/refresh", controllers.RefreshHandle)

	return r
}

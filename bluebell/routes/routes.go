package routes

import (
	"lwz/bluebell/controllers"
	"lwz/bluebell/logger"
	"lwz/bluebell/middlewares"
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

	//前缀路由
	apiV1 := r.Group("/api/v1")
	{
		//默认
		apiV1.GET("/", func(context *gin.Context) {
			context.String(http.StatusOK, "hello gin")
		})
		apiV1.POST("/ping", middlewares.JWTAuthMiddleware(), func(context *gin.Context) {
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
		//user
		apiV1.POST("/signup", controllers.SignUpHandle)
		apiV1.POST("/login", controllers.LoginHandle)
		apiV1.POST("/refresh", controllers.RefreshHandle)
	}

	return r
}

package routes

import (
	"lwz/bluebell/controller"
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
			userID, err := controller.GetCurrentUser(context)
			if err != nil {
				controller.ResponseError(context, controller.CodeNeedAuth)
				return
			}
			context.JSON(http.StatusOK, gin.H{
				"msg":  "hello",
				"data": gin.H{"user_id": userID},
			})
		})
		//user
		apiV1.POST("/signup", controller.SignUpHandle)
		apiV1.POST("/login", controller.LoginHandle)
		apiV1.POST("/refresh", controller.RefreshHandle)

		apiV1.Use(middlewares.JWTAuthMiddleware())
		{
			apiV1.GET("/community", controller.CommunityHandler)
			apiV1.GET("/community/:id", controller.CommunityDetailHandler)

			apiV1.POST("/post", controller.CreatePostHandler)
			apiV1.GET("/post/:id", controller.GetPostDetailHandler)
			apiV1.GET("/posts", controller.GetPostListHandler)
			apiV1.GET("/posts2", controller.GetPostListHandler2)

			apiV1.POST("/vote", controller.PostVoteHandler)
		}
	}

	//404处理
	r.NoRoute(func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}

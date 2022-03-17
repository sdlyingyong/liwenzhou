package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"lwz/project_gobal/logger"
	"net/http"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(zap.L()), gin.Recovery())

	r.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "hello gin")
	})

	return r
}

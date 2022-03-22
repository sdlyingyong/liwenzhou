package controller

import (
	"lwz/bluebell/logic"
	"lwz/bluebell/models"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

//发布帖子
func CreatePostHandler(ctx *gin.Context) {
	//参数处理器
	p := new(models.Post)
	if err := ctx.ShouldBindJSON(p); err != nil {
		zap.L().Error("create post with invalid param", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	//从context获取当前请求的user_id
	userId, err := GetCurrentUser(ctx)
	if err != nil {
		ResponseError(ctx, CodeInvalidAuth)
		return
	}
	p.AuthorID = userId
	//逻辑处理器
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	//响应处理器
	ResponseSuccess(ctx, nil)
}

//获取帖子详情
func GetPostDetailHandler(ctx *gin.Context) {
	//参数处理器
	pidStr := ctx.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	//逻辑处理器
	data, err := logic.GetPostById(pid)
	if err != nil {
		zap.L().Error("logic.GetPostById(pid) failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	//响应处理器
	ResponseSuccess(ctx, data)
}

//获取帖子列表
func GetPostListHandler(ctx *gin.Context) {
	//参数处理
	page, size := GetPageInfo(ctx)
	//逻辑处理
	date, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	//响应处理
	ResponseSuccess(ctx, date)
}

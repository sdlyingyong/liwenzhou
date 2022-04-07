package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Description 发布帖子
// @Param title query string true "标题"
// @Param content query string true "内容"
// @Param community_id query int true "社区id"
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
		ResponseError(ctx, CodeNeedAuth)
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

//升级版帖子接口
//按照时间排序或者分数排序
//api/v1/post2?page=1&size=10&order=time/score
func GetPostListHandler2(c *gin.Context) {
	//参数处理
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	//get请求中参数是在?后,所以用c.ShouldBindQuery()
	//如果请求参数是json格式,用 c.ShouldBindJSON()
	if err := c.ShouldBindQuery(&p); err != nil {
		zap.L().Error("c.ShouldBindQuery(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//逻辑处理
	date, err := logic.GetPostListHandler2(p)
	if err != nil {
		zap.L().Error("logic.GetPostList2(&p) failed",
			zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//响应处理
	ResponseSuccess(c, date)
}

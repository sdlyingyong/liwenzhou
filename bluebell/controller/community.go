package controller

import (
	"lwz/bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ------社区相关的------
func CommunityHandler(ctx *gin.Context) {
	//查询到所有的社区 (community_id .community_name)以列表形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		//不轻易把服务端报错暴露给外面
		zap.L().Error("logic.GetCommunityList()", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	ResponseSuccess(ctx, data)
}

func CommunityDetailHandler(ctx *gin.Context) {
	//参数处理器
	cId := ctx.Param("id")
	id, err := strconv.ParseInt(cId, 10, 64)
	if err != nil {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	//逻辑处理器
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail() failed", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}
	//响应处理器
	ResponseSuccess(ctx, data)
}

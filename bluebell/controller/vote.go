package controller

import (
	"lwz/bluebell/logic"
	"lwz/bluebell/models"

	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

//投票
func PostVoteHandler(c *gin.Context) {
	//参数校验
	p := new(models.ParamVoteDate)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			zap.L().Error("err.(validator.ValidationErrors) failed",
				zap.Any("ParamVoteDate", p),
				zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}
		//翻译并移除错误提示中的结构体信息
		errData := RemoveTopStruct(errs.Translate(trans))
		zap.L().Error("c.ShouldBindJSON(p) failed", zap.Any("ParamVoteDate", p),
			zap.Any("errData", errData))
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
		return
	}
	//逻辑处理
	userID, err := GetCurrentUser(c)
	if err != nil {
		zap.L().Error("GetCurrentUser(c) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	err = logic.VoteForPost(userID, p)
	if err != nil {
		zap.L().Error("logic.VoteForPost(userID, p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//响应处理
	ResponseSuccess(c, nil)
}

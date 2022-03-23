package controller

import (
	"errors"
	mysql "lwz/bluebell/dao/mysql"
	"lwz/bluebell/logic"
	"lwz/bluebell/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type LoginResponse struct {
	//userID在json格式的时候转为string类型
	UserID    int64  `json:"user_id,string"`
	UserName  string `json:"user_name"`
	UserToken string `json:"user_token"`
}

func SignUpHandle(ctx *gin.Context) {
	//1.参数校验器
	p := new(models.ParamSignUp) //创建对象,返回引用
	//字段类型判断和json格式判断器
	if err := ctx.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		errors, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(ctx, CodeInvalidParam)
			return
		} else {
			ResponseErrorWithMsg(ctx, CodeInvalidParam, RemoveTopStruct(errors.Translate(trans)))
			return
		}
	}
	//参数判断手动操作器
	if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.RePassword != p.Password {
		ResponseError(ctx, CodeInvalidParam)
		return
	}
	//2.业务处理器
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp err", zap.Error(err))
		ResponseErrorWithMsg(ctx, CodeServerBusy, "注册失败")
		return
	}
	//3.响应处理器
	ResponseSuccess(ctx, nil)
}

func LoginHandle(ctx *gin.Context) {
	//参数检查器
	p := new(models.ParamLogin) //创建对象,返回引用
	if err := ctx.ShouldBindJSON(p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		errors, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseErrorWithMsg(ctx, CodeServerBusy, err)
			return
		} else {
			ResponseErrorWithMsg(ctx, CodeServerBusy, RemoveTopStruct(errors.Translate(trans)))
			return
		}
	}
	//逻辑处理器
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login err", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(ctx, CodeUserNotExist)
			return
		}
		ResponseError(ctx, CodeInvalidPassword)
		return
	}
	//响应处理器
	//ResponseSuccess(ctx, gin.H{
	//	//因为js解析int类型有上限,超过会转错误值,所以改成字符串类型传递
	//	//"user_id":    strconv.FormatInt(user.UserID, 10),
	//	//"user_id":    fmt.Sprintf("%d", user.UserID),
	//	"user_id":    user.UserID,
	//	"user_name":  user.Username,
	//	"user_token": user.Token,
	//})
	ResponseSuccess(ctx, &LoginResponse{
		UserID:    user.UserID,
		UserName:  user.Username,
		UserToken: user.Token,
	})
}

func RefreshHandle(ctx *gin.Context) {
	//参数检查器
	p := new(models.ParamRefresh) //创建对象,返回引用
	if err := ctx.ShouldBindJSON(p); err != nil {
		zap.L().Error("Refresh token with invalid param", zap.Error(err))
		errors, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseErrorWithMsg(ctx, CodeServerBusy, err)
			return
		} else {
			ResponseErrorWithMsg(ctx, CodeServerBusy, RemoveTopStruct(errors.Translate(trans)))
			return
		}
	}
	//逻辑转发器
	token, err := logic.RefreshToken(p)
	if err != nil {
		ResponseErrorWithMsg(ctx, CodeServerBusy, err.Error())
		return
	}
	//响应处理器
	ResponseSuccess(ctx, gin.H{
		"access_token": token,
	})
}

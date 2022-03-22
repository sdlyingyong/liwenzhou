package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

var ErrorUserNotLogin = errors.New("用户未登录")

const (
	ContextUserIDKey = "userID"
)

//获取当前登陆用户id
func GetCurrentUser(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(ContextUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = uid.(int64) //转换为(类型),返回对应类型的值和是否成功(bool)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

//获取分页参数
func GetPageInfo(ctx *gin.Context) (page, size int64) {
	//参数处理
	pageNumStr := ctx.Query("page")
	pageSizeStr := ctx.Query("size")
	//分页参数
	page, err := strconv.ParseInt(pageNumStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return
}

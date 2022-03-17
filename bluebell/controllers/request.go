package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var ErrorUserNotLogin = errors.New("用户未登录")
const (
	ContextUserIDKey = "userID"
)


func GetCurrentUser(c *gin.Context)(userID int64, err error){
	uid,ok := c.Get(ContextUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID ,ok = uid.(int64)		//转换为(类型),返回对应类型的值和是否成功(bool)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}
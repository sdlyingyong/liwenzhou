package logic

import (
	mysql "lwz/bluebell/dao/msyql"
	"lwz/bluebell/models"
	"lwz/bluebell/pkg/jwt"
	"lwz/bluebell/pkg/snowflake"
)

//用户注册器
func SignUp(p *models.ParamSignUp) (err error) {
	//1.判断用户是否存在
	err = mysql.CheckUserExist(p.Username)
	if err != nil {
		return
	}
	//2.生成uid
	userId, err := snowflake.GenId()
	if err != nil {
		return
	}
	//构造一个user实例
	u := &models.User{
		UserID:   userId,
		Username: p.Username,
		Password: p.Password,
	}
	//3.保存数据库

	if err = mysql.InsertUser(u); err != nil {
		return
	}

	//redis...
	return
}

func Login(p *models.ParamLogin) (token string, err error) {
	u := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err = mysql.Login(u); err != nil {
		return
	}
	//生成jwt返回
	//u.Username
	token, err = jwt.GenToken(u.UserID, u.Username)
	if err != nil {
		return
	}
	return
}

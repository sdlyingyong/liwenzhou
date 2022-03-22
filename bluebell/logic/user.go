package logic

import (
	mysql "lwz/bluebell/dao/mysql"
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

func Login(p *models.ParamLogin) (accessToken, refreshToken string, err error) {
	u := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err = mysql.Login(u); err != nil {
		return
	}
	//生成jwt返回
	//u.Username
	accessToken, err = jwt.GenToken(u.UserID, u.Username)
	if err != nil {
		return
	}
	refreshToken, err = jwt.GenRefreshToken(u.UserID, u.Username)
	if err != nil {
		return
	}
	return
}
func RefreshToken(p *models.ParamRefresh) (newAccessToken string, err error) {
	//如果rToken过期,不正确 =>  重新登陆
	if err = jwt.CheckRefreshToken(p.RefreshToken); err != nil {
		return
	}
	//从旧的access-token 解析出claims数据
	if newAccessToken, err = jwt.RefreshAccessToken(p.AccessToken); err != nil {
		return
	}
	return
}

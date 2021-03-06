package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)

const secret = "ty"

func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	err = db.Get(&count, sqlStr, username)
	if err != nil {
		return
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

func InsertUser(user *models.User) (err error) {
	password := encryptPassword(user.Password)
	sqlStr := "INSERT INTO user (user_id, username, password) VALUES (?,?,?)"
	_, err = db.Exec(sqlStr, user.UserID, user.Username, password)
	if err != nil {
		return
	}
	return
}

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	h.Sum([]byte(oPassword))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *models.User) (err error) {
	oPassword := user.Password //输入的密码
	sqlStr := "SELECT user_id,username,password FROM user WHERE username = ?"
	err = db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		//查询数据库失败
		return
	}
	//判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

func GetUserById(userID int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id , username from user where user_id = ?`
	err = db.Get(user, sqlStr, userID)
	return
}

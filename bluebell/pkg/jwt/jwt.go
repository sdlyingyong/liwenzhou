package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrorInvalidAccessToken  = errors.New("无效的access_token")
	ErrorInvalidRefreshToken = errors.New("无效的refresh_token")
)

var (
	mySecret                   = []byte("夏天夏天悄悄过去")
	AccessTokenExpireDuration  = 1 * time.Hour
	RefreshTokenExpireDuration = 24 * time.Hour
)

type MyClaims struct {
	UserID             int64  `json:"user_id"`
	Username           string `json:"username"`
	jwt.StandardClaims        //官方字段	(ExpiresAt)
}

//生成jwt
func GenToken(userID int64, username string) (token string, err error) {
	claims := MyClaims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenExpireDuration).Unix(),
			Issuer:    "bluebell",
		},
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenObj.SignedString(mySecret)
}

//jwt解析器
func ParseAccessToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	mc, ok := token.Claims.(*MyClaims)
	if !ok || !token.Valid {
		return nil, ErrorInvalidAccessToken
	}
	return mc, nil
}

func ParseExpireAccessToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	v, _ := err.(*jwt.ValidationError) //转换成jwt提供的异常类型
	if v.Errors != jwt.ValidationErrorExpired {
		return nil, ErrorInvalidAccessToken
	}
	fmt.Println("is expired access")
	mc, ok := token.Claims.(*MyClaims)
	if !ok {
		return nil, ErrorInvalidAccessToken
	}
	return mc, nil
}

func GenRefreshToken(userID int64, username string) (token string, err error) {
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(RefreshTokenExpireDuration).Unix(), //过期时间
		Issuer:    "bluebell",
	}).SignedString(mySecret)
	return
}

func CheckRefreshToken(refreshToken string) (err error) {
	token, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return
	}
	_, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		err = ErrorInvalidRefreshToken
		return
	}
	return
}

func RefreshAccessToken(AccessToken string) (newAccessToken string, err error) {
	var mc MyClaims
	_, err = jwt.ParseWithClaims(AccessToken, &mc, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	//依旧有效的access_token 直接返回
	if err == nil {
		newAccessToken = AccessToken
		return
	}
	v, _ := err.(*jwt.ValidationError)
	if v.Errors == jwt.ValidationErrorExpired {
		newAccessToken, err = GenToken(mc.UserID, mc.Username)
		return
	}
	return
}

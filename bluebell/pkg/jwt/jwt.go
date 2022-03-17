package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	mySercet            = []byte("夏天夏天悄悄过去")
	TokenExpireDuration = 60 * time.Minute
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
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "bluebell",
		},
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenObj.SignedString(mySercet)
}

//jwt解析器
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return mySercet, nil
	})
	if err != nil {
		return nil, err
	}
	if v, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return v, nil
	}
	return nil, errors.New("invalid token")

}

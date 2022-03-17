package models

type ParamSignUp struct {
	Username   string `json:"username" binding:"required"` //key : value,value要冒号包裹
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

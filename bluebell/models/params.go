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

type ParamRefresh struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

//投票数据
type ParamVoteDate struct {
	PostID    string `json:"post_id" binding:"required"`                //帖子id
	Direction string `json:"direction" binding:"required,oneof=1 0 -1"` //投票方向 1=赞成票 -1=反对票 0=取消投票
}

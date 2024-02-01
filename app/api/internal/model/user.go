package model

type ParamRegisterUser struct {
	Nickname string `json:"nickname" form:"nickname"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Email    string `json:"email" form:"email"`
}

type ParamLoginUser struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type ApiUser struct {
	Uid   int64  `json:"uid"`
	Token string `json:"token"`
}

package model

type ParamRegisterUser struct {
	Nickname string `json:"nickname" form:"nickname"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Email    string `json:"email" form:"email"`
}

package model

type ParamRegisterUser struct {
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

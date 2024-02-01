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

type ParamGetUserInfo struct {
	Uid int64 `json:"uid" form:"uid"`
}

type ApiUser struct {
	Uid   int64  `json:"uid"`
	Token string `json:"token"`
}

type User struct {
	Uid      int64  `json:"uid" db:"user_id"`
	Username string `json:"username" db:"username"`
	Nickname string `json:"nickname" db:"nickname"`
	Email    string `json:"email" db:"email"`
	Profile  string `json:"profile" db:"profile"`
	JoinedAt int    `json:"joined_at" db:"joined_at"`
}

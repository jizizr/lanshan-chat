package model

import "time"

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

type ParamCheckUsername struct {
	Username string `json:"username" form:"username"`
}

type ParamModifyUserInfo struct {
	Username string `json:"username" form:"username"`
	Nickname string `json:"nickname" form:"nickname"`
	Email    string `json:"email" form:"email"`
	Profile  string `json:"profile" form:"profile"`
}

type ParamModifyPassword struct {
	OldPassword string `json:"old_password" form:"old_password"`
	NewPassword string `json:"new_password" form:"new_password"`
}

type ApiUser struct {
	Uid   int64  `json:"uid"`
	Token string `json:"token"`
}

type UserInfo struct {
	Uid      int64     `json:"user_id" db:"user_id"`
	Username string    `json:"username" db:"username"`
	Nickname string    `json:"nickname" db:"nickname"`
	Email    string    `json:"email" db:"email"`
	Profile  string    `json:"profile" db:"profile"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}

type User struct {
	Uid      int64     `json:"user_id" db:"user_id"`
	Username string    `json:"username" db:"username"`
	Nickname string    `json:"nickname" db:"nickname"`
	Email    string    `json:"email" db:"email"`
	Profile  string    `json:"profile" db:"profile"`
	Password string    `json:"password" db:"password"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}

package mysql

import (
	"lanshan_chat/app/api/global"
	"lanshan_chat/app/api/internal/model"
)

const (
	CountUserByUsernameStr = "SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)"
	AddUserStr             = "INSERT INTO users(username,password,email) VALUES (?, ?, ?)"
)

// CheckUserIsExist 如果用户存在返回 true，否则返回 false
// 如果数据库操作出错返回 error
func CheckUserIsExist(username string) (flag bool, err error) {
	err = global.MDB.Get(&flag, CountUserByUsernameStr, username)
	return
}

func AddUser(u *model.ParamRegisterUser) error {
	_, err := global.MDB.Exec(AddUserStr, u.Username, u.Password, u.Email)
	return err
}

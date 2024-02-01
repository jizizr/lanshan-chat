package mysql

import (
	"lanshan_chat/app/api/global"
)

const (
	CountUserByUsernameStr     = "SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)"
	AddUserStr                 = "INSERT INTO users(username,nickname,password,email) VALUES (?, ?, ?, ?)"
	QueryPasswordByUsernameStr = "SELECT uid,password FROM users WHERE username = ?"
	QueryUserByEmailStr        = "SELECT uid,password FROM users WHERE email = ?"
)

// CheckUserIsExist 如果用户存在返回 true，否则返回 false
// 如果数据库操作出错返回 error
func CheckUserIsExist(username string) (flag bool, err error) {
	err = global.MDB.Get(&flag, CountUserByUsernameStr, username)
	return
}

func AddUser(username, nickname, password, email string) error {
	_, err := global.MDB.Exec(AddUserStr, username, nickname, password, email)
	return err
}

func QueryPasswordByUsername(username string) (uid int64, password string, err error) {
	row := global.MDB.QueryRow(QueryPasswordByUsernameStr, username)
	if err := row.Scan(&uid, &password); err != nil {
		return 0, "", err
	}
	return
}

func QueryUserByEmail(email string) (uid int64, password string, err error) {
	row := global.MDB.QueryRow(QueryUserByEmailStr, email)
	if err := row.Scan(&uid, &password); err != nil {
		return 0, "", err
	}
	return
}

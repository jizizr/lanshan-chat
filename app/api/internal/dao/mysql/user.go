package mysql

import (
	"lanshan_chat/app/api/global"
	"lanshan_chat/app/api/internal/consts"
)

const (
	CheckUserExistByUsernameStr = "SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)"
	AddUserStr                  = "INSERT INTO users(username,nickname,password,email,profile) VALUES (?, ?, ?, ?, ?)"
	QueryPasswordByUsernameStr  = "SELECT user_id,password FROM users WHERE username = ?"
	QueryUserByEmailStr         = "SELECT user_id,password FROM users WHERE email = ?"
)

// CheckUserIsExist 如果用户存在返回 true，否则返回 false
// 如果数据库操作出错返回 error
func CheckUserIsExist(username string) (flag bool, err error) {
	err = global.MDB.Get(&flag, CheckUserExistByUsernameStr, username)
	return
}

func AddUser(username, nickname, password, email string) (int64, error) {
	result, err := global.MDB.Exec(AddUserStr, username, nickname, password, email, consts.DefultProfile)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
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

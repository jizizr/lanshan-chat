package mysql

import (
	"lanshan_chat/app/api/global"
	"lanshan_chat/app/api/internal/consts"
	"lanshan_chat/app/api/internal/model"
	"time"
)

const (
	CheckUserExistByUsernameStr = "SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)"
	CheckUserExistByUIDStr      = "SELECT EXISTS(SELECT 1 FROM users WHERE user_id = ?)"
	AddUserStr                  = "INSERT INTO users(username,nickname,password,email,profile,joined_at) VALUES (?, ?, ?, ?, ?, ?)"
	QueryPasswordByUsernameStr  = "SELECT user_id,password FROM users WHERE username = ?"
	QueryUserByEmailStr         = "SELECT user_id,password FROM users WHERE email = ?"
	QueryUserByUIDStr           = "SELECT user_id,username,nickname,email,profile,joined_at FROM users WHERE user_id = ?"
	ModifyUserInfoStr           = `UPDATE users 
						           SET username = :username,
						               nickname = :nickname,
						               email = :email,
						               profile = :profile
						           WHERE user_id = :user_id`
	QueryUserStr = "SELECT * FROM users WHERE user_id = ?"
)

// CheckUserIsExistByUsername 如果用户存在返回 true，否则返回 false
// 如果数据库操作出错返回 error
func CheckUserIsExistByUsername(username string) (flag bool, err error) {
	err = global.MDB.Get(&flag, CheckUserExistByUsernameStr, username)
	return
}

func CheckUserIsExistByUID(uid int64) (flag bool, err error) {
	err = global.MDB.Get(&flag, CheckUserExistByUIDStr, uid)
	return
}

func AddUser(username, nickname, password, email string) (int64, time.Time, error) {
	t := time.Now()
	result, err := global.MDB.Exec(AddUserStr, username, nickname, password, email, consts.DefultProfile, t)
	if err != nil {
		return -1, time.Time{}, err
	}
	uid, err := result.LastInsertId()
	return uid, t, err
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

func QueryUserByUID(uid int64) (*model.UserInfo, error) {
	user := new(model.UserInfo)
	err := global.MDB.Get(user, QueryUserByUIDStr, uid)
	return user, err
}

func ModifyUserInfo(u *model.User) error {
	_, err := global.MDB.NamedExec(ModifyUserInfoStr, u)
	return err
}

func QueryUser(uid int64) (*model.User, error) {
	user := new(model.User)
	err := global.MDB.Get(user, QueryUserStr, uid)
	return user, err
}

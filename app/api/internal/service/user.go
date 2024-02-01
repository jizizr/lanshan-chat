package service

import (
	"go.uber.org/zap"
	"lanshan_chat/app/api/global"
	"lanshan_chat/app/api/internal/consts"
	"lanshan_chat/app/api/internal/dao/mysql"
	"lanshan_chat/app/api/internal/dao/redis"
	"lanshan_chat/app/api/internal/model"
	"lanshan_chat/utils"
)

func Register(user *model.ParamRegisterUser) (int64, error) {
	flag, err := mysql.CheckUserIsExistByUsername(user.Username)
	if err != nil {
		return -1, err
	}
	if flag {
		return -1, consts.UserExistError
	}
	password := utils.CryptoPassword(user.Password)
	uid, t, err := mysql.AddUser(user.Username, user.Nickname, password, user.Email)
	if err != nil {
		return -1, err
	}
	err = redis.AddUser(uid, user.Username, user.Nickname, password, user.Email, t)
	return uid, err
}

func Login(user *model.ParamLoginUser) (int64, error) {
	// 判断用户是否存在
	flag, err := mysql.CheckUserIsExistByUsername(user.Username)
	if err != nil {
		global.Logger.Error("login failed", zap.Error(err))
		return -1, err
	}
	if !flag {
		return -1, consts.UserNotExistError
	}

	var (
		uid      int64
		password string
	)

	// 判断用户输入的是用户名还是邮箱
	if utils.Check_email(user.Username) {
		uid, password, err = mysql.QueryUserByEmail(user.Username)
	} else {
		uid, password, err = mysql.QueryPasswordByUsername(user.Username)
	}
	if err != nil {
		global.Logger.Error("login failed", zap.Error(err))
		return -1, err
	}
	if password != utils.CryptoPassword(user.Password) {
		return -1, consts.PasswordWrongError
	}
	return uid, nil
}

func GetUserProfile(uid int64) (*model.User, error) {
	// 判断用户是否存在
	flag, err := redis.CheckUserIsExist(uid)
	if err != nil {
		return nil, err
	}
	if flag {
		return redis.GetUserProfile(uid)
	}
	flag, err = mysql.CheckUserIsExistByUID(uid)
	if err != nil {
		return nil, err
	}
	if !flag {
		return nil, consts.UserNotExistError
	}
	return mysql.QueryUserByUID(uid)
}

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
	// 检查username是否被占用
	flag, err := CheckUsername(user.Username)
	if err != nil {
		return -1, err
	}
	if flag {
		return -1, consts.UserExistError
	}
	// 添加用户
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

func GetUserInfo(uid int64) (*model.UserInfo, error) {
	// 判断用户是否存在
	flag, err := redis.CheckUserIsExist(uid)
	if err != nil {
		return nil, err
	}
	if flag {
		return redis.GetUserInfo(uid)
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

func CheckUsername(username string) (flag bool, err error) {
	if flag, err = redis.CheckUsername(username); err != nil {
		return
	}
	if !flag {
		flag, err = mysql.CheckUserIsExistByUsername(username)
	}

	// 如果数据库中存在该用户，将其添加到redis缓存中
	if flag {
		redis.AddToSet(username)
	}
	return
}

func getUser(uid int64) (*model.User, error) {
	// 判断用户是否存在
	flag, err := redis.CheckUserIsExist(uid)
	if err != nil {
		return nil, err
	}
	if flag {
		return redis.GetUser(uid)
	}

	flag, err = mysql.CheckUserIsExistByUID(uid)
	if err != nil {
		return nil, err
	}
	if !flag {
		return nil, consts.UserNotExistError
	}
	return mysql.QueryUser(uid)
}

func ModifyUserInfo(uid int64, u *model.ParamModifyUserInfo) error {
	user, err := getUser(uid)
	if err != nil {
		return err
	}

	if u.Username != "" {
		flag, err := CheckUsername(u.Username)
		if err != nil {
			return err
		}
		if flag {
			return consts.UserExistError
		}

		defer func(oldUsername, newUsername string) {
			redis.DelFromSet(oldUsername)
			redis.AddToSet(newUsername)
		}(user.Username, u.Username)

		user.Username = u.Username
	}
	if u.Nickname != "" {
		user.Nickname = u.Nickname
	}
	if u.Email != "" {
		user.Email = u.Email
	}
	if u.Profile != "" {
		user.Profile = u.Profile
	}

	if err := mysql.ModifyUserInfo(user); err != nil {
		return err
	}
	return redis.ModifyUserInfo(user)
}

func ModifyPassword(uid int64, u *model.ParamModifyPassword) error {
	user, err := getUser(uid)
	if err != nil {
		return err
	}
	if user.Password != utils.CryptoPassword(u.OldPassword) {
		return consts.PasswordWrongError
	}
	user.Password = utils.CryptoPassword(u.NewPassword)
	user.Uid = uid
	if err := mysql.ModifyUserInfo(user); err != nil {
		return err
	}
	return redis.ModifyUserInfo(user)
}

func SearchUser(keyword string) ([]model.UserInfo, error) {
	return mysql.SearchUser(keyword)
}

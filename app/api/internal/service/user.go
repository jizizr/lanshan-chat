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

func Register(user *model.ParamRegisterUser) error {
	flag, err := redis.CheckUserIsExist(user.Username)
	if err != nil {
		global.Logger.Error("register failed", zap.Error(err))
		return err
	}
	if !flag {
		flag, err = mysql.CheckUserIsExist(user.Username)
		if err != nil {
			global.Logger.Error("register failed", zap.Error(err))
			return err
		}
	}
	if flag {
		return consts.UserExistError
	}
	password := utils.CryptoPassword(user.Password)
	err = mysql.AddUser(user.Username, user.Nickname, password, user.Email)
	if err != nil {
		return err
	}
	err = redis.AddUser(user.Username, user.Nickname, password, user.Email)
	if err != nil {
		return err
	}
	return nil
}

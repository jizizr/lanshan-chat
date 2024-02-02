package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"lanshan_chat/app/api/global"
	"lanshan_chat/app/api/internal/consts"
	"lanshan_chat/app/api/internal/model"
	"lanshan_chat/app/api/internal/service"
	"lanshan_chat/utils"
)

func Register(c *gin.Context) {
	u := new(model.ParamRegisterUser)
	if err := c.ShouldBind(u); err != nil {
		RespFailed(c, 400, consts.CodeShouldBind)
		return
	}
	if u.Username == "" || u.Password == "" || u.Nickname == "" || u.Email == "" {
		RespFailed(c, 400, consts.CodeParamEmpty)
		return
	}
	if !utils.Check_email(u.Email) {
		RespFailed(c, 400, consts.CodeEmailWrongFormat)
		return
	}
	if utils.Check_email(u.Username) {
		RespFailed(c, 400, consts.CodeUsernameWrongFormat)
		return
	}
	uid, err := service.Register(u)
	if err != nil {
		if errors.Is(err, consts.UserExistError) {
			RespFailed(c, 400, consts.CodeUserAlreadyExist)
			return
		}
		RespFailed(c, 500, consts.CodeDBCheckUser)
		global.Logger.Error("register failed", zap.Error(err))
		return
	}

	token, err := utils.GenToken(uid)
	if err != nil {
		RespFailed(c, 500, consts.CodeServerBusy)
		global.Logger.Error("register failed", zap.Error(err))
		return
	}
	RespSuccess(c, &model.ApiUser{
		Uid:   uid,
		Token: token,
	})
}

func Login(c *gin.Context) {
	u := new(model.ParamLoginUser)
	if err := c.ShouldBind(u); err != nil {
		RespFailed(c, 400, consts.CodeShouldBind)
		return
	}
	if u.Username == "" || u.Password == "" {
		RespFailed(c, 400, consts.CodeParamEmpty)
		return
	}
	uid, err := service.Login(u)
	if err != nil {
		if errors.Is(err, consts.UserNotExistError) {
			RespFailed(c, 400, consts.CodeUserNotExist)
			return
		}
		if errors.Is(err, consts.PasswordWrongError) {
			RespFailed(c, 400, consts.CodeWrongPassword)
			return
		}
		RespFailed(c, 500, consts.CodeDBCheckUser)
		global.Logger.Error("login failed", zap.Error(err))
		return
	}
	token, err := utils.GenToken(uid)
	if err != nil {
		RespFailed(c, 500, consts.CodeServerBusy)
		global.Logger.Error("login failed", zap.Error(err))
		return
	}
	RespSuccess(c, &model.ApiUser{
		Uid:   uid,
		Token: token,
	})
}

func GetUserInfo(c *gin.Context) {
	u := new(model.ParamGetUserInfo)
	if err := c.ShouldBind(u); err != nil {
		RespFailed(c, 400, consts.CodeShouldBind)
		return
	}
	if u.Uid == 0 {
		RespFailed(c, 400, consts.CodeParamEmpty)
		return
	}
	user, err := service.GetUserProfile(u.Uid)
	if err != nil {
		if errors.Is(err, consts.UserNotExistError) {
			RespFailed(c, 400, consts.CodeUserNotExist)
			return
		}
		RespFailed(c, 500, consts.CodeDBCheckUser)
		global.Logger.Error("get user profile failed", zap.Error(err))
		return
	}
	RespSuccess(c, user)
}

func CheckUsername(c *gin.Context) {
	u := new(model.ParamCheckUsername)
	if err := c.ShouldBind(u); err != nil {
		RespFailed(c, 400, consts.CodeShouldBind)
		return
	}
	if u.Username == "" {
		RespFailed(c, 400, consts.CodeParamEmpty)
		return
	}
	flag, err := service.CheckUsername(u.Username)
	if err != nil {
		RespFailed(c, 500, consts.CodeDBCheckUser)
		global.Logger.Error("check user is exist failed", zap.Error(err))
		return
	}
	RespSuccess(c, flag)

}

package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"lanshan_chat/app/api/internal/consts"
	"lanshan_chat/app/api/internal/model"
	"lanshan_chat/app/api/internal/service"
	"lanshan_chat/utils"
)

func Register(c *gin.Context) {
	var u *model.ParamRegisterUser
	if err := c.ShouldBindJSON(u); err != nil {
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
	if err := service.Register(u); err != nil {
		if errors.Is(err, consts.UserExistError) {
			RespFailed(c, 400, consts.CodeUserAlreadyExist)
			return
		}
		RespFailed(c, 500, consts.CodeDBCheckUser)
	}
	RespSuccess(c, nil)

}

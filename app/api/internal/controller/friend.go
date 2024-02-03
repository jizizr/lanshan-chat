package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"lanshan_chat/app/api/global"
	"lanshan_chat/app/api/internal/consts"
	"lanshan_chat/app/api/internal/model"
	"lanshan_chat/app/api/internal/service"
)

func AddFriend(c *gin.Context) {
	f := new(model.ParamAddFriend)
	if err := c.ShouldBind(f); err != nil {
		RespFailed(c, 400, consts.CodeShouldBind)
		return
	}
	if f.FriendID == 0 {
		RespFailed(c, 400, consts.CodeParamEmpty)
		return
	}
	uid, ok := GetUID(c)
	if !ok {
		RespFailed(c, 400, consts.CodeServerBusy)
		return
	}
	err := service.AddFriend(uid, f.FriendID)
	if err != nil {
		if errors.Is(err, consts.FriendExistError) {
			RespFailed(c, 400, consts.CodeFriendAlreadyExist)
			return
		}
		RespFailed(c, 500, consts.CodeDBCheckUser)
		global.Logger.Error("add friend failed", zap.Error(err))
		return
	}
	RespSuccess(c, nil)
}

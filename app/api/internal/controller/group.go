package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"lanshan_chat/app/api/global"
	"lanshan_chat/app/api/internal/consts"
	"lanshan_chat/app/api/internal/model"
	"lanshan_chat/app/api/internal/service"
	"lanshan_chat/utils"
)

func CreateGroup(c *gin.Context) {
	g := new(model.ParamCreateGroup)
	if err := c.ShouldBind(g); err != nil {
		RespFailed(c, 400, consts.CodeShouldBind)
		return
	}
	if g.GroupName == "" || g.Description == "" || g.Type == "" {
		RespFailed(c, 400, consts.CodeParamEmpty)
		return
	}
	avatar, err := c.FormFile("avatar")
	if err != nil {
		return
	}

	// 检查文件大小
	if avatar.Size > 5<<20 {
		RespFailed(c, 400, consts.CodeFileTooLarge)
		return
	}
	if avatar.Size == 0 {
		RespFailed(c, 400, consts.CodeFileEmpty)
		return
	}

	// 压缩图片
	avatarCompressed, err := utils.CompressImage(avatar)
	if err != nil {
		RespFailed(c, 500, consts.CodeCompressFailed)
		return
	}

	// 上传图片
	url, err := UploadBin(avatarCompressed, avatar.Filename)
	if err != nil {
		RespFailed(c, 500, consts.CodeServerBusy)
		return
	}

	userID, ok := GetUID(c)
	if !ok {
		RespFailed(c, 400, consts.CodeServerBusy)
		return
	}
	if err := service.CreateGroup(g, url, userID); err != nil {
		RespFailed(c, 500, consts.CodeDBCreateGroup)
		global.Logger.Error("CreateGroup failed", zap.Error(err))
		return
	}
	RespSuccess(c, nil)
}

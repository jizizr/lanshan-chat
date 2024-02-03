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

func JoinPublicGroup(c *gin.Context) {
	g := new(model.ParamJoinPublicGroup)
	if err := c.ShouldBind(g); err != nil {
		RespFailed(c, 400, consts.CodeShouldBind)
		return
	}
	if g.GroupID == 0 {
		RespFailed(c, 400, consts.CodeParamEmpty)
		return
	}
	userID, ok := GetUID(c)
	if !ok {
		RespFailed(c, 400, consts.CodeServerBusy)
		return
	}
	if err := service.JoinPublicGroup(userID, g.GroupID); err != nil {
		if errors.Is(err, consts.GroupAlreadyJoinError) {
			RespFailed(c, 400, consts.CodeGroupAlreadyJoin)
		} else if errors.Is(err, consts.BanedError) {
			RespFailed(c, 400, consts.CodeUserInGroupBanned)
		} else if errors.Is(err, consts.GroupNotExistError) {
			RespFailed(c, 400, consts.CodeGroupNotExist)
		} else if errors.Is(err, consts.GroupIsPrivateError) {
			RespFailed(c, 400, consts.CodeGroupIsPrivate)
		} else {
			RespFailed(c, 500, consts.CodeDBJoinGroup)
			global.Logger.Error("JoinGroup failed", zap.Error(err))
		}
		return
	}
	RespSuccess(c, nil)
}

func JoinPrivateGroup(c *gin.Context) {
	t := new(model.ParamJoinPrivateGroup)
	if err := c.ShouldBind(t); err != nil {
		RespFailed(c, 400, consts.CodeShouldBind)
		return
	}
	if t.Token == "" {
		RespFailed(c, 400, consts.CodeParamEmpty)
		return
	}
	userID, ok := GetUID(c)
	if !ok {
		RespFailed(c, 400, consts.CodeServerBusy)
		return
	}
	if err := service.JoinPrivateGroup(userID, t.Token); err != nil {
		if errors.Is(err, consts.GroupAlreadyJoinError) {
			RespFailed(c, 400, consts.CodeGroupAlreadyJoin)
		} else if errors.Is(err, consts.BanedError) {
			RespFailed(c, 400, consts.CodeUserInGroupBanned)
		} else if errors.Is(err, consts.InvalidTokenError) {
			RespFailed(c, 400, consts.CodeInvalidToken)
		} else if errors.Is(err, consts.GroupNotExistError) {
			RespFailed(c, 400, consts.CodeGroupNotExist)
		} else {
			RespFailed(c, 500, consts.CodeDBJoinGroup)
			global.Logger.Error("JoinGroup failed", zap.Error(err))
		}
		return
	}
	RespSuccess(c, nil)
}

func GetPrivateGroupToken(c *gin.Context) {
	g := new(model.ParamGetPrivateGroupToken)
	if err := c.ShouldBind(g); err != nil {
		RespFailed(c, 400, consts.CodeShouldBind)
		return
	}
	if g.GroupID == 0 || g.ExpiresTime == 0 {
		RespFailed(c, 400, consts.CodeParamEmpty)
		return
	}
	userID, ok := GetUID(c)
	if !ok {
		RespFailed(c, 400, consts.CodeServerBusy)
		return
	}
	token, err := service.GetPrivateGroupToken(userID, g.GroupID, g.ExpiresTime)
	if err != nil {
		if errors.Is(err, consts.GroupIsPublicError) {
			RespFailed(c, 400, consts.CodeGroupIsPublic)
		} else if errors.Is(err, consts.PermissionDeniedError) {
			RespFailed(c, 400, consts.CodePermissionDenied)
		} else {
			RespFailed(c, 500, consts.CodeServerBusy)
			global.Logger.Error("GetPrivateGroupToken failed", zap.Error(err))
		}
		return
	}
	RespSuccess(c, gin.H{"token": token})
}

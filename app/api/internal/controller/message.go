package controller

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"lanshan_chat/app/api/global"
	"lanshan_chat/app/api/internal/consts"
	"lanshan_chat/app/api/internal/model"
	"lanshan_chat/app/api/internal/service"
	"lanshan_chat/utils"
)

func handleMessages(c *gin.Context, m *model.ParamSendMessage) error {
	if err := c.ShouldBind(m); err != nil {
		RespFailed(c, 400, consts.CodeShouldBind)
		return err
	}
	if m.GroupID == 0 || m.Message == "" || m.Type == "" {
		RespFailed(c, 400, consts.CodeParamEmpty)
		return errors.New("")
	}
	switch m.Type {
	case "text":
		if global.Filter.GetFilter().IsSensitive(m.Message) {
			utils.Report()
		}
	case "image", "file", "video", "audio":
		fh, err := c.FormFile("file")
		if err != nil {
			RespFailed(c, 400, consts.CodeShouldBind)
			return err
		}
		if fh.Size == 0 {
			RespFailed(c, 400, consts.CodeFileEmpty)
			return err
		}
		f, err := fh.Open()
		if err != nil {
			RespFailed(c, 400, consts.CodeServerBusy)
			return err
		}
		defer f.Close()
		// 读取文件
		buf := make([]byte, fh.Size)
		_, err = f.Read(buf)
		if err != nil {
			RespFailed(c, 400, consts.CodeServerBusy)
			return err
		}
		// 上传文件
		url, err := UploadBin(buf, fh.Filename)
		if err != nil {
			RespFailed(c, 500, consts.CodeServerBusy)
			return err
		}
		m.Url = url
	default:
		RespFailed(c, 400, consts.CodeNotInEnum)
		return errors.New("")
	}
	userID, ok := GetUID(c)
	if !ok {
		RespFailed(c, 400, consts.CodeServerBusy)
		return errors.New("")
	}
	m.SenderID = userID
	return nil
}

func SendMessage(c *gin.Context) {
	m := new(model.ParamSendMessage)
	if err := handleMessages(c, m); err != nil {
		return
	}
	if err := service.SendMessage(m); err != nil {
		if errors.Is(err, consts.PermissionDeniedError) {
			RespFailed(c, 400, consts.CodePermissionDenied)
			return
		} else {
			RespFailed(c, 500, consts.CodeDBSendMessage)
			global.Logger.Error("SendMessage failed", zap.Error(err))
		}
		return
	}
	RespSuccess(c, m)
}

func DeleteMessage(c *gin.Context) {
	m := new(model.ParamDeleteMessage)
	if err := c.ShouldBind(m); err != nil {
		RespFailed(c, 400, consts.CodeShouldBind)
		return
	}
	if m.GroupID == 0 {
		RespFailed(c, 400, consts.CodeParamEmpty)
		return
	}
	userID, ok := GetUID(c)
	if !ok {
		RespFailed(c, 400, consts.CodeServerBusy)
		return
	}
	if err := service.DeleteMessage(userID, m); err != nil {
		if errors.Is(err, consts.PermissionDeniedError) {
			RespFailed(c, 400, consts.CodePermissionDenied)
		} else if errors.Is(err, sql.ErrNoRows) {
			RespFailed(c, 400, consts.CodeMessageNotExist)
		} else {
			RespFailed(c, 500, consts.CodeDBDeleteMessage)
			global.Logger.Error("DeleteMessage failed", zap.Error(err))
		}
		return
	}
	RespSuccess(c, nil)
}

func EditMessage(c *gin.Context) {
	m := new(model.ParamSendMessage)
	if err := handleMessages(c, m); err != nil {
		return
	}
	if err := service.EditMessage(m); err != nil {
		if errors.Is(err, consts.PermissionDeniedError) {
			RespFailed(c, 400, consts.CodePermissionDenied)
			return
		} else if errors.Is(err, sql.ErrNoRows) {
			RespFailed(c, 400, consts.CodeMessageNotExist)
		} else {
			RespFailed(c, 500, consts.CodeDBSendMessage)
			global.Logger.Error("SendMessage failed", zap.Error(err))
		}
		return
	}
	RespSuccess(c, nil)
}

func GetGroupMessage(c *gin.Context) {
	m := new(model.ParamGetMessage)
	if err := c.ShouldBind(m); err != nil {
		RespFailed(c, 400, consts.CodeShouldBind)
		return
	}
	if m.GroupID == 0 || m.StartID == 0 || m.Limit == 0 {
		RespFailed(c, 400, consts.CodeParamEmpty)
		return
	}
	userID, ok := GetUID(c)
	if !ok {
		RespFailed(c, 400, consts.CodeServerBusy)
		return
	}
	messages, err := service.GetMessage(userID, m)
	if err != nil {
		if errors.Is(err, consts.PermissionDeniedError) {
			RespFailed(c, 400, consts.CodePermissionDenied)
		} else {
			RespFailed(c, 500, consts.CodeDBCheckUser)
			global.Logger.Error("get message failed", zap.Error(err))
		}
		return
	}
	RespSuccess(c, messages)
}

func ReadMessage(c *gin.Context) {
	m := new(model.ParamReadMessage)
	if err := c.ShouldBind(m); err != nil {
		RespFailed(c, 400, consts.CodeShouldBind)
		return
	}
	if m.GroupID == 0 || m.LastRead == 0 {
		RespFailed(c, 400, consts.CodeParamEmpty)
		return
	}
	userID, ok := GetUID(c)
	if !ok {
		RespFailed(c, 400, consts.CodeServerBusy)
		return
	}
	if err := service.ReadMessage(userID, m); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			RespFailed(c, 400, consts.CodeUserNotInGroup)
		} else {
			RespFailed(c, 500, consts.CodeDBCheckUser)
			global.Logger.Error("read message failed", zap.Error(err))
		}
		return
	}
	RespSuccess(c, nil)
}

func GetLastRead(c *gin.Context) {
	m := new(model.ParamGetLastMessageID)
	if err := c.ShouldBind(m); err != nil {
		RespFailed(c, 400, consts.CodeShouldBind)
		return
	}
	if m.GroupID == 0 {
		RespFailed(c, 400, consts.CodeParamEmpty)
		return
	}
	uid, ok := GetUID(c)
	if !ok {
		RespFailed(c, 400, consts.CodeServerBusy)
		return
	}
	lastRead, err := service.GetLastRead(uid, m.GroupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			RespFailed(c, 400, consts.CodeUserNotInGroup)
		} else {
			RespFailed(c, 500, consts.CodeDBCheckUser)
			global.Logger.Error("get last read failed", zap.Error(err))
		}
		return
	}
	RespSuccess(c, &model.ApiMessageID{MessageID: lastRead})
}

func GetLastMessage(c *gin.Context) {
	m := new(model.ParamGetLastMessageID)
	if err := c.ShouldBind(m); err != nil {
		RespFailed(c, 400, consts.CodeShouldBind)
		return
	}
	if m.GroupID == 0 {
		RespFailed(c, 400, consts.CodeParamEmpty)
		return
	}
	lastID, err := service.GetGroupLastMessageID(m.GroupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			RespFailed(c, 400, consts.CodeGroupNotExist)
		} else {
			RespFailed(c, 500, consts.CodeDBCheckUser)
			global.Logger.Error("get last message failed", zap.Error(err))
		}
		return
	}
	RespSuccess(c, &model.ApiMessageID{MessageID: lastID})
}

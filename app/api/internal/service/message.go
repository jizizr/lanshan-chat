package service

import (
	"lanshan_chat/app/api/internal/consts"
	"lanshan_chat/app/api/internal/dao/mysql"
	"lanshan_chat/app/api/internal/dao/redis"
	"lanshan_chat/app/api/internal/model"
	"time"
)

func SendMessage(m *model.ParamSendMessage) (err error) {
	status, _, err := getUserSRInGroup(m.SenderID, m.GroupID)
	if err != nil {
		return
	}
	if status != "ok" {
		return consts.PermissionDeniedError
	}
	m.MessageID, err = redis.GetMessageID(m.GroupID)
	if err != nil {
		return
	}
	t := time.Now()
	m.SendDate = t
	_, err = mysql.SendMessageToGroup(m)
	return
}

func DeleteMessage(userID int64, deleteMsg *model.ParamDeleteMessage) error {
	id, senderID, err := mysql.QuerySenderIDFromGroup(deleteMsg.GroupID, deleteMsg.MessageID)
	if err != nil {
		return err
	}

	// 检验是否有权限删除消息
	if senderID != userID {
		_, role, err := getUserSRInGroup(userID, deleteMsg.GroupID)
		if err != nil {
			return err
		}
		if role != "admin" {
			return consts.PermissionDeniedError
		}
	}

	return mysql.DeleteMessageFromGroup(id)
}

func EditMessage(m *model.ParamSendMessage) error {
	id, senderID, err := mysql.QuerySenderIDFromGroup(m.GroupID, m.MessageID)
	if err != nil {
		return err
	}

	// 检验是否有权限修改消息
	if senderID != m.SenderID {
		return consts.PermissionDeniedError
	}
	m.ID = id
	m.SendDate = time.Now()
	return mysql.UpdateMessage(m)
}

func GetMessage(userID int64, m *model.ParamGetMessage) (messages []model.Message, err error) {
	status, _, err := getUserSRInGroup(userID, m.GroupID)
	if err != nil {
		return
	}
	if status == "banned" || status == "" {
		return nil, consts.PermissionDeniedError
	}
	messages, err = mysql.GetMessages(m.GroupID, m.StartID, m.Limit)
	return
}

package service

import (
	"lanshan_chat/app/api/internal/consts"
	"lanshan_chat/app/api/internal/dao/mysql"
	"lanshan_chat/app/api/internal/dao/redis"
	"lanshan_chat/app/api/internal/model"
	"time"
)

func GetGroupLastMessageID(groupID int64) (messageID int64, err error) {
	messageID, err = redis.GetMessageIDNow(groupID)
	if err == nil && messageID != 0 {
		return
	}
	return mysql.QueryLastMessageID(groupID)
}

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
	if err != nil {
		return
	}
	redis.SetMessageUniqueID(m.GroupID, m.MessageID, m.ID, m.SenderID)
	return
}

func getMessageISInGroup(groupID, messageID int64) (ID, senderID int64, err error) {
	ID, senderID = redis.GetMessageUniqueID(groupID, messageID)
	if ID != 0 && senderID != 0 {
		return
	}
	return mysql.QuerySenderIDFromGroup(groupID, messageID)
}

func deleteMessageFromGroup(ID, groupID, messageID int64) error {
	if err := mysql.DeleteMessageFromGroup(ID); err != nil {
		return err
	}
	return redis.DelMessage(groupID, messageID)
}

func DeleteMessage(userID int64, deleteMsg *model.ParamDeleteMessage) error {
	id, senderID, err := getMessageISInGroup(deleteMsg.GroupID, deleteMsg.MessageID)
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

	return deleteMessageFromGroup(id, deleteMsg.GroupID, deleteMsg.MessageID)
}

func EditMessage(m *model.ParamSendMessage) error {
	id, senderID, err := getMessageISInGroup(m.GroupID, m.MessageID)
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

func GetLastRead(userID, groupID int64) (lastRead int64, err error) {
	lastRead, err = redis.GetLastRead(groupID, userID)
	if err == nil && lastRead != 0 {
		return
	}
	return mysql.QueryLastRead(groupID, userID)
}

func ReadMessage(userID int64, m *model.ParamReadMessage) error {
	if err := mysql.UpdateLastRead(m.GroupID, userID, m.LastRead); err != nil {
		return err
	}
	return redis.SetLastRead(m.GroupID, userID, m.LastRead)
}

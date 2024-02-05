package mysql

import (
	"lanshan_chat/app/api/global"
	"lanshan_chat/app/api/internal/model"
)

const (
	SendMessageToGroupStr = `INSERT INTO group_message
    						(group_id, message_id, sender_id, reply_message_id, message, type, url, file_name, send_date)
    						VALUES (:group_id, :message_id, :sender_id, :reply_message_id, :message, :type, :url, :filename, :send_date)`
	DeleteMessageFromGroupStr = "DELETE FROM group_message WHERE id = ?"
	QuerySenderIDFromGroupStr = "SELECT id,sender_id FROM group_message WHERE group_id = ? AND message_id =? FOR UPDATE"
	UpdateMessageStr          = `UPDATE group_message 
								 SET reply_message_id = :reply_message_id,
								     message = :message,
								     type = :type,
								     url = :url,
								     file_name = :filename,
								     update_date = :send_date
								     WHERE id = :id`
	GetMessageStr = `SELECT group_id, message_id, sender_id, reply_message_id, message, type, url, file_name, send_date, update_date 
					 FROM group_message 
					 WHERE group_id = ? AND message_id > ?  
					 ORDER BY message_id LIMIT ?`
)

func SendMessageToGroup(param *model.ParamSendMessage) (id int64, err error) {
	result, err := global.MDB.NamedExec(SendMessageToGroupStr, param)
	if err != nil {
		return
	}
	id, err = result.LastInsertId()
	return
}

func QuerySenderIDFromGroup(groupID, messageID int64) (id int64, senderID int64, err error) {
	err = global.MDB.QueryRow(QuerySenderIDFromGroupStr, groupID, messageID).Scan(&id, &senderID)
	return
}

func DeleteMessageFromGroup(id int64) (err error) {
	_, err = global.MDB.Exec(DeleteMessageFromGroupStr, id)
	return
}

func UpdateMessage(param *model.ParamSendMessage) (err error) {
	_, err = global.MDB.NamedExec(UpdateMessageStr, param)
	return
}

func GetMessages(groupID, startID int64, limit int) (messages []model.Message, err error) {
	err = global.MDB.Select(&messages, GetMessageStr, groupID, startID, limit)
	return
}

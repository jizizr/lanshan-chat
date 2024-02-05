package model

import (
	"database/sql"
	"time"
)

type ParamSendMessage struct {
	ID             int64     `db:"id,omitempty"`
	GroupID        int64     `json:"group_id" form:"group_id" db:"group_id"`
	MessageID      int64     `json:"message_id" form:"message_id" db:"message_id"`
	SenderID       int64     `db:"sender_id"`
	ReplyMessageID int64     `json:"reply_message_id" form:"reply_message_id" db:"reply_message_id"`
	Message        string    `json:"message" form:"message" db:"message"`
	Type           string    `json:"type" form:"type" db:"type"`
	Url            string    `json:"url" form:"url" db:"url"`
	Filename       string    `json:"filename" form:"filename" db:"filename"`
	SendDate       time.Time `db:"send_date"`
}

type ParamDeleteMessage struct {
	GroupID   int64 `json:"group_id" form:"group_id" db:"group_id"`
	MessageID int64 `json:"message_id" form:"message_id" db:"message_id"`
}

type ParamGetMessage struct {
	GroupID int64 `json:"group_id" form:"group_id" db:"group_id"`
	StartID int64 `json:"start_id" form:"start_id" db:"start_id"`
	Limit   int   `json:"limit" form:"limit" db:"limit"`
}

type ParamReadMessage struct {
	GroupID  int64 `json:"group_id" form:"group_id" db:"group_id"`
	LastRead int64 `json:"last_read" form:"last_read" db:"last_read"`
}

type ParamGetLastMessageID struct {
	GroupID int64 `json:"group_id" form:"group_id" db:"group_id"`
}

type ApiMessageID struct {
	MessageID int64 `json:"message_id" form:"message_id" db:"message_id"`
}

type Message struct {
	GroupID        int64        `json:"group_id" db:"group_id"`
	MessageID      int64        `json:"message_id" db:"message_id"`
	SenderID       int64        `json:"sender_id" db:"sender_id"`
	ReplyMessageID int64        `json:"reply_message_id" db:"reply_message_id"`
	Message        string       `json:"message" db:"message"`
	Type           string       `json:"type" db:"type"`
	Url            string       `json:"url" db:"url"`
	FileName       string       `json:"filename" db:"file_name"`
	SendDate       time.Time    `json:"send_date" db:"send_date"`
	UpdateDate     sql.NullTime `json:"update_date" db:"update_date"`
}

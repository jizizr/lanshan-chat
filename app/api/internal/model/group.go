package model

import "time"

type ParamCreateGroup struct {
	GroupName   string `json:"group_name" form:"group_name"`
	Description string `json:"description" form:"description"`
	// ENUM [public,private]
	Type string `json:"type" form:"type"`
}

type Group struct {
	GroupID     int64     `json:"group_id" db:"group_id"`
	GroupName   string    `json:"group_name" db:"group_name"`
	Avatar      string    `json:"avatar" db:"avatar"`
	Description string    `json:"description" db:"description"`
	Type        string    `json:"type" db:"type"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

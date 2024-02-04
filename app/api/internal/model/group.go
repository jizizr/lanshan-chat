package model

import "time"

type ParamCreateGroup struct {
	GroupName   string `json:"group_name" form:"group_name"`
	Description string `json:"description" form:"description"`
	// ENUM [public,private]
	Type string `json:"type" form:"type"`
}

type ParamJoinPublicGroup struct {
	GroupID int64 `json:"group_id" form:"group_id"`
}

type ParamJoinPrivateGroup struct {
	Token string `json:"token" form:"token"`
}

type ParamGetPrivateGroupToken struct {
	GroupID     int64 `json:"group_id" form:"group_id"`
	ExpiresTime int64 `json:"expires_time" form:"expires_time"`
}

type ParamInviteUser struct {
	GroupID  int64 `json:"group_id" form:"group_id"`
	InviteID int64 `json:"invite_id" form:"invite_id"`
}

type ParamChangeMemberStatus struct {
	GroupID  int64  `json:"group_id" form:"group_id"`
	ChangeID int64  `json:"change_id" form:"change_id"`
	Status   string `json:"status" form:"status"`
}

type ParamKickFromGroup struct {
	GroupID int64 `json:"group_id" form:"group_id"`
	KickID  int64 `json:"kick_id" form:"kick_id"`
}

type ParamLeaveGroup struct {
	GroupID int64 `json:"group_id" form:"group_id"`
}

type PrivateGroupID struct {
	GroupID int64 `json:"group_id" form:"group_id"`
}

type Group struct {
	GroupID     int64     `json:"group_id" db:"group_id"`
	GroupName   string    `json:"group_name" db:"group_name"`
	Avatar      string    `json:"avatar" db:"avatar"`
	Description string    `json:"description" db:"description"`
	Type        string    `json:"type" db:"type"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

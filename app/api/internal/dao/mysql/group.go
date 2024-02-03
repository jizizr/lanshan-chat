package mysql

import (
	"lanshan_chat/app/api/global"
	"lanshan_chat/app/api/internal/model"
	"time"
)

const (
	CreateGroupStr = "INSERT INTO `groups` (group_name, avatar, description, type, created_at) VALUES (?, ?, ?, ?, ?)"
	JoinGroupStr   = "INSERT INTO `user_groups`(user_id, group_id, role, joined_at) VALUES (?, ?, ?, ?)"
)

func CreateGroup(g *model.ParamCreateGroup, url string, uid int64) (group *model.Group, err error) {
	tx, err := global.MDB.Beginx()
	if err != nil {
		return
	}
	t := time.Now()
	result, err := tx.Exec(CreateGroupStr, g.GroupName, url, g.Description, g.Type, t)
	if err != nil {
		_ = tx.Rollback()
		return
	}
	groupID, err := result.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return
	}
	_, err = tx.Exec(JoinGroupStr, uid, groupID, "admin", t)
	if err != nil {
		_ = tx.Rollback()
		return
	}
	group = &model.Group{
		GroupID:     groupID,
		GroupName:   g.GroupName,
		Avatar:      url,
		Description: g.Description,
		Type:        g.Type,
		CreatedAt:   t,
	}
	err = tx.Commit()
	return
}

package mysql

import (
	"database/sql"
	"errors"
	"lanshan_chat/app/api/global"
	"lanshan_chat/app/api/internal/model"
	"time"
)

const (
	CreateGroupStr        = "INSERT INTO `groups` (group_name, avatar, description, type, created_at) VALUES (?, ?, ?, ?, ?)"
	CheckGroupIsExistStr  = "SELECT type FROM `groups` WHERE group_id = ?"
	JoinGroupStr          = "INSERT INTO `user_groups`(user_id, group_id, role, joined_at,last_read) VALUES (?, ?, ?, ?, ?)"
	QueryUserSRInGroupStr = "SELECT status,role FROM `user_groups` WHERE group_id = ? AND user_id = ?"
	ChangeMemberStatusStr = "UPDATE `user_groups` SET status = ? WHERE group_id = ? AND user_id = ?"
	DeleteGroupUserStr    = "DELETE FROM `user_groups` WHERE group_id = ? AND user_id = ?"
	SearchGroupStr        = "SELECT group_id,group_name,avatar,description,type,created_at FROM `groups` WHERE type = 'public' AND (group_name LIKE ? OR description LIKE ?) LIMIT 10"
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

	if _, err = tx.Exec(JoinGroupStr, uid, groupID, "admin", t, 0); err != nil {
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

func QueryGroupType(groupID int64) (groupType string, err error) {
	err = global.MDB.Get(&groupType, CheckGroupIsExistStr, groupID)
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	return
}

func JoinGroup(uid, groupID, lastRead int64, role string, t time.Time) (err error) {
	_, err = global.MDB.Exec(JoinGroupStr, uid, groupID, role, t, lastRead)
	return
}

func QueryUserSRInGroup(uid, groupID int64) (status, role string, err error) {
	err = global.MDB.QueryRow(QueryUserSRInGroupStr, groupID, uid).Scan(&status, &role)
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	return
}

func ChangeMemberStatus(groupID, changeUserID int64, status string) (err error) {
	_, err = global.MDB.Exec(ChangeMemberStatusStr, status, groupID, changeUserID)
	return
}

func DeleteGroupUser(groupID, uid int64) (err error) {
	_, err = global.MDB.Exec(DeleteGroupUserStr, groupID, uid)
	return
}

func SearchGroup(keyword string) (groups []model.Group, err error) {
	groups = make([]model.Group, 10)
	err = global.MDB.Select(&groups, SearchGroupStr, "%"+keyword+"%", "%"+keyword+"%")
	return
}

package service

import (
	"lanshan_chat/app/api/internal/consts"
	"lanshan_chat/app/api/internal/dao/mysql"
	"lanshan_chat/app/api/internal/dao/redis"
	"lanshan_chat/app/api/internal/model"
	"lanshan_chat/utils"
	"time"
)

func CreateGroup(g *model.ParamCreateGroup, url string, uid int64) error {
	group, err := mysql.CreateGroup(g, url, uid)
	if err != nil {
		return err
	}
	redis.CreateGroup(group)
	redis.JoinGroup(uid, group.GroupID, 0, "admin", group.CreatedAt)
	_, err = redis.GetMessageID(group.GroupID)
	if err != nil {
		return err
	}
	return nil
}

func checkGroupType(groupID int64) (string, error) {
	t := redis.GetGroupType(groupID)
	if t != "" {
		return t, nil
	}
	return mysql.QueryGroupType(groupID)
}

func joinGroup(uid, groupID int64) (err error) {
	status, _, err := getUserSRInGroup(uid, groupID)
	if err != nil {
		return err
	}
	switch status {
	case "ok":
		return consts.GroupAlreadyJoinError
	case "banned":
		return consts.BanedError
	}

	t := time.Now()

	// lastRead 应该是此时群组最后一条消息的id
	lastRead, _ := GetGroupLastMessageID(groupID)
	if err := mysql.JoinGroup(uid, groupID, lastRead, "member", t); err != nil {
		return err
	}
	redis.JoinGroup(uid, groupID, lastRead, "member", t)
	return nil
}

func JoinPublicGroup(uid, groupID int64) error {
	// 检查群组是否是公开群组
	t, err := checkGroupType(groupID)
	if err != nil {
		return err
	}
	if t == "" {
		return consts.GroupNotExistError
	}
	if t == "private" {
		return consts.GroupIsPrivateError
	}
	return joinGroup(uid, groupID)
}

func JoinPrivateGroup(uid int64, token string) error {
	// 解析token
	claim, err := utils.ParseCustomToken(token)
	if err != nil {
		return consts.InvalidTokenError
	}
	return joinGroup(uid, int64(claim.ID.(map[string]interface{})["group_id"].(float64)))
}

// getUserSRInGroup 获取用户在群组中的状态和角色
func getUserSRInGroup(uid, groupID int64) (status, role string, err error) {
	status, role = redis.GetUserSRInGroup(uid, groupID)
	if status == "" {
		status, role, err = mysql.QueryUserSRInGroup(uid, groupID)
		if err != nil {
			return
		}
	}
	return
}

func GetPrivateGroupToken(uid, groupID, expiresTime int64) (string, error) {
	t, err := checkGroupType(groupID)
	if t == "public" {
		return "", consts.GroupIsPublicError
	}
	_, role, err := getUserSRInGroup(uid, groupID)
	if err != nil {
		return "", err
	}
	if role != "admin" {
		return "", consts.PermissionDeniedError
	}
	return utils.GenCustomToken(&model.PrivateGroupID{GroupID: groupID}, expiresTime)
}

// changeUserStatus 解除用户在群组中的禁言,需要提前检查用户是否有权限
func changeUserStatus(groupID, unbanUID int64, status string) error {
	err := mysql.ChangeMemberStatus(groupID, unbanUID, status)
	if err != nil {
		return err
	}
	return redis.ChangeMemberStatus(groupID, unbanUID, status)
}

func InviteToGroup(uid, groupID, inviteUID int64) error {
	t, err := checkGroupType(groupID)
	if t == "" {
		return consts.GroupNotExistError
	}
	// 检查用户是否有权限邀请
	_, role, err := getUserSRInGroup(uid, groupID)
	if err != nil {
		return err
	}
	switch role {
	case "member":
		if t == "private" {
			return consts.PermissionDeniedError
		}
		return joinGroup(inviteUID, groupID)
	case "admin":
		status, _, err := getUserSRInGroup(inviteUID, groupID)
		if err != nil {
			return err
		}
		switch status {
		case "ok":
			return consts.GroupAlreadyJoinError
		case "banned":
			err = changeUserStatus(groupID, inviteUID, "ok")
			if err != nil {
				return err
			}
		case "":
			return joinGroup(inviteUID, groupID)
		}
	case "":
		return consts.PermissionDeniedError
	}
	return consts.EnumError
}

func checkChangeRight(uid, groupID, changeUID int64) error {
	_, role, err := getUserSRInGroup(uid, groupID)
	if err != nil {
		return err
	}
	if role != "admin" {
		return consts.PermissionDeniedError
	}
	_, role, err = getUserSRInGroup(changeUID, groupID)
	if err != nil {
		return err
	}
	if role != "member" {
		return consts.PermissionDeniedError
	}
	return nil
}

func ChangeMemberStatus(uid, groupID, changeUID int64, status string) error {
	err := checkChangeRight(uid, groupID, changeUID)
	if err != nil {
		return err
	}
	return changeUserStatus(groupID, changeUID, status)
}

func delUser(uid, groupID int64) error {
	if err := mysql.DeleteGroupUser(groupID, uid); err != nil {
		return err
	}
	return redis.DelGroupUser(groupID, uid)
}

func KickFromGroup(uid, groupID, kickUID int64) error {
	if err := checkChangeRight(uid, groupID, kickUID); err != nil {
		return err
	}
	return delUser(kickUID, groupID)
}

func LeaveGroup(uid, groupID int64) error {
	status, _, err := getUserSRInGroup(uid, groupID)
	if err != nil {
		return err
	}
	switch status {
	case "":
		return consts.GroupNotExistError
	case "banned":
		return consts.BanedError
	}
	return delUser(uid, groupID)
}

func SearchGroup(keyword string) ([]model.Group, error) {
	return mysql.SearchGroup(keyword)
}

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

	// lastRead 应该是此时群组最后一条消息的id，这里先写为0
	if err := mysql.JoinGroup(uid, groupID, 0, "member", t); err != nil {
		return err
	}
	redis.JoinGroup(uid, groupID, 0, "member", t)
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

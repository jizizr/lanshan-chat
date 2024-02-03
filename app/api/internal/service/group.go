package service

import (
	"lanshan_chat/app/api/internal/dao/mysql"
	"lanshan_chat/app/api/internal/dao/redis"
	"lanshan_chat/app/api/internal/model"
)

func CreateGroup(g *model.ParamCreateGroup, url string, uid int64) error {
	group, err := mysql.CreateGroup(g, url, uid)
	if err != nil {
		return err
	}
	redis.CreateGroup(group)
	redis.JoinGroup(uid, group.GroupID, "admin", group.CreatedAt, 0)
	return nil
}

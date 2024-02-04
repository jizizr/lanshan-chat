package redis

import (
	"context"
	"fmt"
	"lanshan_chat/app/api/global"
	"lanshan_chat/app/api/internal/model"
	"time"
)

func CreateGroup(g *model.Group) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("group:%d", g.GroupID)
	field := map[string]interface{}{
		"group_name":  g.GroupName,
		"avatar":      g.Avatar,
		"description": g.Description,
		"type":        g.Type,
		"created_at":  g.CreatedAt.Unix(),
	}
	global.RDB.HMSet(ctx, key, field)
}

func GetGroupType(groupID int64) string {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("group:%d", groupID)
	return global.RDB.HGet(ctx, key, "type").Val()
}

func JoinGroup(uid, groupID, lastRead int64, role string, t time.Time) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("group:%d:%d", groupID, uid)
	field := map[string]interface{}{
		"last_read":   lastRead,
		"role":        role,
		"joined_at":   t.Unix(),
		"status":      "ok",
		"muted_until": nil,
	}
	global.RDB.HMSet(ctx, key, field)
}

func GetUserSRInGroup(uid, groupID int64) (status string, role string) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("group:%d:%d", groupID, uid)
	result := global.RDB.HMGet(ctx, key, "status", "role").Val()
	status, _ = result[0].(string)
	role, _ = result[1].(string)
	return
}

func ChangeMemberStatus(groupID, changeUserID int64, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("group:%d:%d", groupID, changeUserID)
	return global.RDB.HSet(ctx, key, "status", status).Err()
}

func DelGroupUser(groupID, kickID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("group:%d:%d", groupID, kickID)
	return global.RDB.Del(ctx, key).Err()
}

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

func JoinGroup(uid, groupID int64, role string, t time.Time, lastRead int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("group:%d:%d", groupID, uid)
	field := map[string]interface{}{
		"last_read":   lastRead,
		"role":        role,
		"joined_at":   t.Unix(),
		"status":      "ok",
		"baned_until": nil,
	}
	global.RDB.HMSet(ctx, key, field)
}

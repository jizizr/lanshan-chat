package redis

import (
	"context"
	"fmt"
	"lanshan_chat/app/api/global"
	"lanshan_chat/app/api/internal/consts"
	"time"
)

func CheckUserIsExist(uid int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("user:%d", uid)
	flag, err := global.RDB.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return flag == 1, nil
}

func AddUser(uid int64, username, nickname, password, email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("user:%d", uid)
	field := map[string]interface{}{
		"username": username,
		"nickname": nickname,
		"password": password,
		"email":    email,
		"profile":  consts.DefultProfile,
	}
	if err := global.RDB.HMSet(ctx, key, field).Err(); err != nil {
		return err
	}
	// 设置过期时间
	return global.RDB.Expire(ctx, key, 24*time.Hour).Err()
}

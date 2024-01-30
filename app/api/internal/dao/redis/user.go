package redis

import (
	"context"
	"fmt"
	"lanshan_chat/app/api/global"
	"time"
)

func CheckUserIsExist(username string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("user:%s", username)
	flag, err := global.RDB.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return flag == 1, nil
}

func AddUser(username, nickname, password, email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("user:%s", username)
	field := map[string]interface{}{
		"nickname": nickname,
		"password": password,
		"email":    email}
	if err := global.RDB.HMSet(ctx, key, field).Err(); err != nil {
		return err
	}
	// 设置过期时间
	if err := global.RDB.Expire(ctx, key, 2*time.Hour).Err(); err != nil {
		return err
	}
	return nil
}

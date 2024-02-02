package redis

import (
	"context"
	"fmt"
	"lanshan_chat/app/api/global"
	"lanshan_chat/app/api/internal/consts"
	"lanshan_chat/app/api/internal/model"
	"strconv"
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

func CheckUsername(username string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	return global.RDB.SIsMember(ctx, "users", username).Result()
}

func AddToSet(username string) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	global.RDB.SAdd(ctx, "users", username)
}
func AddUser(uid int64, username, nickname, password, email string, t time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()
	AddToSet(username)
	key := fmt.Sprintf("user:%d", uid)
	field := map[string]interface{}{
		"username":  username,
		"nickname":  nickname,
		"password":  password,
		"email":     email,
		"profile":   consts.DefultProfile,
		"joined_at": t.Unix(),
	}
	if err := global.RDB.HMSet(ctx, key, field).Err(); err != nil {
		return err
	}
	// 设置过期时间
	return global.RDB.Expire(ctx, key, 24*time.Hour).Err()
}

func GetUserProfile(uid int64) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("user:%d", uid)
	field, err := global.RDB.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	t, _ := strconv.ParseInt(field["joined_at"], 10, 64)
	return &model.User{
		Uid:      uid,
		Username: field["username"],
		Nickname: field["nickname"],
		Email:    field["email"],
		Profile:  field["profile"],
		JoinedAt: time.Unix(t, 0),
	}, nil
}

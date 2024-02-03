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

func DelFromSet(username string) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	global.RDB.SRem(ctx, "users", username)
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
	return global.RDB.HMSet(ctx, key, field).Err()
}

func GetUserInfo(uid int64) (*model.UserInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("user:%d", uid)
	field, err := global.RDB.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	t, _ := strconv.ParseInt(field["joined_at"], 10, 64)
	return &model.UserInfo{
		Uid:      uid,
		Username: field["username"],
		Nickname: field["nickname"],
		Email:    field["email"],
		Profile:  field["profile"],
		JoinedAt: time.Unix(t, 0),
	}, nil
}

func GetUser(uid int64) (*model.User, error) {
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
		Password: field["password"],
		JoinedAt: time.Unix(t, 0),
	}, nil

}

func ModifyUserInfo(u *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("user:%d", u.Uid)
	field := map[string]interface{}{
		"username":  u.Username,
		"nickname":  u.Nickname,
		"email":     u.Email,
		"password":  u.Password,
		"profile":   u.Profile,
		"joined_at": u.JoinedAt.Unix(),
	}
	return global.RDB.HMSet(ctx, key, field).Err()
}

package redis

import (
	"context"
	"fmt"
	"lanshan_chat/app/api/global"
	"time"
)

func CheckFriendIsExist(userID, friendID int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("group:%d:%d", userID, friendID)
	flag, err := global.RDB.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return flag == 1, nil
}

func AddFriend(userID, friendID int64, t time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("group:%d:%d", userID, friendID)
	field := map[string]interface{}{
		"create_at": t,
	}
	return global.RDB.HMSet(ctx, key, field).Err()
}

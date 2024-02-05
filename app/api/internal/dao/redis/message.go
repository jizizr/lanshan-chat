package redis

import (
	"context"
	"fmt"
	"lanshan_chat/app/api/global"
	"strconv"
	"time"
)

func GetMessageID(groupID int64) (msgID int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("group:msg:%d", groupID)
	msgID, err = global.RDB.Incr(ctx, key).Result()
	return
}

func GetMessageIDNow(groupID int64) (msgID int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("group:msg:%d", groupID)
	msgID, err = global.RDB.Get(ctx, key).Int64()
	return
}

func SetMessageUniqueID(groupID, msgID, ID, SenderID int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("id:%d:%d", groupID, msgID)
	field := map[string]int64{
		"ID":       ID,
		"SenderID": SenderID,
	}
	global.RDB.HMSet(ctx, key, field)
	global.RDB.Expire(ctx, key, 24*time.Hour)
}

func DelMessage(groupID, msgID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("id:%d:%d", groupID, msgID)
	return global.RDB.Del(ctx, key).Err()
}

func GetMessageUniqueID(groupID, msgID int64) (ID, senderID int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("id:%d:%d", groupID, msgID)
	field, err := global.RDB.HGetAll(ctx, key).Result()
	if err != nil || len(field) == 0 {
		return
	}
	global.RDB.Expire(ctx, key, 24*time.Hour)
	ID, _ = strconv.ParseInt(field["ID"], 10, 64)
	senderID, _ = strconv.ParseInt(field["SenderID"], 10, 64)
	return
}

func SetLastRead(groupID, userID, lastRead int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("read:%d:%d", groupID, userID)
	return global.RDB.Set(ctx, key, lastRead, 0).Err()
}

func GetLastRead(groupID, userID int64) (lastRead int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("read:%d:%d", groupID, userID)
	lastRead, err = global.RDB.Get(ctx, key).Int64()
	return
}

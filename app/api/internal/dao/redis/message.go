package redis

import (
	"context"
	"fmt"
	"lanshan_chat/app/api/global"
	"time"
)

func GetMessageID(groupID int64) (msgID int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := fmt.Sprintf("group:msg:%d", groupID)
	msgID, err = global.RDB.Incr(ctx, key).Result()
	return
}

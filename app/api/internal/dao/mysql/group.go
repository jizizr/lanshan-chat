package mysql

import (
	"lanshan_chat/app/api/global"
	"time"
)

const (
	CheckFriendIsExistStr = "SELECT EXISTS(SELECT 1 FROM friends WHERE user_id = ? AND friend_id = ?)"
	AddFriendStr          = "INSERT INTO friends (user_id, friend_id,created_at) VALUES (?, ?, ?)"
)

func CheckFriendIsExist(userID, friendID int64) (flag bool, err error) {
	err = global.MDB.Get(&flag, CheckFriendIsExistStr, userID, friendID)
	return
}

func AddFriend(userID, friendID int64) (t time.Time, err error) {
	t = time.Now()
	_, err = global.MDB.Exec(AddFriendStr, userID, friendID, t)
	return
}

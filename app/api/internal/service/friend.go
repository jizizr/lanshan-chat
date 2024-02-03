package service

import (
	"lanshan_chat/app/api/internal/consts"
	"lanshan_chat/app/api/internal/dao/mysql"
	"lanshan_chat/app/api/internal/dao/redis"
)

func AddFriend(userID, friendID int64) error {
	// 检查好友是否已经被添加
	flag, err := redis.CheckFriendIsExist(userID, friendID)
	if err != nil {
		return err
	}
	if !flag {
		flag, err = mysql.CheckFriendIsExist(userID, friendID)
		if err != nil {
			return err
		}
	}
	if flag {
		return consts.FriendExistError
	}
	// 添加好友
	t, err := mysql.AddFriend(userID, friendID)
	if err != nil {
		return err
	}
	return redis.AddFriend(userID, friendID, t)
}

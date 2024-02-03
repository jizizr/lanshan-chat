package consts

import "errors"

var (
	UserExistError        = errors.New("用户已存在")
	UserNotExistError     = errors.New("用户不存在")
	PasswordWrongError    = errors.New("密码错误")
	FriendExistError      = errors.New("好友已存在")
	FriendNotExistError   = errors.New("好友不存在")
	GroupAlreadyJoinError = errors.New("已经加入群组")
	BanedError            = errors.New("用户已被封禁")
	InvalidTokenError     = errors.New("无效的token")
	GroupNotExistError    = errors.New("群组不存在")
	GroupIsPublicError    = errors.New("群组是公开群")
	GroupIsPrivateError   = errors.New("群组是私有群")
	PermissionDeniedError = errors.New("权限不足")
)

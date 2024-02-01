package consts

import "errors"

var (
	UserExistError     = errors.New("用户已存在")
	UserNotExistError  = errors.New("用户不存在")
	PasswordWrongError = errors.New("密码错误")
)

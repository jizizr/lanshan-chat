package consts

type RespCode int16

const Success RespCode = 0

// service error
const (
	CodeShouldBind RespCode = 1000 + iota
	CodeUserAlreadyExist
	CodeUserNotExist
	CodeParamEmpty
	CodeEmailWrongFormat
	CodeUsernameWrongFormat
	CodeNeedLogin
	CodeInvalidToken
	CodeWrongPassword
	CodeServerBusy
	CodeFriendAlreadyExist
	CodeFriendNotExist
	CodeGroupAlreadyExist
	CodeGroupNotExist
	CodeFileTooLarge
	CodeFileEmpty
	CodeCompressFailed
	CodeGroupAlreadyJoin
	CodeUserInGroupBanned
	CodeGroupIsPublic
	CodeGroupIsPrivate
	CodePermissionDenied
)

// database error
const (
	CodeDBCheckUser RespCode = 2000 + iota
	CodeDBAddUser
	CodeDBCreateGroup
	CodeDBJoinGroup
)

var codeMsgMap = map[RespCode]string{
	Success:                 "success",
	CodeShouldBind:          "请求参数错误",
	CodeUserAlreadyExist:    "用户已存在",
	CodeUserNotExist:        "用户不存在",
	CodeDBCheckUser:         "数据库查询用户失败",
	CodeParamEmpty:          "参数不能为空",
	CodeEmailWrongFormat:    "邮箱格式错误",
	CodeUsernameWrongFormat: "用户名不能为邮箱",
	CodeNeedLogin:           "请先登录",
	CodeInvalidToken:        "无效的token",
	CodeWrongPassword:       "密码错误",
	CodeServerBusy:          "服务器繁忙",
	CodeDBAddUser:           "数据库添加用户失败",
	CodeFriendAlreadyExist:  "好友已存在",
	CodeFriendNotExist:      "好友不存在",
	CodeGroupAlreadyExist:   "群组已存在",
	CodeGroupNotExist:       "群组不存在",
	CodeFileTooLarge:        "文件过大",
	CodeFileEmpty:           "文件为空",
	CodeCompressFailed:      "压缩失败",
	CodeDBCreateGroup:       "数据库创建群组失败",
	CodeDBJoinGroup:         "数据库加入群组失败",
	CodeGroupAlreadyJoin:    "已经加入群组",
	CodeUserInGroupBanned:   "用户在群组中被封禁",
	CodeGroupIsPublic:       "群组是公开的",
	CodeGroupIsPrivate:      "群组是私有的",
	CodePermissionDenied:    "权限不足",
}

func (code RespCode) Msg() string {
	return codeMsgMap[code]
}

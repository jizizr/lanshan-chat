package consts

type RespCode int16

const Success RespCode = 0

// service error
const (
	CodeShouldBind RespCode = 1000 + iota
	CodeUserAlreadyExist
	CodeParamEmpty
	CodeEmailWrongFormat
	CodeNeedLogin
	CodeInvalidToken
)

// database error
const (
	CodeDBCheckUser RespCode = 2000 + iota
)

var codeMsgMap = map[RespCode]string{
	Success:              "success",
	CodeShouldBind:       "请求参数错误",
	CodeUserAlreadyExist: "用户已存在",
	CodeDBCheckUser:      "数据库查询用户失败",
	CodeParamEmpty:       "参数不能为空",
	CodeEmailWrongFormat: "邮箱格式错误",
	CodeNeedLogin:        "请先登录",
	CodeInvalidToken:     "无效的token",
}

func (code RespCode) Msg() string {
	return codeMsgMap[code]
}

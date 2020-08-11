package app

// Code 响应码
type Code int

// 主要Code
const (
	CodeSucceed           Code = 200
	CodeFailed            Code = 500
	CodeForbidden         Code = 403
	CodeNotFound          Code = 404
	CodeRequestParamError Code = 405
	CodeDataHasRelation   Code = 406
	CodeCustom            Code = 999 // 自定义错误
)

// 账号相关Code
const (
	CodeUserLoginSucceed Code = iota + 1000
	CodeUserNotFound
	CodeUserPasswordError
	CodeUserForbidden
)

// Token相关Code
const (
	CodeTokenExpired Code = iota + 2000
	CodeRefreshTokenExpired
	CodeRefreshTokenInvalidated
)

// 权限相关code
const (
	CodeNotAllowed Code = iota + 5000 // 没有权限
)

// codeMap Code对应信息
var codeMap = map[Code]string{
	CodeSucceed:           "成功",
	CodeFailed:            "失败",
	CodeForbidden:         "禁止操作",
	CodeNotFound:          "未找到或不存在",
	CodeRequestParamError: "请求参数错误",
	CodeDataHasRelation:   "数据存在关联，操作失败",

	CodeUserLoginSucceed:  "登录成功",
	CodeUserNotFound:      "账号不存在",
	CodeUserPasswordError: "密码错误",
	CodeUserForbidden:     "账号已禁用",

	CodeTokenExpired:            "当前会话已过期",
	CodeRefreshTokenExpired:     "刷新Token已过期，请重新登陆",
	CodeRefreshTokenInvalidated: "刷新Token不合法",

	CodeNotAllowed: "没有操作权限",
}

// 定义系统内部消息，模拟枚举类型，code = CodeCustom时返回msg
func (code Code) String(msg string) string {
	switch code {
	// 自定义错误
	case CodeCustom:
		return msg

	default:
		if msg, ok := codeMap[code]; ok {
			return msg
		}
		return "未知错误"
	}
}

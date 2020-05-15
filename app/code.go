package app

// 主要Code
const (
	CodeSucceed   = 200
	CodeFailed    = 500
	CodeForbidden = 403
	CodeNotFound  = 404
	CodeCustom    = 999 // 自定义错误
)

// 会员相关Code
const (
	CodeUserLoginSucceed Code = iota + 1000
	CodeUserNotFound
	CodeUserPasswordError
)

// Token相关Code
const (
	CodeTokenExpired Code = iota + 2000
)

// 权限相关code
const (
	CodeNotAllowed Code = iota + 5000 // 没有权限
)

// Code 响应码
type Code int

// 定义系统内部消息，模拟枚举类型，code = CodeCustom时返回msg
func (code Code) String(msg string) string {
	switch code {
	// 自定义错误
	case CodeCustom:
		return msg
	case CodeSucceed:
		return "成功"
	case CodeFailed:
		return "失败"
	case CodeForbidden:
		return "禁止操作"
	case CodeNotFound:
		return "信息不存在"

	case CodeUserLoginSucceed:
		return "登录成功"
	case CodeUserNotFound:
		return "账号不存在"
	case CodeUserPasswordError:
		return "密码错误"

	case CodeTokenExpired:
		return "当前会话已过期"

	case CodeNotAllowed:
		return "没有操作权限"

	default:
		return "未知错误"
	}
}

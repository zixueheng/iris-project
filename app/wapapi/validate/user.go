/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-18 09:34:00
 * @LastEditTime: 2021-04-28 14:50:49
 */
package validate

// UsernameLogin 会员用户名登录请求格式
type UsernameLogin struct {
	Username string `json:"username" validate:"required,gte=2,lte=50" comment:"用户名"`
	Password string `json:"password" validate:"required,gte=3" comment:"密码"`
}

// UsernameRegister 会员用户名注册请求格式
type UsernameRegister struct {
	Username string `json:"username" validate:"required,gte=2,lte=50" comment:"用户名"`
	Password string `json:"password" validate:"required,gte=3" comment:"密码"`
}

// UserphoneLogin 会员手机号登录请求格式
type UserphoneLogin struct {
	Phone    string `json:"phone" validate:"required,len=11" comment:"手机号"`
	Password string `json:"password" validate:"required,gte=3" comment:"密码"`
}

// UserphoneRegister 会员手机号注册请求格式
type UserphoneRegister struct {
	Phone      string `json:"phone" validate:"required,len=11" comment:"手机号"`
	Password   string `json:"password" validate:"required,gte=3" comment:"密码"`
	Verifycode string `json:"verifycode" validate:"required,gte=4" comment:"验证码"`
}

// UserphoneLoginRegister 会员手机号登录注册请求格式
type UserphoneLoginRegister struct {
	Phone      string `json:"phone" validate:"required,len=11" comment:"手机号"`
	Verifycode string `json:"verifycode" validate:"required,gte=4" comment:"验证码"`
}

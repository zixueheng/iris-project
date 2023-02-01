/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:27:34
 * @LastEditTime: 2021-04-28 14:46:35
 */
package middleware

import (
	"iris-project/app/config"

	"github.com/iris-contrib/middleware/jwt"
)

// JwtHandler JWT 中间件
func JwtHandler() *jwt.Middleware {
	var mySecret = []byte(config.App.Jwtsecret) // jwt验证密钥
	return jwt.New(jwt.Config{
		// CredentialsOptional: false,
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return mySecret, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		Expiration:    false, // 这里设置不主动验证token的过期时间，留在后续中间件中手动验证
	})
}

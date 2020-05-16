package middleware

import (
	"iris-project/config"

	"github.com/iris-contrib/middleware/jwt"
)

// JwtHandler JWT 中间件
func JwtHandler() *jwt.Middleware {
	var mySecret = []byte(config.App.Jwtsecret) // jwt验证密钥
	return jwt.New(jwt.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return mySecret, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		Expiration:    false, // 这里设置不主动验证token的过期时间，留在后续中间件中手动验证
	})
}

package middleware

import (
	"iris-project/config"

	"github.com/iris-contrib/middleware/jwt"
)

// JwtHandler JWT 中间件
func JwtHandler() *jwt.Middleware {
	var mySecret = []byte(config.App.Jwtsecret)
	return jwt.New(jwt.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return mySecret, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
}

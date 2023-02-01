/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-18 09:34:00
 * @LastEditTime: 2021-04-28 14:50:55
 */
package middleware

import (
	"errors"
	"iris-project/app"
	"iris-project/app/config"
	"iris-project/global"
	"iris-project/lib/util"
	"time"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

// AuthIfLogin WapApi端，如果登录获取UserID，不强制
// func AuthIfLogin(ctx iris.Context) {
// 	fmt.Println("验证1")
// 	value := ctx.Values().Get("jwt").(*jwt.Token)
// 	data := value.Claims.(jwt.MapClaims)

// 	var useID, exp string
// 	if value, ok := data[global.WapUserJWTKey]; ok {
// 		useID = value.(string)
// 	} else {
// 		ctx.Next()
// 		return
// 	}

// 	if value, ok := data["exp"]; ok {
// 		exp = value.(string)
// 	} else {
// 		ctx.Next()
// 		return
// 	}

// 	expObj, err := time.ParseInLocation(config.App.Timeformat, exp, time.Local)
// 	if err != nil { // 过期时间解析错误，返回 BadRequest
// 		ctx.Next()
// 		return
// 	}

// 	if expObj.Before(time.Now()) { // Token 超时
// 		ctx.Next()
// 		return
// 	}

// 	if expObj.Before(time.Now()) { // Token 超时
// 		ctx.Next()
// 		return
// 	}

// 	ctx.Values().Set(global.WapUserJWTKey, uint32(util.ParseInt(useID))) // 将 userID 存储到 ctx 中 以共享
// 	ctx.Next()
// }

// Auth WapApi端，强制Token验证
func Auth(ctx iris.Context) {
	value := ctx.Values().Get("jwt").(*jwt.Token)
	data := value.Claims.(jwt.MapClaims)
	// for key, value := range data {
	// 	ctx.Writef("%s = %s\n", key, value)
	// }

	var userID, exp string
	if value, ok := data[global.WapUserJWTKey]; ok {
		userID = value.(string)
	} else {
		app.ResponseProblemHTTPCode(ctx, iris.StatusBadRequest, errors.New("Token中没有"+global.WapUserJWTKey))
		return
	}

	if value, ok := data["exp"]; ok {
		exp = value.(string)
	} else {
		app.ResponseProblemHTTPCode(ctx, iris.StatusBadRequest, errors.New("Token中没有exp"))
		return
	}

	expObj, err := time.ParseInLocation(config.App.Timeformat, exp, time.Local)
	if err != nil { // 过期时间解析错误，返回 BadRequest
		app.ResponseProblemHTTPCode(ctx, iris.StatusBadRequest, err)
		return
	}

	if expObj.Before(time.Now()) { // Token 超时
		ctx.JSON(app.APIData(false, app.CodeTokenExpired, "", nil))
		ctx.StopExecution()
		return
	}
	ctx.Values().Set(global.WapUserJWTKey, uint32(util.ParseInt(userID))) // 将 userID 存储到 ctx 中 以共享
	ctx.Next()
}

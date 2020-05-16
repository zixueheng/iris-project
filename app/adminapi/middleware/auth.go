package middleware

import (
	"iris-project/app"
	"iris-project/config"
	"iris-project/lib/util"
	"time"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

// Auth Token验证
func Auth(ctx iris.Context) {
	value := ctx.Values().Get("jwt").(*jwt.Token)

	data := value.Claims.(jwt.MapClaims)
	// for key, value := range data {
	// 	ctx.Writef("%s = %s\n", key, value)
	// }
	adminUseID := data["admin_user_id"].(string)
	exp := data["exp"].(string)

	// fmt.Println(adminUseID, exp)

	expObj, err := time.ParseInLocation(config.App.Timeformat, exp, time.Local)
	if err != nil { // 过期时间解析错误，返回 BadRequest
		ctx.Application().Logger().Error(err)
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.StopExecution()
		return
	}
	// fmt.Println(util.TimeFormat(expObj, ""), util.TimeFormat(time.Now(), ""))

	if expObj.Before(time.Now()) { // Token 超时
		ctx.JSON(app.APIData(false, app.CodeTokenExpired, "", nil))
		ctx.StopExecution()
		return
	}
	ctx.Values().Set("auth_admin_user_id", util.ParseInt(adminUseID, 0)) // 将 admin_user_id 存储到 ctx 中 以共享
	ctx.Next()
}

/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:28:08
 * @LastEditTime: 2022-11-15 14:34:56
 */
package middleware

import (
	"encoding/json"
	"errors"
	"iris-project/app"
	"iris-project/app/adminapi/model"
	"iris-project/app/config"
	"iris-project/global"
	"iris-project/lib/cache"
	"iris-project/lib/util"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

// Auth AdminApi Token验证
func Auth(ctx iris.Context) {
	value := ctx.Values().Get("jwt").(*jwt.Token)
	data := value.Claims.(jwt.MapClaims)
	// for key, value := range data {
	// 	ctx.Writef("%s = %s\n", key, value)
	// }

	var adminUseID, exp string
	if value, ok := data[global.AdminUserJWTKey]; ok {
		adminUseID = value.(string)
	} else {
		app.ResponseProblemHTTPCode(ctx, iris.StatusBadRequest, errors.New("Token中没有"+global.AdminUserJWTKey))
	}

	if value, ok := data["exp"]; ok {
		exp = value.(string)
	} else {
		app.ResponseProblemHTTPCode(ctx, iris.StatusBadRequest, errors.New("Token中没有exp"))
	}

	expObj, err := time.ParseInLocation(config.App.Timeformat, exp, time.Local)
	if err != nil { // 过期时间解析错误，返回 BadRequest
		app.ResponseProblemHTTPCode(ctx, iris.StatusBadRequest, err)
	}

	if expObj.Before(time.Now()) { // Token 超时
		ctx.JSON(app.APIData(false, app.CodeTokenExpired, "", nil))
		ctx.StopExecution()
		return
	}

	var adminUser = new(model.AdminUser)
	cacheAdminUser, err := cache.Get(config.App.Appname + global.AdminUserCacheKeyPrefix + adminUseID) // 加载缓存中账号信息
	// fmt.Println("缓存", cacheAdminUser)
	if err == redis.Nil {
		adminUser.ID = uint32(util.ParseInt(adminUseID))
		if !adminUser.GetOne(true) {
			ctx.JSON(app.APIData(false, app.CodeUserNotFound, "", nil)) // 账号不存在
			ctx.StopExecution()
			return
		}

		if len(adminUser.Roles) == 0 { // 账号没有关联角色，即没有权限
			ctx.JSON(app.APIData(false, app.CodeNotAllowed, "", nil))
			ctx.StopExecution()
			return
		}
		json, _ := json.Marshal(adminUser)
		cache.Set(config.App.Appname+global.AdminUserCacheKeyPrefix+adminUseID, string(json), time.Minute*time.Duration(global.AdminUserCacheMinutes)) // 账号信息保存到缓存
		if !adminUser.SuperAdmin {                                                                                                                     // 不是超级管理员 检查权限
			if !checkRight(adminUser, ctx.GetCurrentRoute().ResolvePath(), ctx.GetCurrentRoute().Method()) {
				ctx.JSON(app.APIData(false, app.CodeNotAllowed, "", nil))
				ctx.StopExecution()
				return
			}
		}

	} else {
		json.Unmarshal([]byte(cacheAdminUser), adminUser)
		if !adminUser.SuperAdmin { // 不是超级管理员 检查权限
			if !checkRight(adminUser, ctx.GetCurrentRoute().ResolvePath(), ctx.GetCurrentRoute().Method()) {
				ctx.JSON(app.APIData(false, app.CodeNotAllowed, "", nil))
				ctx.StopExecution()
				return
			}
		}
	}
	ctx.Values().Set("auth_admin_user", adminUser) // 将 admin_user 存储到 ctx 中 以共享
	ctx.Next()
}

// checkRight 检查权限
func checkRight(adminUser *model.AdminUser, path, method string) (hasRight bool) {
	// return true

	// fmt.Println("检查权限：", adminUser.Username, path, method)

	hasRight = false
	menus := adminUser.Menus
	for _, menu := range menus {
		if menu.Type == "api" && menu.Status == 1 && menu.APIPath == path && strings.ToUpper(menu.Method) == strings.ToUpper(method) {
			hasRight = true
			return
		}
	}
	return

}

package middleware

import (
	"encoding/json"
	"errors"
	"iris-project/app"
	"iris-project/app/adminapi/model"
	"iris-project/config"
	"iris-project/global"
	"iris-project/lib/util"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
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

	var adminUseID, exp string
	if value, ok := data["admin_user_id"]; ok {
		adminUseID = value.(string)
	} else {
		app.ResponseProblemHTTPCode(ctx, iris.StatusBadRequest, errors.New("Token中没有admin_user_id"))
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
	cacheAdminUser, err := global.Redis.Get("vo_admin_user_" + adminUseID).Result() // 加载redis中账号信息
	if err == redis.Nil {
		// fmt.Println("走数据库")
		if !adminUser.GetAdminUserByID(util.ParseInt(adminUseID)) {
			ctx.JSON(app.APIData(false, app.CodeUserNotFound, "", nil)) // 账号不存在
			ctx.StopExecution()
			return
		}

		var role = new(model.Role)
		if role.GetRoleByID(adminUser.RoleID) { // 加载角色包括角色的菜单
			adminUser.Role = *role
			json, _ := json.Marshal(adminUser)
			global.Redis.Set("vo_admin_user_"+adminUseID, string(json), time.Minute*time.Duration(global.AdminUserCacheMinutes)) // 账号信息保存到redis

			if adminUser.Role.Tag != "superadmin" { // 不是超级管理员 检查权限
				if !checkRight(adminUser, ctx.GetCurrentRoute().ResolvePath(), ctx.GetCurrentRoute().Method()) {
					ctx.JSON(app.APIData(false, app.CodeNotAllowed, "", nil))
					ctx.StopExecution()
					return
				}
			}
		} else { // 账号没有关联角色，即没有权限
			ctx.JSON(app.APIData(false, app.CodeNotAllowed, "", nil))
			ctx.StopExecution()
			return
		}
	} else {
		// fmt.Println("走缓存")
		json.Unmarshal([]byte(cacheAdminUser), adminUser)
		if adminUser.Role.Tag != "superadmin" { // 不是超级管理员 检查权限
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
	// fmt.Println("检查权限：", adminUser.Username, path, method)
	hasRight = false
	menus := adminUser.Role.Menus
	for _, menu := range menus {
		if menu.APIPath == path && strings.ToUpper(menu.Method) == strings.ToUpper(method) {
			hasRight = true
			return
		}
	}

	return
}

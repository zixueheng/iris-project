package controller

import (
	"errors"
	"iris-project/app"
	"iris-project/app/adminapi/model"

	"github.com/kataras/iris/v12"
)

// GetAuthAdminUser 获取AuthAdminUser
func GetAuthAdminUser(ctx iris.Context) *model.AdminUser {
	entry, ok := ctx.Values().GetEntry("auth_admin_user")
	if !ok {
		app.ResponseProblemHTTPCode(ctx, iris.StatusInternalServerError, errors.New("未取到 auth中间件 设置的 auth_admin_user"))
	}
	return entry.ValueRaw.(*model.AdminUser)
}

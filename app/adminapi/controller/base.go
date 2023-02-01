/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:28:08
 * @LastEditTime: 2021-04-28 14:57:42
 */
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

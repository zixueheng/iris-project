/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-18 09:34:00
 * @LastEditTime: 2021-04-28 14:51:38
 */
package controller

import (
	"iris-project/global"

	"github.com/kataras/iris/v12"
)

// GetAuthUserID 获取登录会员ID
func GetAuthUserID(ctx iris.Context) uint32 {
	entry, ok := ctx.Values().GetEntry(global.WapUserJWTKey)
	if !ok {
		// app.ResponseProblemHTTPCode(ctx, iris.StatusInternalServerError, errors.New("未取到 auth中间件 设置的"+global.WapUserJWTKey))
		return 0
	}
	return entry.ValueRaw.(uint32)
}

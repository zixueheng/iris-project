/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:27:34
 * @LastEditTime: 2023-03-02 10:16:23
 */
package routes

import (
	"github.com/kataras/iris/v12"
	// "github.com/kataras/iris/v12/middleware/requestid"
)

// InitRoute 加载路由
func InitRoute(app *iris.Application) {
	root := app.Party("/")
	// root.Use(requestid.New())
	InitAdminAPI(root) // 加载 adminapi 模块路由
	InitWapAPI(root)   // 加载 wapapi 模块 路由

	// InitWebSocket(root) // 加载websocket
}

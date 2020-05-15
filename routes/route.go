package routes

import "github.com/kataras/iris/v12"

// InitRoute 加载路由
func InitRoute(app *iris.Application) {
	root := app.Party("/")
	InitAdminAPI(root) // 加载 adminapi 模块路由
	InitWapAPI(root)   // 加载 wapapi 模块 路由
}

package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"

	"iris-project/app/wapapi/controller"
	"iris-project/middleware"
)

// InitWapAPI 初始化 WapApi 模块路由
func InitWapAPI(app iris.Party) {
	adminParty := app.Party("/wapapi", middleware.CrsAuth()).AllowMethods(iris.MethodOptions)
	mvc.Configure(adminParty, loadWapAPIController)
}

// loadWapAPIController 加载 wapapi 控制器
func loadWapAPIController(app *mvc.Application) {
	app.Handle(new(controller.User))
	app.Handle(new(controller.Car))
}

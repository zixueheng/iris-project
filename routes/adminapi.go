package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"

	"iris-project/app/adminapi/controller"
	adminapimiddleware "iris-project/app/adminapi/middleware"
	"iris-project/middleware"
)

// InitAdminAPI 初始化 AdminApi 模块路由
func InitAdminAPI(app iris.Party) {
	adminParty := app.Party("/adminapi", middleware.CrsAuth()).AllowMethods(iris.MethodOptions)
	// {
	// 	adminParty.Get("/login", controller.Login)
	// }

	mvc.Configure(adminParty, loadAdminAPIController)
}

// loadAdminAPIController 加载 adminapi 控制器
func loadAdminAPIController(app *mvc.Application) {
	app.Handle(new(controller.Public))

	// 以下控制器路由都需要验证 token
	app.Router.Use(middleware.JwtHandler().Serve, adminapimiddleware.Auth)
	app.Handle(new(controller.User))
	app.Handle(new(controller.Car))
}

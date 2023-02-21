/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:27:34
 * @LastEditTime: 2023-02-21 09:48:35
 */
package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"

	"iris-project/app/adminapi/controller"
	adminapimiddleware "iris-project/app/adminapi/middleware"
	"iris-project/global"
	"iris-project/middleware"
)

// InitAdminAPI 初始化 AdminApi 模块路由
func InitAdminAPI(app iris.Party) {
	party := app.Party(global.AdminAPI, middleware.CrsAuth()).AllowMethods(iris.MethodOptions)
	// {
	// 	adminParty.Get("/login", controller.Login)
	// }

	mvc.Configure(party, loadAdminAPIController)
}

// loadAdminAPIController 加载 adminapi 控制器
func loadAdminAPIController(app *mvc.Application) {
	app.Handle(new(controller.Public))
	app.Handle(new(controller.Init))

	// 以下控制器路由都需要验证 token
	app.Router.Use(middleware.JwtHandler().Serve, adminapimiddleware.Auth)
	app.Handle(new(controller.AdminUser))
	app.Handle(new(controller.Role))
	app.Handle(new(controller.Menu))

	app.Handle(new(controller.Home))
	app.Handle(new(controller.System))

	app.Handle(new(controller.File))
	app.Handle(new(controller.Config))
}

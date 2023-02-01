/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:27:34
 * @LastEditTime: 2022-09-29 10:59:47
 */
package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"

	"iris-project/app/wapapi/controller"
	wapapimiddleware "iris-project/app/wapapi/middleware"
	"iris-project/middleware"
)

// InitWapAPI 初始化 WapApi 模块路由
func InitWapAPI(app iris.Party) {
	party := app.Party("/wapapi", middleware.CrsAuth()).AllowMethods(iris.MethodOptions)
	mvc.Configure(party, loadWapAPIController)
}

// loadWapAPIController 加载 wapapi 控制器
func loadWapAPIController(app *mvc.Application) {
	app.Handle(new(controller.Public))
	app.Handle(new(controller.Version))
	app.Handle(new(controller.Home))

	app.Router.Use(middleware.JwtHandler().Serve, wapapimiddleware.Auth)
	app.Handle(new(controller.Address))
	app.Handle(new(controller.User))

}
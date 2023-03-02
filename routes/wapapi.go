/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:27:34
 * @LastEditTime: 2023-02-27 10:03:27
 */
package routes

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"

	"iris-project/app/wapapi/controller"
	wapapimiddleware "iris-project/app/wapapi/middleware"
	"iris-project/global"
	"iris-project/middleware"
)

// InitWapAPI 初始化 WapApi 模块路由
func InitWapAPI(app iris.Party) {
	party := app.Party(global.WapAPI, middleware.CrsAuth(), middleware.Sentinel()).AllowMethods(iris.MethodOptions)
	// party.UseRouter(middleware.AcLog.Handler) // 使用httplog
	mvc.Configure(party, loadWapAPIController)
}

// loadWapAPIController 加载 wapapi 控制器
func loadWapAPIController(app *mvc.Application) {
	app.Handle(new(controller.Public))
	app.Handle(new(controller.Version))
	app.Handle(new(controller.Home))
	app.Handle(new(controller.Test))

	app.Router.Use(middleware.JwtHandler().Serve, wapapimiddleware.Auth)
	app.Handle(new(controller.Address))
	app.Handle(new(controller.User))

}

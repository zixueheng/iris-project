package controller

import (
	"fmt"
	"iris-project/app"
	"iris-project/app/adminapi/model"
	"iris-project/app/adminapi/validate"
	"iris-project/global"
	"iris-project/middleware"

	"github.com/go-redis/redis/v7"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// Public 公共控制器，不需要验证token
type Public struct {
	Ctx iris.Context // IRIS框架会自动注入 Context
}

// GetInit 初始化数据
// func (p *Public) GetInit() string {
// 	adminUser := model.AdminUser{
// 		Username: "admin",
// 		Password: util.HashPassword("123456"),
// 		Role: model.Role{
// 			Name:   "超级管理员",
// 			Tag:    "superadmin",
// 			Status: 1,
// 		},
// 		// RoleID:   role.ID,
// 		Phone:  "15215657185",
// 		Status: 1,
// 	}
// 	global.Db.Create(&adminUser)

// 	return "ok"
// }

// AfterActivation 前置方法
func (p *Public) AfterActivation(a mvc.AfterActivation) {
	// 给单独的控制器方法添加中间件
	// select the route based on the method name you want to modify.
	index := a.GetRoute("PostRefreshtoken") // 根据 方法名 获取 方法的路由
	// just prepend the handler(s) as middleware(s) you want to use. or append for "done" handlers.
	index.Handlers = append([]iris.Handler{middleware.JwtHandler().Serve}, index.Handlers...) // 将中间件 追加到 路由的 Handleers 字段中
}

// PostLogin 登录
func (p *Public) PostLogin() {
	loginInfo := new(validate.LoginRequest)

	errmsg := app.CheckRequest(p.Ctx, loginInfo)
	if len(errmsg) != 0 {
		p.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}

	adminUser := &model.AdminUser{}
	response, ok, code := adminUser.CheckLogin(loginInfo)

	p.Ctx.JSON(app.APIData(ok, code, "", response))
}

// PostRefreshtoken 刷新缓存
func (p *Public) PostRefreshtoken() {
	value := p.Ctx.Values().Get("jwt").(*jwt.Token) // 这里要先给这个路由方法添加JWT的中间件才能获取到 jwt变量

	data := value.Claims.(jwt.MapClaims)
	// for key, value := range data {
	// 	ctx.Writef("%s = %s\n", key, value)
	// }
	adminUseID := data["admin_user_id"].(string)

	refreshToken, err := global.Redis.Get("refresh_token_admin_" + adminUseID).Result()
	if err == redis.Nil {
		fmt.Println("refresh token not existing")
		p.Ctx.JSON(app.APIData(false, app.CodeRefreshTokenExpired, "", nil))
	} else {
		p.Ctx.Writef("ID: %s, RefreshToken: %s", adminUseID, refreshToken)
	}
}

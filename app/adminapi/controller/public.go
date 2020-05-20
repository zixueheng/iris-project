package controller

import (
	"iris-project/app"
	"iris-project/app/adminapi/model"
	"iris-project/app/adminapi/validate"
	"iris-project/global"
	"iris-project/lib/util"
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
func (p *Public) GetInit() string {
	// adminUser := model.AdminUser{
	// 	Username: "admin",
	// 	Password: util.HashPassword("123456"),
	// 	Role: model.Role{
	// 		Name:   "超级管理员",
	// 		Tag:    "superadmin",
	// 		Status: 1,
	// 	},
	// 	// RoleID:   role.ID,
	// 	Phone:  "15215657185",
	// 	Status: 1,
	// }
	// global.Db.Create(&adminUser)

	// role := model.Role{
	// 	Name: "商品管理员",
	// 	Tag:  "goods_manager",
	// 	Menus: []model.Menu{
	// 		{ID: 1, PID: 0, Name: "商品列表", Type: "menu", APIPath: "/adminapi/goodslist/%v/%v", Method: "GET", Sort: 1, Status: 1},
	// 		{ID: 2, PID: 1, Name: "商品详情", Type: "api", APIPath: "/adminapi/goods/%v", Method: "GET", Sort: 0, Status: 1},
	// 		{ID: 3, PID: 1, Name: "商品编辑", Type: "api", APIPath: "/adminapi/goods/%v", Method: "POST", Sort: 0, Status: 1},
	// 		{ID: 4, PID: 1, Name: "商品删除", Type: "api", APIPath: "/adminapi/goods/%v", Method: "DELETE", Sort: 0, Status: 1},
	// 		{ID: 5, PID: 0, Name: "商品分类", Type: "menu", APIPath: "/adminapi/categorylist/%v/%v", Method: "GET", Sort: 2, Status: 1},
	// 		{ID: 6, PID: 5, Name: "商品分类详情", Type: "api", APIPath: "/adminapi/category/%v", Method: "GET", Sort: 0, Status: 1},
	// 		{ID: 7, PID: 5, Name: "商品分类编辑", Type: "api", APIPath: "/adminapi/category/%v", Method: "POST", Sort: 0, Status: 1},
	// 		{ID: 8, PID: 5, Name: "商品分类删除", Type: "api", APIPath: "/adminapi/category/%v", Method: "DELETE", Sort: 0, Status: 1},
	// 	},
	// 	Status: 1,
	// }
	// global.Db.Create(&role)

	// goodseditor := model.AdminUser{
	// 	Username: "goodseditor",
	// 	Password: util.HashPassword("123456"),
	// 	// Role: model.Role{
	// 	// 	Model: gorm.Model{ID:2},
	// 	// },
	// 	RoleID: 2,
	// 	Phone:  "13721047437",
	// 	Status: 1,
	// }
	// global.Db.Create(&goodseditor)

	// admin1 := new(model.AdminUser)
	// admin1.ID = 1
	// global.Db.Model(admin1).Update("phone", "16666666666")

	return "ok"
}

// AfterActivation 前置方法
func (p *Public) AfterActivation(a mvc.AfterActivation) {
	// 给单独的控制器方法添加中间件
	// select the route based on the method name you want to modify.
	refreshtokenRoute := a.GetRoute("PostRefreshtoken") // 根据 方法名 获取 方法的路由
	// just prepend the handler(s) as middleware(s) you want to use. or append for "done" handlers.
	refreshtokenRoute.Handlers = append([]iris.Handler{middleware.JwtHandler().Serve}, refreshtokenRoute.Handlers...) // 将中间件 追加到 路由的 Handleers 字段中
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

	param := struct {
		RefreshToken string `json:"refresh_token"`
	}{}

	if err := p.Ctx.ReadJSON(&param); err != nil {
		p.Ctx.JSON(app.APIData(false, app.CodeRequestParamError, "", nil))
		return
	}
	if param.RefreshToken == "" {
		p.Ctx.JSON(app.APIData(false, app.CodeRequestParamError, "", nil))
		return
	}

	refreshToken, err := global.Redis.Get("refresh_token_admin_" + adminUseID).Result()
	if err == redis.Nil {
		p.Ctx.JSON(app.APIData(false, app.CodeRefreshTokenExpired, "", nil))
	} else if param.RefreshToken != refreshToken {
		p.Ctx.JSON(app.APIData(false, app.CodeRefreshTokenInvalidated, "", nil))
	} else {
		token, refreshToken := app.GenTokenAndRefreshToken("admin_user_id", util.ParseInt(adminUseID), global.AdminTokenMinutes, global.AdminRefreshTokenMinutes)
		response := struct {
			Token        string `json:"token"`
			RefreshToken string `json:"refresh_token"`
		}{
			Token:        token,
			RefreshToken: refreshToken,
		}
		p.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", response))
	}
}

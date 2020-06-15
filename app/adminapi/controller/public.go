package controller

import (
	"errors"
	"iris-project/app"
	"iris-project/app/adminapi/model"
	"iris-project/app/adminapi/validate"
	"iris-project/config"
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
	global.Db.AutoMigrate(
		// &model.User{},
		&model.Menu{},
		&model.Role{},
		&model.AdminUser{},
		&model.AdminUserRole{},
		&model.RoleMenu{},
	)
	// AutoMigrate 会忽略外键，需手动添加（建议直接到数据库中添加）
	// 参数分别为模型外键，关联表主键，删除级联，修改级联
	global.Db.Model(&model.AdminUserRole{}).AddForeignKey("admin_user_id", "iris_admin_user(id)", "CASCADE", "CASCADE")
	global.Db.Model(&model.AdminUserRole{}).AddForeignKey("role_id", "iris_role(id)", "CASCADE", "CASCADE")

	global.Db.Model(&model.RoleMenu{}).AddForeignKey("role_id", "iris_role(id)", "CASCADE", "CASCADE")
	global.Db.Model(&model.RoleMenu{}).AddForeignKey("menu_id", "iris_menu(id)", "CASCADE", "CASCADE")

	adminUser := model.AdminUser{
		Username: "admin",
		Password: util.HashPassword("123456"),
		Roles: []model.Role{
			{
				Name:   "超级管理员",
				Tag:    "superadmin",
				Status: 1,
			},
		},
		Phone:  "15215657185",
		Status: 1,
	}
	global.Db.Create(&adminUser)

	menus := []*model.Menu{
		{ID: 1, PID: 0, Name: "主页", Icon: "md-home", Type: "menu", MenuPath: "/admin/home/", APIPath: "", Method: "", UniqueAuthKey: "admin-home", Sort: 1, Status: 1},
		{ID: 2, PID: 1, Name: "首页统计", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/statistic", Method: "GET", UniqueAuthKey: "admin-home-statistic", Sort: 0, Status: 1},

		{ID: 3, PID: 0, Name: "权限管理", Icon: "md-settings", Type: "menu", MenuPath: "/admin/setting/", APIPath: "", Method: "", UniqueAuthKey: "admin-setting", Sort: 100, Status: 1},

		{ID: 4, PID: 3, Name: "管理员", Icon: "", Type: "menu", MenuPath: "/admin/setting/adminuser", APIPath: "", Method: "", UniqueAuthKey: "admin-setting-adminuser", Sort: 0, Status: 1},
		{ID: 5, PID: 4, Name: "管理员列表", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/adminuser/list/%v/%v", Method: "GET", UniqueAuthKey: "admin-setting-adminuser-list", Sort: 0, Status: 1},
		{ID: 6, PID: 4, Name: "管理员详情", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/adminuser/%v", Method: "GET", UniqueAuthKey: "admin-setting-adminuser-info", Sort: 0, Status: 1},
		{ID: 7, PID: 4, Name: "管理员添加编辑", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/adminuser", Method: "POST", UniqueAuthKey: "admin-setting-adminuser-save", Sort: 0, Status: 1},
		{ID: 8, PID: 4, Name: "管理员删除", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/adminuser/%v", Method: "DELETE", UniqueAuthKey: "admin-setting-adminuser-delete", Sort: 0, Status: 1},
		{ID: 9, PID: 4, Name: "管理员禁用启用", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/adminuser/status/%v", Method: "GET", UniqueAuthKey: "admin-setting-adminuser-status", Sort: 0, Status: 1},

		{ID: 10, PID: 3, Name: "角色", Icon: "", Type: "menu", MenuPath: "/admin/setting/role", APIPath: "", Method: "", UniqueAuthKey: "admin-setting-role", Sort: 0, Status: 1},
		{ID: 11, PID: 10, Name: "角色列表", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/role/list/%v/%v", Method: "GET", UniqueAuthKey: "admin-setting-role-list", Sort: 0, Status: 1},
		{ID: 12, PID: 10, Name: "角色添加编辑", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/role", Method: "POST", UniqueAuthKey: "admin-setting-role-save", Sort: 0, Status: 1},
		{ID: 13, PID: 10, Name: "角色详情", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/role/%v", Method: "GET", UniqueAuthKey: "admin-setting-role-info", Sort: 0, Status: 1},
		{ID: 14, PID: 10, Name: "角色删除", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/role/%v", Method: "DELETE", UniqueAuthKey: "admin-setting-role-delete", Sort: 0, Status: 1},
		{ID: 15, PID: 10, Name: "角色禁用或启用", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/role/status/%v", Method: "GET", UniqueAuthKey: "admin-setting-role-status", Sort: 0, Status: 1},

		{ID: 16, PID: 3, Name: "菜单", Icon: "", Type: "menu", MenuPath: "/admin/setting/menu", APIPath: "", Method: "", UniqueAuthKey: "admin-setting-menu", Sort: 0, Status: 1},
		{ID: 17, PID: 16, Name: "菜单列表", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/menu/tree", Method: "GET", UniqueAuthKey: "admin-setting-menu-tree", Sort: 0, Status: 1},
		{ID: 18, PID: 16, Name: "菜单添加编辑", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/menu", Method: "POST", UniqueAuthKey: "admin-setting-menu-save", Sort: 0, Status: 1},
		{ID: 19, PID: 16, Name: "菜单删除", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/menu/%v", Method: "DELETE", UniqueAuthKey: "admin-setting-menu-delete", Sort: 0, Status: 1},
		{ID: 20, PID: 16, Name: "菜单禁用或启用", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/menu/status/%v", Method: "GET", UniqueAuthKey: "admin-setting-menu-status", Sort: 0, Status: 1},
	}
	for _, m := range menus {
		global.Db.Create(m)
	}

	role := model.Role{
		Name:   "子管理员",
		Tag:    "goods_manager",
		Menus:  menus,
		Status: 1,
	}
	global.Db.Create(&role)

	goodseditor := model.AdminUser{
		Username: "subadmin",
		Password: util.HashPassword("123456"),
		// Role: model.Role{
		// 	Model: gorm.Model{ID:2},
		// },
		Roles: []model.Role{
			role,
		},
		Phone:  "13721047437",
		Status: 1,
	}
	global.Db.Create(&goodseditor)

	return "ok"
}

// GetInit2 初始化数据
func (p *Public) GetInit2() string {
	menus := []*model.Menu{
		{ID: 1, PID: 0, Name: "主页", Icon: "md-home", Type: "menu", MenuPath: "/admin/home/", APIPath: "", Method: "", UniqueAuthKey: "admin-home", Sort: 1, Status: 1},
		{ID: 2, PID: 1, Name: "首页统计接口", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/statistic", Method: "GET", UniqueAuthKey: "admin-home-statistic", Sort: 0, Status: 1},

		{ID: 3, PID: 0, Name: "权限管理", Icon: "md-settings", Type: "menu", MenuPath: "/admin/setting/", APIPath: "", Method: "", UniqueAuthKey: "admin-setting", Sort: 2, Status: 1},
		{ID: 4, PID: 3, Name: "管理员", Icon: "", Type: "menu", MenuPath: "/admin/setting/adminuser", APIPath: "", Method: "", UniqueAuthKey: "admin-setting-adminuser", Sort: 0, Status: 1},
		{ID: 5, PID: 4, Name: "管理员列表", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/adminuser/list/%v/%v", Method: "GET", UniqueAuthKey: "admin-setting-adminuser-list", Sort: 0, Status: 1},
		{ID: 6, PID: 4, Name: "管理员详情", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/adminuser/%v", Method: "GET", UniqueAuthKey: "admin-setting-adminuser-info", Sort: 0, Status: 1},
		{ID: 7, PID: 4, Name: "管理员添加编辑", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/adminuser", Method: "POST", UniqueAuthKey: "admin-setting-adminuser-save", Sort: 0, Status: 1},
		{ID: 8, PID: 4, Name: "管理员删除", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/adminuser/%v", Method: "DELETE", UniqueAuthKey: "admin-setting-adminuser-delete", Sort: 0, Status: 1},
		{ID: 9, PID: 4, Name: "管理员禁用启用", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/adminuser/status/%v", Method: "GET", UniqueAuthKey: "admin-setting-adminuser-status", Sort: 0, Status: 1},
	}
	// for _, m := range menus {
	// 	global.Db.Create(m)
	// }

	role := model.Role{
		Name:   "子管理员2",
		Tag:    "info_manager",
		Menus:  menus,
		Status: 1,
	}
	global.Db.Create(&role)

	goodseditor := model.AdminUser{
		Username: "subadmin2",
		Password: util.HashPassword("123456"),
		// Role: model.Role{
		// 	Model: gorm.Model{ID:2},
		// },
		Roles: []model.Role{
			role,
		},
		Phone:  "16963256896",
		Status: 1,
	}
	global.Db.Create(&goodseditor)
	return "ok"
}

// AfterActivation 后置方法
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
	// adminUseID := data["admin_user_id"].(string)
	var adminUseID string
	if value, ok := data[global.AdminUserJWTKey]; ok {
		adminUseID = value.(string)
	} else {
		app.ResponseProblemHTTPCode(p.Ctx, iris.StatusBadRequest, errors.New("Token中没有"+global.AdminUserJWTKey))
	}

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

	refreshToken, err := global.Redis.Get(config.App.Appname + ":refresh_token_admin_" + adminUseID).Result()
	if err == redis.Nil {
		p.Ctx.JSON(app.APIData(false, app.CodeRefreshTokenExpired, "", nil))
	} else if param.RefreshToken != refreshToken {
		p.Ctx.JSON(app.APIData(false, app.CodeRefreshTokenInvalidated, "", nil))
	} else {
		token, refreshToken := app.GenTokenAndRefreshToken(global.AdminUserJWTKey, util.ParseInt(adminUseID), global.AdminTokenMinutes, global.AdminRefreshTokenMinutes)
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

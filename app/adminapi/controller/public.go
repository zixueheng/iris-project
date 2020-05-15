package controller

import (
	"iris-project/app"
	"iris-project/app/adminapi/model"
	"iris-project/app/adminapi/validate"

	"github.com/kataras/iris/v12"
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

// PostLogin 登录
func (p *Public) PostLogin() {
	loginInfo := new(validate.LoginRequest)

	// if err := p.Ctx.ReadJSON(aul); err != nil {
	// 	// p.Ctx.StatusCode(iris.StatusOK)
	// 	_, _ = p.Ctx.JSON(app.APIData(false, nil, err.Error()))
	// 	return
	// }

	// err := global.Validate.Struct(*aul)
	// if err != nil {
	// 	errs := err.(validator.ValidationErrors)
	// 	for _, e := range errs.Translate(global.ValidateTrans) {
	// 		if len(e) > 0 {
	// 			// p.Ctx.StatusCode(iris.StatusOK)
	// 			_, _ = p.Ctx.JSON(app.APIData(false, nil, e))
	// 			return
	// 		}
	// 	}
	// }

	errmsg := app.CheckRequest(p.Ctx, loginInfo)
	if len(errmsg) != 0 {
		p.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}

	adminUser := &model.AdminUser{}
	response, ok, code := adminUser.CheckLogin(loginInfo)

	p.Ctx.JSON(app.APIData(ok, code, "", response))
}

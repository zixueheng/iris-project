package controller

import (
	"iris-project/app/adminapi/model"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// AdminUser 管理员控制器
type AdminUser struct {
	Ctx       iris.Context     // IRIS框架会自动注入 Context
	AdminUser *model.AdminUser // 通过执行依赖函数 GetAuthAdminUser 动态注入
}

// BeforeActivation 前置方法
func (au *AdminUser) BeforeActivation(b mvc.BeforeActivation) {
	b.Dependencies().Add(GetAuthAdminUser) // 注入依赖函数 GetAuthAdminUser
}

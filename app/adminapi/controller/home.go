package controller

import (
	"iris-project/app/adminapi/model"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// Home 控制器
type Home struct {
	Ctx           iris.Context     // IRIS框架会自动注入 Context
	AuthAdminUser *model.AdminUser // 通过执行依赖函数 GetAuthAdminUser 动态注入
}

// BeforeActivation 前置方法
func (h *Home) BeforeActivation(b mvc.BeforeActivation) {
	b.Dependencies().Add(GetAuthAdminUser) // 注入依赖函数 GetAuthAdminUser
}

// GetHomeStatistic 首页统计
func (h *Home) GetHomeStatistic() {

}

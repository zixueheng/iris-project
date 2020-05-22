package controller

import (
	"iris-project/app"
	"iris-project/app/adminapi/model"
	"iris-project/global"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// Role 控制器
type Role struct {
	Ctx           iris.Context     // IRIS框架会自动注入 Context
	AuthAdminUser *model.AdminUser // 通过执行依赖函数 GetAuthAdminUser 动态注入
}

// BeforeActivation 前置方法
func (r *Role) BeforeActivation(b mvc.BeforeActivation) {
	b.Dependencies().Add(GetAuthAdminUser) // 注入依赖函数 GetAuthAdminUser
}

// GetRoleListBy 角色列表
func (r *Role) GetRoleListBy(page, size uint) {
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}
	var where = make(map[string]interface{})
	if name := r.Ctx.URLParamDefault("name", ""); len(name) > 0 {
		where["name"] = name
	}

	var (
		roles []model.Role
		total uint
	)
	global.Db.Where(where).Offset((page - 1) * size).Limit(size).Find(&roles).Offset(-1).Limit(-1).Count(&total)

	r.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", app.List{List: roles, Total: total}))
}

// PostRole 创建或更新角色
func (r *Role) PostRole() {

}

// DeleteRoleBy 删除角色
func (r *Role) DeleteRoleBy(id uint) {

}

// GetRoleStatusBy 禁用或启用角色
func (r *Role) GetRoleStatusBy(id uint) {
	role := new(model.Role)
	if !role.GetRoleByID(id) {
		r.Ctx.JSON(app.APIData(false, app.CodeNotFound, "", nil))
		return
	}

	if role.Status == 1 {
		if global.Db.Model(role).Update("status", -1).RowsAffected > 0 {
			r.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
			return
		}
	} else if role.Status == -1 {
		if global.Db.Model(role).Update("status", 1).RowsAffected > 0 {
			r.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
			return
		}
	}
	r.Ctx.JSON(app.APIData(false, app.CodeFailed, "", nil))
}

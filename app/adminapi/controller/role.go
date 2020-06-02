package controller

import (
	"iris-project/app"
	"iris-project/app/adminapi/model"
	"iris-project/app/adminapi/validate"
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

// GetRoleBy 角色详情
func (r *Role) GetRoleBy(id uint) {
	role := new(model.Role)
	if !role.GetRoleMenusTreeByID(id) {
		r.Ctx.JSON(app.APIData(false, app.CodeNotFound, "", nil))
		return
	}
	r.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", role))
}

// PostRole 创建或更新角色
func (r *Role) PostRole() {
	postInfo := new(validate.RoleRequest)

	errmsg := app.CheckRequest(r.Ctx, postInfo)
	if len(errmsg) != 0 {
		r.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}
	role := &model.Role{
		ID:     postInfo.ID,
		Name:   postInfo.Name,
		Tag:    postInfo.Tag,
		Status: postInfo.Status,
	}
	menus := make([]*model.Menu, 0)
	for _, mid := range postInfo.MenuIds {
		menu := new(model.Menu)
		if menu.GetMenuByID(mid) {
			menus = append(menus, menu)
		} else {
			r.Ctx.JSON(app.APIData(false, app.CodeRequestParamError, "", nil))
			return
		}
	}
	role.Menus = menus

	if err := role.CreateUpdateRole(); err != nil {
		r.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	r.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// DeleteRoleBy 删除角色
func (r *Role) DeleteRoleBy(id uint) {
	role := new(model.Role)
	if !role.GetRoleByID(id) {
		r.Ctx.JSON(app.APIData(false, app.CodeNotFound, "", nil))
		return
	}

	if role.Tag == global.SuperAdminUserTag { // 超级管理员禁止删除
		r.Ctx.JSON(app.APIData(false, app.CodeForbidden, "", nil))
		return
	}
	if global.Db.Unscoped().Delete(role).RowsAffected > 0 {
		r.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
		return
	}
	r.Ctx.JSON(app.APIData(false, app.CodeFailed, "", nil))
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

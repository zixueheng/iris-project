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
		where["name like"] = "%" + name + "%"
	}
	if status := r.Ctx.URLParamIntDefault("status", 0); status != 0 {
		where["status"] = status
	}
	conditionString, conditionValues, _ := app.BuildCondition(where)

	var (
		roles []model.Role
		total uint
	)
	global.Db.Where(conditionString, conditionValues...).Order("id desc").Offset((page - 1) * size).Limit(size).Find(&roles).Offset(-1).Limit(-1).Count(&total)
	for i, role := range roles {
		obj := &model.Role{ID: role.ID}
		obj.GetRoleMenus()
		if obj.Tag == global.SuperAdminUserTag {
			roles[i].MenuNames = "所有权限"
		} else {
			for j, m := range obj.Menus {
				if j+1 == len(obj.Menus) {
					roles[i].MenuNames = roles[i].MenuNames + m.Name
				} else {
					roles[i].MenuNames = roles[i].MenuNames + m.Name + ","
				}
			}
		}
	}

	r.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", app.List{List: roles, Total: total}))
}

// GetRolesSelect 角色下拉选项
func (au *AdminUser) GetRolesSelect() {
	var roles []model.Role
	global.Db.Where("status=?", 1).Find(&roles)
	if roles == nil {
		roles = make([]model.Role, 0)
	}
	au.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", roles))
}

// GetRoleInfoBy 角色详情
func (r *Role) GetRoleInfoBy(id uint) {
	role := new(model.Role)
	role.ID = id
	// if !role.GetRoleMenusTree() {
	// 	r.Ctx.JSON(app.APIData(false, app.CodeNotFound, "", nil))
	// 	return
	// }
	if !role.GetRole() {
		r.Ctx.JSON(app.APIData(false, app.CodeNotFound, "", nil))
		return
	}
	var allMenus []*model.Menu
	global.Db.Where("status=?", 1).Find(&allMenus)
	allTree := model.GetTreeMenus(allMenus) // 加载所有菜单树
	var menuIDS = make([]uint, 0)
	if role.Tag == global.SuperAdminUserTag {
		for _, v := range allMenus {
			menuIDS = append(menuIDS, v.ID)
		}
	} else {
		global.Db.Find(&model.RoleMenu{}).Where("role_id=?", id).Pluck("menu_id", &menuIDS) // 该角色所有菜单ID
	}
	role.MenuIDS = menuIDS

	var (
		inArray func(uint, []uint) bool
		fn      func([]*model.MenuTree)
	)
	inArray = func(value uint, array []uint) bool {
		for _, a := range array {
			if value == a {
				return true
			}
		}
		return false
	}
	fn = func(tree []*model.MenuTree) {
		for i, t := range tree {
			if inArray(t.ID, menuIDS) || role.Tag == global.SuperAdminUserTag { // 判断是否拥有该菜单
				tree[i].Checked = true
			}
			if len(t.Children) > 0 {
				fn(t.Children)
			}
		}
	}
	fn(allTree)
	role.MenusTree = allTree

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
	if postInfo.ID != 0 {
		obj := &model.Role{ID: postInfo.ID}
		obj.GetRole()
		if obj.Tag == global.SuperAdminUserTag {
			r.Ctx.JSON(app.APIData(false, app.CodeCustom, "超级管理员角色禁止编辑", nil))
			return
		}
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
		menu.ID = mid
		if menu.GetMenu() {
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
	role.ID = id
	if !role.GetRole() {
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
	role.ID = id
	if !role.GetRole() {
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

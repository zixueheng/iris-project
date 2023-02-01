/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:28:08
 * @LastEditTime: 2021-10-12 17:40:13
 */
package controller

import (
	"iris-project/app"
	"iris-project/app/adminapi/model"
	"iris-project/app/adminapi/validate"
	"iris-project/app/dao"
	"iris-project/global"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// Role 控制器
type Role struct {
	Ctx           iris.Context
	AuthAdminUser *model.AdminUser
}

// BeforeActivation 前置方法
func (r *Role) BeforeActivation(b mvc.BeforeActivation) {
	b.Dependencies().Register(GetAuthAdminUser)
}

// GetRoleList 角色列表
func (r *Role) GetRoleList() {
	page, size := app.Pager(r.Ctx)

	var where = make(map[string]interface{})
	if name := r.Ctx.URLParamDefault("name", ""); len(name) > 0 {
		where["name like"] = "%" + name + "%"
	}
	if status := r.Ctx.URLParamIntDefault("status", 0); status != 0 {
		where["status"] = status
	}

	var (
		roles []model.Role
		total int64
	)
	searchList := &dao.SearchListData{
		Where:    where,
		Fields:   []string{},
		OrderBys: []string{"id desc"},
		Preloads: []string{},
		Page:     page,
		Size:     size,
	}
	if err := searchList.GetList(&roles, &total); err != nil {
		r.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

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
func (r *Role) GetRoleInfoBy(id uint32) {
	role := new(model.Role)
	role.ID = id
	// if !role.GetRoleMenusTree() {
	// 	r.Ctx.JSON(app.APIData(false, app.CodeNotFound, "", nil))
	// 	return
	// }
	if !role.GetOne(false) { // 兼容前端，接受0
		// r.Ctx.JSON(app.APIData(false, app.CodeNotFound, "", nil))
		// return
		role.Status = 1
	}
	var allMenus []*model.Menu
	global.Db.Where("status=?", 1).Order("sort asc").Find(&allMenus)
	allTree := model.GetTreeMenus(allMenus) // 加载所有菜单树
	var menuIDS = make([]uint32, 0)
	if role.Tag == global.SuperAdminUserTag {
		for _, v := range allMenus {
			menuIDS = append(menuIDS, v.ID)
		}
	} else {
		global.Db.Model(&model.RoleMenu{}).Where("role_id=?", id).Pluck("menu_id", &menuIDS) // 该角色所有菜单ID
	}
	role.MenuIDS = menuIDS

	var (
		inArray func(uint32, []uint32) bool
		fn      func([]*model.MenuTree)
	)
	inArray = func(value uint32, array []uint32) bool {
		for _, a := range array {
			if value == a {
				return true
			}
		}
		return false
	}
	fn = func(tree []*model.MenuTree) {
		for i, t := range tree {
			// if inArray(t.ID, menuIDS) || role.Tag == global.SuperAdminUserTag { // 判断是否拥有该菜单
			// 	// tree[i].Checked = false
			// 	// tree[i].Selected = false
			// }
			if len(t.Children) > 0 {
				fn(t.Children)
			} else { // 叶子节点
				if inArray(t.ID, menuIDS) || role.Tag == global.SuperAdminUserTag { // 判断是否拥有该菜单
					tree[i].Checked = true
					// tree[i].Selected = false
				} else {
					tree[i].Checked = false
				}

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
		obj.GetOne(false)
		if obj.Tag == global.SuperAdminUserTag {
			r.Ctx.JSON(app.APIData(false, app.CodeCustom, "超级管理员角色禁止编辑", nil))
			return
		}
	}
	role := &model.Role{
		ID:       postInfo.ID,
		Name:     postInfo.Name,
		JumpPage: postInfo.JumpPage,
		Tag:      postInfo.Tag,
		Status:   postInfo.Status,
	}

	menus := make([]*model.Menu, 0)
	for _, mid := range postInfo.MenuIds {
		menu := new(model.Menu)
		menu.ID = mid
		if menu.GetOne(false) {
			menus = append(menus, menu)
		} else {
			r.Ctx.JSON(app.APIData(false, app.CodeRequestParamError, "", nil))
			return
		}
	}
	role.Menus = menus

	if err := role.CreateUpdate(); err != nil {
		r.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	r.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// DeleteRoleBy 删除角色
func (r *Role) DeleteRoleBy(id uint32) {
	role := new(model.Role)
	role.ID = id
	if !role.GetOne(false) {
		r.Ctx.JSON(app.APIData(false, app.CodeNotFound, "", nil))
		return
	}

	if role.Tag == global.SuperAdminUserTag { // 超级管理员禁止删除
		r.Ctx.JSON(app.APIData(false, app.CodeForbidden, "", nil))
		return
	}

	if err := role.Delete(); err != nil {
		r.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}
	r.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// GetRoleStatusBy 禁用或启用角色
func (r *Role) GetRoleStatusBy(id uint32) {
	role := new(model.Role)
	role.ID = id
	if !role.GetOne(false) {
		r.Ctx.JSON(app.APIData(false, app.CodeNotFound, "", nil))
		return
	}

	if role.Status == 1 {
		if err := role.Updates(map[string]interface{}{"status": -1}); err != nil {
			r.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
			return
		}
	} else if role.Status == -1 {
		if err := role.Updates(map[string]interface{}{"status": 1}); err != nil {
			r.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
			return
		}
	}
	r.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

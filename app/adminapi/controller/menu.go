/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:28:08
 * @LastEditTime: 2022-01-04 17:16:38
 */
package controller

import (
	"iris-project/app"
	"iris-project/app/adminapi/model"
	"iris-project/app/adminapi/validate"
	"iris-project/app/dao"
	"iris-project/global"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// Menu 控制器
type Menu struct {
	Ctx           iris.Context
	AuthAdminUser *model.AdminUser
}

// BeforeActivation 前置方法
func (m *Menu) BeforeActivation(b mvc.BeforeActivation) {
	b.Dependencies().Register(GetAuthAdminUser)
}

// GetMenuTree 所有的菜单转成树状结构，一次返回
func (m *Menu) GetMenuTree() {
	tree := model.GetTreeMenus(nil)
	m.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", tree))
}

// GetMenuSelect 菜单下拉选择项
func (m *Menu) GetMenuSelect() {
	var list []*model.Menu
	global.Db.Where("status=?", 1).Order("sort asc").Find(&list)

	var box = make([]*model.Menu, 0)
	var root = &model.Menu{ID: 0, Name: "顶级按钮", Level: 0, HTML: ""}
	box = append(box, root)
	var fn func([]*model.Menu, *model.Menu)
	fn = func(ls []*model.Menu, node *model.Menu) {
		for _, v := range ls {
			if v.PID == node.ID {
				v.Level = node.Level + 1
				v.HTML = strings.Repeat("|----", v.Level)
				box = append(box, v)
				fn(ls, v)
			}
		}
	}
	fn(list, root)

	m.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", box))
}

// GetMenuSelect 菜单下拉选择项
func (m *Menu) GetMenuSelect2() {
	type (
		Menu struct {
			ID   uint32 `json:"value"`
			PID  uint32 `json:"-"`
			Name string `json:"label"`
		}
		MenuTree struct {
			*Menu
			Children []*MenuTree `json:"children,omitempty"`
		}
	)
	var (
		allMenuTree []*MenuTree
		allMenus    []*Menu
	)
	dao.FindAll(nil, &allMenus, map[string]interface{}{"status": 1}, []string{}, []string{"sort asc"})
	rootMenuTree := &MenuTree{
		Menu: &Menu{ID: 0, Name: "根节点"},
	}
	// allMenuTree = append(allMenuTree, rootMenuTree)
	for _, m := range allMenus {
		mt := &MenuTree{
			Menu: m,
		}
		allMenuTree = append(allMenuTree, mt)
	}
	var (
		makeMenuTree func(list []*MenuTree, menuTree *MenuTree)
		hasMenuChild func(list []*MenuTree, node *MenuTree) (children []*MenuTree, has bool)
	)
	makeMenuTree = func(list []*MenuTree, menuTree *MenuTree) {
		if children, has := hasMenuChild(list, menuTree); has {
			menuTree.Children = append(menuTree.Children, children...)
			for _, v := range children {
				makeMenuTree(list, v)
			}
		} else {
			menuTree.Children = make([]*MenuTree, 0) // 没有子节点将 nil 转成空切片，输出json 是[] 而不是 null
		}
	}
	hasMenuChild = func(list []*MenuTree, node *MenuTree) (children []*MenuTree, has bool) {
		for _, m := range list {
			if m.PID == node.ID {
				children = append(children, m)
			}
		}
		if len(children) > 0 {
			has = true
		}
		return
	}
	makeMenuTree(allMenuTree, rootMenuTree)
	m.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", rootMenuTree.Children))
}

// PostMenu 创建或更新菜单
// menu.save
func (m *Menu) PostMenu() {
	postInfo := new(validate.MenuRequest)

	errmsg := app.CheckRequest(m.Ctx, postInfo)
	if len(errmsg) != 0 {
		m.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}

	menu := &model.Menu{
		ID:            postInfo.ID,
		PID:           postInfo.PID,
		Name:          postInfo.Name,
		Icon:          postInfo.Icon,
		Type:          postInfo.Type,
		MenuPath:      postInfo.MenuPath,
		APIPath:       postInfo.APIPath,
		Method:        postInfo.Method,
		UniqueAuthKey: postInfo.UniqueAuthKey,
		Action:        postInfo.Action,
		Subject:       postInfo.Subject,
		Header:        postInfo.Header,
		IsHeader:      postInfo.IsHeader,
		Sort:          postInfo.Sort,
		Status:        postInfo.Status,
	}

	if err := menu.CreateUpdate(); err != nil {
		m.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	m.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// DeleteMenuBy 删除菜单(包含子菜单)
func (m *Menu) DeleteMenuBy(id uint32) {
	menu := new(model.Menu)
	menu.ID = id
	if !menu.GetOne(false) {
		m.Ctx.JSON(app.APIData(false, app.CodeNotFound, "", nil))
		return
	}

	if err := menu.Delete(); err != nil {
		m.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}
	m.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))

}

// GetMenuStatusBy 禁用或启用菜单
func (m *Menu) GetMenuStatusBy(id uint32) {
	menu := new(model.Menu)
	menu.ID = id
	if !menu.GetOne(false) {
		m.Ctx.JSON(app.APIData(false, app.CodeNotFound, "", nil))
		return
	}

	if menu.Status == 1 {
		if err := menu.Updates(map[string]interface{}{"status": -1}); err != nil {
			m.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
			return
		}
	} else if menu.Status == -1 {
		if err := menu.Updates(map[string]interface{}{"status": 1}); err != nil {
			m.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
			return
		}
	}
	m.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

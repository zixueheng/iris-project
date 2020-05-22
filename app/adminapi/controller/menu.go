package controller

import (
	"iris-project/app"
	"iris-project/app/adminapi/model"
	"iris-project/app/adminapi/validate"
	"iris-project/global"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// Menu 控制器
type Menu struct {
	Ctx           iris.Context     // IRIS框架会自动注入 Context
	AuthAdminUser *model.AdminUser // 通过执行依赖函数 GetAuthAdminUser 动态注入
}

// BeforeActivation 前置方法
func (m *Menu) BeforeActivation(b mvc.BeforeActivation) {
	b.Dependencies().Add(GetAuthAdminUser) // 注入依赖函数 GetAuthAdminUser
}

// GetMenuTree 所有的菜单转成树状结构，一次返回
func (m *Menu) GetMenuTree() {
	menu := &model.Menu{}
	tree := menu.GetTreeMenus()
	m.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", tree))
}

// PostMenu 创建或更新菜单
func (m *Menu) PostMenu() {
	postInfo := new(validate.MenuRequest)

	errmsg := app.CheckRequest(m.Ctx, postInfo)
	if len(errmsg) != 0 {
		m.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}

	menu := &model.Menu{
		ID:       postInfo.ID,
		PID:      postInfo.PID,
		Name:     postInfo.Name,
		Icon:     postInfo.Icon,
		Type:     postInfo.Type,
		MenuPath: postInfo.MenuPath,
		APIPath:  postInfo.APIPath,
		Method:   postInfo.Method,
		Sort:     postInfo.Sort,
		Status:   postInfo.Status,
	}

	if err := menu.CreateUpdateMenu(); err != nil {
		m.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	m.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// DeleteMenuBy 删除菜单(包含子菜单)
func (m *Menu) DeleteMenuBy(id uint) {
	menu := new(model.Menu)
	if !menu.GetMenuByID(id) {
		m.Ctx.JSON(app.APIData(false, app.CodeNotFound, "", nil))
		return
	}
	var delIDs []uint
	delIDs = append(delIDs, menu.ID)
	getChild(&delIDs, menu.ID) // 查找所有子菜单
	// fmt.Println(delIDs)

	if global.Db.Unscoped().Where("id IN (?)", delIDs).Delete(&model.Menu{}).RowsAffected > 0 {
		m.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
		return
	}
	m.Ctx.JSON(app.APIData(false, app.CodeFailed, "", nil))
}

func hasChild(id uint) (ids []uint, has bool) {
	global.Db.Model(&model.Menu{}).Where("p_id=?", id).Pluck("id", &ids)
	if len(ids) > 0 {
		has = true
	}
	return
}

func getChild(allIDs *[]uint, id uint) {
	if cids, has := hasChild(id); has {
		*allIDs = append(*allIDs, cids...)
		for _, cid := range cids {
			getChild(allIDs, cid)
		}
	}
}

// GetMenuStatusBy 禁用或启用菜单
func (m *Menu) GetMenuStatusBy(id uint) {
	menu := new(model.Menu)
	if !menu.GetMenuByID(id) {
		m.Ctx.JSON(app.APIData(false, app.CodeNotFound, "", nil))
		return
	}

	if menu.Status == 1 {
		if global.Db.Model(menu).Update("status", -1).RowsAffected > 0 {
			m.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
			return
		}
	} else if menu.Status == -1 {
		if global.Db.Model(menu).Update("status", 1).RowsAffected > 0 {
			m.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
			return
		}
	}
	m.Ctx.JSON(app.APIData(false, app.CodeFailed, "", nil))
}

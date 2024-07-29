/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:28:08
 * @LastEditTime: 2024-07-29 11:27:27
 */
package controller

import (
	"iris-project/app"
	"iris-project/app/adminapi/model"
	"iris-project/app/adminapi/validate"
	"iris-project/app/dao"
	"iris-project/global"
	"iris-project/lib/util"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// AdminUser 管理员控制器
type AdminUser struct {
	Ctx           iris.Context
	AuthAdminUser *model.AdminUser
}

// BeforeActivation 前置方法
func (au *AdminUser) BeforeActivation(b mvc.BeforeActivation) {
	b.Dependencies().Register(GetAuthAdminUser)
}

// @Tags        管理员
// @Summary		管理员列表
// @Description	管理员列表
// @Accept		json
// @Produce		json
// @Param		username query string false	"用户名"
// @Param		status   query int    false	"状态"
// @Param		page     query string false	"页码"
// @Param		size     query int    false	"页大小"
// @Success		200		{object}	app.Response		""
// @Failure		200		{object}	app.Response	    ""
// @Router		/adminapi/adminuser/list [get]
func (au *AdminUser) GetAdminuserList() {
	page, size := app.Pager(au.Ctx)

	var where = make(map[string]interface{})
	if username := au.Ctx.URLParamDefault("username", ""); len(username) > 0 {
		where["username like"] = "%" + username + "%"
	}
	if status := au.Ctx.URLParamIntDefault("status", 0); status != 0 {
		where["status"] = status
	}

	var (
		adminUsers []*model.AdminUser
		total      int64
	)

	searchList := &dao.SearchListData{
		Where:    where,
		Fields:   []string{},
		OrderBys: []string{"id desc"},
		Preloads: []string{"Roles"},
		Page:     page,
		Size:     size,
	}
	if err := searchList.GetList(&adminUsers, &total); err != nil {
		au.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	for i, adminUser := range adminUsers {
		adminUser.GetOne(false)
		for _, role := range adminUser.Roles {
			if role.Tag == global.SuperAdminUserTag {
				adminUsers[i].SuperAdmin = true
			}
		}
	}

	au.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", app.List{List: adminUsers, Total: total}))
}

// @Tags        管理员
// @Summary		管理员详情
// @Description	管理员详情
// @Accept		json
// @Produce		json
// @Param		id path int true	"ID"
// @Success		200		{object}	app.Response		""
// @Failure		200		{object}	app.Response	    ""
// @Router		/adminapi/adminuser/info/{id} [get]
func (au *AdminUser) GetAdminuserInfoBy(id uint32) {
	adminUser := new(model.AdminUser)
	adminUser.ID = id
	if !adminUser.GetOne(false) {
		au.Ctx.JSON(app.APIData(true, app.CodeUserNotFound, "", nil))
		return
	}

	au.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", *adminUser))
}

// @Tags        管理员
// @Summary		创建或更新管理员
// @Description	创建或更新管理员
// @Accept		json
// @Produce		json
// @Param		account	body		validate.CreateUpdateAdminUserRequest	true	"请求体"
// @Success		200		{object}	app.Response		""
// @Failure		200		{object}	app.Response	    ""
// @Router		/adminapi/adminuser [post]
func (au *AdminUser) PostAdminuser() {
	postInfo := new(validate.CreateUpdateAdminUserRequest)

	errmsg := app.CheckRequest(au.Ctx, postInfo)
	if len(errmsg) != 0 {
		au.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}

	adminUser := &model.AdminUser{
		ID:       postInfo.ID,
		Username: postInfo.Username,
		Realname: postInfo.Realname,
		Phone:    postInfo.Phone,
		Status:   postInfo.Status,
	}
	roles := make([]model.Role, 0)
	for _, rid := range postInfo.RoleIds {
		role := new(model.Role)
		role.ID = rid
		role.GetOne(false)
		if role.Tag == global.SuperAdminUserTag && len(postInfo.RoleIds) > 1 {
			au.Ctx.JSON(app.APIData(false, app.CodeCustom, "超级管理员不可同时拥有其他角色", nil))
			return
		}
		roles = append(roles, *role)
	}
	adminUser.Roles = roles

	if postInfo.Password != "" {
		adminUser.Password = util.HashPassword(postInfo.Password)
	}

	if err := adminUser.CreateUpdate(); err != nil {
		au.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	au.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// @Tags        管理员
// @Summary		删除管理员
// @Description	删除管理员
// @Accept		json
// @Produce		json
// @Param		id path int true	"ID"
// @Success		200		{object}	app.Response		""
// @Failure		200		{object}	app.Response	    ""
// @Router		/adminapi/adminuser/{id} [delete]
func (au *AdminUser) DeleteAdminuserBy(id uint32) {
	adminUser := new(model.AdminUser)
	adminUser.ID = id
	if !adminUser.GetOne(false) {
		au.Ctx.JSON(app.APIData(false, app.CodeUserNotFound, "", nil))
		return
	}

	if adminUser.SuperAdmin { // 超级管理员禁止删除
		au.Ctx.JSON(app.APIData(false, app.CodeForbidden, "", nil))
		return
	}
	// if global.Db.Unscoped().Delete(adminUser).RowsAffected > 0 {
	// 	au.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
	// 	return
	// }
	if err := adminUser.Delete(); err != nil {
		au.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}
	au.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// @Tags        管理员
// @Summary		禁用或启用管理员
// @Description	禁用或启用管理员
// @Accept		json
// @Produce		json
// @Param		id path int true	"ID"
// @Success		200		{object}	app.Response		""
// @Failure		200		{object}	app.Response	    ""
// @Router		/adminapi/adminuser/{id} [get]
func (au *AdminUser) GetAdminuserStatusBy(id uint32) {
	adminUser := new(model.AdminUser)
	adminUser.ID = id
	if !adminUser.GetOne(false) {
		au.Ctx.JSON(app.APIData(false, app.CodeUserNotFound, "", nil))
		return
	}

	// if adminUser.SuperAdmin {
	// 	au.Ctx.JSON(app.APIData(false, app.CodeForbidden, "", nil))
	// 	return
	// }

	if adminUser.Status == 1 {
		if err := adminUser.Updates(map[string]interface{}{"status": -1}); err != nil {
			au.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
			return
		}
	} else if adminUser.Status == -1 {
		if err := adminUser.Updates(map[string]interface{}{"status": 1}); err != nil {
			au.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
			return
		}
	}
	au.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

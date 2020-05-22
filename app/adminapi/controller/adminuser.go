package controller

import (
	"iris-project/app"
	"iris-project/app/adminapi/model"
	"iris-project/app/adminapi/validate"
	"iris-project/global"
	"iris-project/lib/util"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// AdminUser 管理员控制器
type AdminUser struct {
	Ctx           iris.Context     // IRIS框架会自动注入 Context
	AuthAdminUser *model.AdminUser // 通过执行依赖函数 GetAuthAdminUser 动态注入
}

// BeforeActivation 前置方法
func (au *AdminUser) BeforeActivation(b mvc.BeforeActivation) {
	b.Dependencies().Add(GetAuthAdminUser) // 注入依赖函数 GetAuthAdminUser
}

// GetAdminuserListBy 管理员列表
func (au *AdminUser) GetAdminuserListBy(page, size uint) {
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}
	var where = make(map[string]interface{})
	if username := au.Ctx.URLParamDefault("username", ""); len(username) > 0 {
		where["username"] = username
	}

	var (
		adminUsers []model.AdminUser
		total      uint
	)
	global.Db.Where(where).Preload("Role").Offset((page - 1) * size).Limit(size).Find(&adminUsers).Offset(-1).Limit(-1).Count(&total)

	au.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", app.List{List: adminUsers, Total: total}))
}

// GetAdminuserBy 获取管理员详情
func (au *AdminUser) GetAdminuserBy(id uint) {
	var (
		adminUser = new(model.AdminUser)
		role      = new(model.Role)
	)
	if !adminUser.GetAdminUserByID(id) {
		au.Ctx.JSON(app.APIData(true, app.CodeUserNotFound, "", nil))
		return
	}

	if role.GetRoleByID(adminUser.RoleID) {
		adminUser.Role = *role
	}
	au.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", *adminUser))
}

// PostAdminuser 创建或更新管理员
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
		RoleID:   postInfo.RoleID,
		Phone:    postInfo.Phone,
		Status:   postInfo.Status,
	}
	if postInfo.Password != "" {
		adminUser.Password = util.HashPassword(postInfo.Password)
	}

	if err := adminUser.CreateUpdateAdminUser(); err != nil {
		au.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	au.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// DeleteAdminuserBy 删除管理员
func (au *AdminUser) DeleteAdminuserBy(id uint) {
	adminUser := new(model.AdminUser)
	if !adminUser.GetAdminUserByID(id) {
		au.Ctx.JSON(app.APIData(false, app.CodeUserNotFound, "", nil))
		return
	}
	role := new(model.Role)
	if role.GetRoleByID(adminUser.RoleID) && role.Tag == global.SuperAdminUserTag { // 超级管理员禁止删除
		au.Ctx.JSON(app.APIData(false, app.CodeForbidden, "", nil))
		return
	}
	if global.Db.Unscoped().Delete(adminUser).RowsAffected > 0 {
		au.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
		return
	}
	au.Ctx.JSON(app.APIData(false, app.CodeFailed, "", nil))
}

// GetAdminuserStatusBy 禁用或启用管理员
func (au *AdminUser) GetAdminuserStatusBy(id uint) {
	adminUser := new(model.AdminUser)
	if !adminUser.GetAdminUserByID(id) {
		au.Ctx.JSON(app.APIData(false, app.CodeUserNotFound, "", nil))
		return
	}
	var role = new(model.Role)
	if role.GetRoleByID(adminUser.RoleID) && role.Tag == global.SuperAdminUserTag {
		au.Ctx.JSON(app.APIData(false, app.CodeForbidden, "", nil))
		return
	}

	if adminUser.Status == 1 {
		if global.Db.Model(adminUser).Update("status", -1).RowsAffected > 0 {
			au.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
			return
		}
	} else if adminUser.Status == -1 {
		if global.Db.Model(adminUser).Update("status", 1).RowsAffected > 0 {
			au.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
			return
		}
	}
	au.Ctx.JSON(app.APIData(false, app.CodeFailed, "", nil))
}

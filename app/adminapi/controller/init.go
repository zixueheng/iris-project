/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-08 14:17:15
 * @LastEditTime: 2023-02-20 15:30:59
 */
package controller

import (
	"fmt"
	adminmodel "iris-project/app/adminapi/model"
	appmodel "iris-project/app/model"
	"iris-project/global"
	"iris-project/lib/util"

	"github.com/kataras/iris/v12"
)

// Init 数据库初始化，不需要验证token
type Init struct {
	Ctx iris.Context // IRIS框架会自动注入 Context
}

// GetInit 初始化数据
func (p *Public) GetInit() string {
	err := global.Db.AutoMigrate(
		// &model.User{},
		&adminmodel.Menu{},
		&adminmodel.Role{},
		&adminmodel.AdminUser{},
		// &adminmodel.AdminUserRole{},
		// &adminmodel.RoleMenu{},
	)
	// return "ok"
	if err != nil {
		fmt.Println(err.Error())
	}
	// AutoMigrate 会创建表、缺失的外键、约束、列和索引。 如果大小、精度、是否为空可以更改，则 AutoMigrate 会改变列的类型。 出于保护您数据的目的，它 不会 删除未使用的列

	adminUser := adminmodel.AdminUser{
		Username: "admin",
		Realname: "总裁",
		Password: util.HashPassword("123456"),
		Roles: []adminmodel.Role{
			{
				Name:   "超级管理员",
				Tag:    "superadmin",
				Status: 1,
			},
		},
		Phone:  "15215657185",
		Status: 1,
	}
	global.Db.Create(&adminUser)

	menus := []*adminmodel.Menu{
		{ID: 1, PID: 0, Name: "主页", Icon: "md-home", Type: "menu", MenuPath: "/admin/home/", APIPath: "", Method: "", UniqueAuthKey: "admin-home", Header: "home", IsHeader: 1, Sort: 1, Status: 1},
		{ID: 2, PID: 1, Name: "首页统计", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/home/statistic", Method: "GET", UniqueAuthKey: "admin-home-statistic", Header: "home", IsHeader: 0, Sort: 0, Status: 1},

		{ID: 3, PID: 0, Name: "系统管理", Icon: "md-settings", Type: "menu", MenuPath: "/admin/setting", APIPath: "", Method: "", UniqueAuthKey: "admin-setting", Header: "settings", IsHeader: 1, Sort: 5, Status: 1},

		{ID: 4, PID: 3, Name: "管理员", Icon: "", Type: "menu", MenuPath: "/admin/setting/adminuser", APIPath: "", Method: "", UniqueAuthKey: "admin-setting-adminuser", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},
		{ID: 5, PID: 4, Name: "管理员列表", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/adminuser/list", Method: "GET", UniqueAuthKey: "admin-setting-adminuser-list", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},
		{ID: 6, PID: 4, Name: "管理员详情", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/adminuser/info/%v", Method: "GET", UniqueAuthKey: "admin-setting-adminuser-info", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},
		{ID: 7, PID: 4, Name: "管理员添加编辑", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/adminuser", Method: "POST", UniqueAuthKey: "admin-setting-adminuser-save", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},
		{ID: 8, PID: 4, Name: "管理员删除", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/adminuser/%v", Method: "DELETE", UniqueAuthKey: "admin-setting-adminuser-delete", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},
		{ID: 9, PID: 4, Name: "管理员禁用启用", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/adminuser/status/%v", Method: "GET", UniqueAuthKey: "admin-setting-adminuser-status", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},
		{ID: 10, PID: 4, Name: "选择角色下拉选项", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/roles/select", Method: "GET", UniqueAuthKey: "admin-setting-adminuser-roles-select", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},

		{ID: 11, PID: 3, Name: "角色", Icon: "", Type: "menu", MenuPath: "/admin/setting/role", APIPath: "", Method: "", UniqueAuthKey: "admin-setting-role", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},
		{ID: 12, PID: 11, Name: "角色列表", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/role/list", Method: "GET", UniqueAuthKey: "admin-setting-role-list", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},
		{ID: 13, PID: 11, Name: "角色添加编辑", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/role", Method: "POST", UniqueAuthKey: "admin-setting-role-save", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},
		{ID: 14, PID: 11, Name: "角色详情", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/role/info/%v", Method: "GET", UniqueAuthKey: "admin-setting-role-info", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},
		{ID: 15, PID: 11, Name: "角色删除", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/role/%v", Method: "DELETE", UniqueAuthKey: "admin-setting-role-delete", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},
		{ID: 16, PID: 11, Name: "角色禁用或启用", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/role/status/%v", Method: "GET", UniqueAuthKey: "admin-setting-role-status", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},

		{ID: 17, PID: 3, Name: "权限", Icon: "", Type: "menu", MenuPath: "/admin/setting/menu", APIPath: "", Method: "", UniqueAuthKey: "admin-setting-menu", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},
		{ID: 18, PID: 17, Name: "权限列表", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/menu/tree", Method: "GET", UniqueAuthKey: "admin-setting-menu-tree", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},
		{ID: 19, PID: 17, Name: "权限添加编辑", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/menu", Method: "POST", UniqueAuthKey: "admin-setting-menu-save", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},
		{ID: 20, PID: 17, Name: "权限删除", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/menu/%v", Method: "DELETE", UniqueAuthKey: "admin-setting-menu-delete", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},
		{ID: 21, PID: 17, Name: "权限禁用或启用", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/menu/status/%v", Method: "GET", UniqueAuthKey: "admin-setting-menu-status", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},
		{ID: 22, PID: 17, Name: "权限下拉选择项", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/menu/select", Method: "GET", UniqueAuthKey: "admin-setting-menu-select", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},

		{ID: 23, PID: 3, Name: "数据备份", Icon: "", Type: "menu", MenuPath: "/admin/setting/dbbackup", APIPath: "", Method: "", UniqueAuthKey: "admin-setting-dbbackup", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},
		{ID: 24, PID: 23, Name: "数据库备份", Icon: "", Type: "api", MenuPath: "", APIPath: "/adminapi/system/db/backup", Method: "GET", UniqueAuthKey: "admin-setting-dbbackup-request", Header: "settings", IsHeader: 0, Sort: 0, Status: 1},
	}
	for _, m := range menus {
		global.Db.Create(m)
	}

	role := adminmodel.Role{
		Name:   "子管理员",
		Tag:    "submanager",
		Menus:  menus,
		Status: 1,
	}
	global.Db.Create(&role)

	goodseditor := adminmodel.AdminUser{
		Username: "subadmin",
		Realname: "小王",
		Password: util.HashPassword("123456"),
		// Role: model.Role{
		// 	Model: gorm.Model{ID:2},
		// },
		Roles: []adminmodel.Role{
			role,
		},
		Phone:  "13721047437",
		Status: 1,
	}
	global.Db.Create(&goodseditor)

	return "ok"
}

// GetInit2 初始化数据
func (p *Public) GetInit2() string {
	err := global.Db.AutoMigrate(
		&appmodel.User{},
		&appmodel.UserRefreshToken{},
		&appmodel.UserWechat{},
		&appmodel.UserUsername{},
		&appmodel.UserPhone{},
		&appmodel.Verifycode{},

		&appmodel.Config{},
		&appmodel.File{},
		&appmodel.FileCategory{},
		&appmodel.FileFavor{},

		&appmodel.Address{},
	)
	if err != nil {
		fmt.Println(err.Error())
	}
	return "ok"
}

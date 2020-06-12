package model

import (
	"encoding/json"
	"errors"
	"iris-project/app"
	"iris-project/app/adminapi/validate"
	"iris-project/config"
	"iris-project/global"
	"iris-project/lib/util"
	"time"

	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/gorm"
)

// AdminUser 模型
//
// Token使用JWT(生成后由客户端持有，请求时发送到服务器进行验证)无需保存数据库；
// RefreshToken 随机生成32位字符串保存到Redis中，也无需保存到数据库
type AdminUser struct {
	// gorm.Model
	ID             uint           `gorm:"primary_key" json:"id"`
	CreatedAt      global.SQLTime `gorm:"type:datetime;" json:"created_at"`
	UpdatedAt      global.SQLTime `gorm:"type:datetime;" json:"updated_at"`
	Username       string         `gorm:"type:varchar(100);unique_index;not null" json:"username"`
	Password       string         `gorm:"type:varchar(100);not null" json:"-"`
	Roles          []Role         `gorm:"many2many:admin_user_role;association_autoupdate:false;" json:"roles"`
	SuperAdmin     bool           `gorm:"-" json:"super_admin"`      // 是否超级管理员
	Menus          []*Menu        `gorm:"-" json:"menus"`            // 所有菜单
	MenusTree      []*MenuTree    `gorm:"-" json:"menus_tree"`       // 菜单树
	UniqueAuthKeys []string       `gorm:"-" json:"unique_auth_keys"` // 所有鉴权key
	Phone          string         `gorm:"type:char(11);not null" json:"phone"`
	Status         int8           `gorm:"type:tinyint(1);default:1" json:"status"`
}

// CheckLogin 登录检查
func (au *AdminUser) CheckLogin(loginInfo *validate.LoginRequest) (interface{}, bool, app.Code) {
	if !au.GetAdminUserByUsername(loginInfo.Username) {
		return nil, false, app.CodeUserNotFound
	} else if au.Status != 1 {
		return nil, false, app.CodeUserForbidden
	} else {
		if bcrypt.Match(loginInfo.Password, au.Password) {
			token, refreshToken := app.GenTokenAndRefreshToken(global.AdminUserJWTKey, int(au.ID), global.AdminTokenMinutes, global.AdminRefreshTokenMinutes)
			au.GetAdminUserByID(au.ID)

			json, _ := json.Marshal(au)
			global.Redis.Set(config.App.Appname+":vo_admin_user_"+util.ParseString(int(au.ID)), string(json), time.Minute*time.Duration(global.AdminUserCacheMinutes)) // 账号信息保存到redis

			response := struct {
				Token        string `json:"token"`
				RefreshToken string `json:"refresh_token"`
				AdminUser    `json:"admin_user"`
			}{
				Token:        token,
				RefreshToken: refreshToken,
				AdminUser:    *au,
			}

			return response, true, app.CodeUserLoginSucceed

		}
		return nil, false, app.CodeUserPasswordError
	}

}

// GetAdminUserByID 根据ID获取管理员（包括角色、菜单、菜单树）
func (au *AdminUser) GetAdminUserByID(id uint) bool {
	if err := global.Db.Where("id=?", id).Preload("Roles", "status=?", 1).First(au).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}

	for _, role := range au.Roles {
		if role.Tag == global.SuperAdminUserTag {
			au.SuperAdmin = true
			menus := make([]*Menu, 0)
			global.Db.Where("status=?", 1).Find(&menus) // 所有菜单
			role.Menus = menus
		} else {
			role.GetRoleMenusByID(role.ID) // 加载角色菜单
		}

		for _, menu := range role.Menus {
			if !checkMenus(au.Menus, *menu) {
				au.Menus = append(au.Menus, menu)
			}
			if len(menu.UniqueAuthKey) > 0 && !checkUniqueAuthKeys(au.UniqueAuthKeys, menu.UniqueAuthKey) {
				au.UniqueAuthKeys = append(au.UniqueAuthKeys, menu.UniqueAuthKey)
			}
		}
	}
	au.MenusTree = GetTreeMenus(au.Menus)
	return true
	// return FindOne(au, map[string]interface{}{"username": username}, "Role")
}

func checkMenus(menus []*Menu, menu Menu) bool {
	for _, m := range menus {
		if m.ID == menu.ID {
			return true
		}
	}
	return false
}
func checkUniqueAuthKeys(keys []string, key string) bool {
	for _, s := range keys {
		if s == key {
			return true
		}
	}
	return false
}

// GetAdminUserByUsername 根据用户名获取管理员
func (au *AdminUser) GetAdminUserByUsername(username string) bool {
	if err := global.Db.Where("username=?", username).First(au).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
	// return FindOne(au, map[string]interface{}{"username": username}, "Role")
}

// CreateUpdateAdminUser 创建或更新管理员
func (au *AdminUser) CreateUpdateAdminUser() error {
	if au.ID == 0 {
		var count uint
		global.Db.Model(&AdminUser{}).Where("username=?", au.Username).Count(&count)
		if count > 0 {
			return errors.New("Username重复")
		}
		if err := global.Db.Create(au).Error; err != nil {
			return err
		}
	} else {
		var count uint
		global.Db.Model(&AdminUser{}).Where("username=? and id<>?", au.Username, au.ID).Count(&count)
		if count > 0 {
			return errors.New("Username重复")
		}
		// Todo
		// var role = new(Role)
		// if role.GetRoleByID(au.RoleID) && role.Tag == global.SuperAdminUserTag {
		// 	return errors.New("超级管理员禁止编辑")
		// }
		global.Db.Unscoped().Where("admin_user_id=?", au.ID).Delete(&AdminUserRole{}) // 删除原来的关联
		if err := global.Db.Model(au).Updates(*au).Error; err != nil {
			return err
		}
	}

	return nil
}

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
	Realname       string         `gorm:"type:varchar(100);not null" json:"realname"`
	Password       string         `gorm:"type:varchar(100);not null" json:"-"`
	Roles          []Role         `gorm:"many2many:admin_user_role;association_autoupdate:false;" json:"roles"`
	SuperAdmin     bool           `gorm:"-" json:"super_admin"` // 是否超级管理员
	Menus          []*Menu        `gorm:"-" json:"menus"`       // 所有菜单和接口
	MenusTree      []*MenuTree    `gorm:"-" json:"-"`           // 所有菜单和接口树
	OnlyMenusTree  []*MenuTree    `gorm:"-" json:"-"`           // 只有菜单树
	UniqueAuthKeys []string       `gorm:"-" json:"-"`           // 所有鉴权key
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

			au.GetAdminUser()

			json, _ := json.Marshal(au)
			global.Redis.Set(config.App.Appname+":vo_admin_user_"+util.ParseString(int(au.ID)), string(json), time.Minute*time.Duration(global.AdminUserCacheMinutes)) // 账号信息保存到redis

			type userInfo struct {
				ID         uint   `json:"id"`
				SuperAdmin bool   `json:"super_admin"`
				Username   string `json:"account"`
				Realname   string `json:"realname"`
				HeadPic    string `json:"head_pic"`
			}

			expiresTime := time.Now().Add(time.Minute * time.Duration(global.AdminTokenMinutes)) // token 的过期时间

			response := struct {
				Token        string      `json:"token"`
				ExpiresTime  int64       `json:"expires_time"`
				RefreshToken string      `json:"refresh_token"`
				Menus        []*MenuTree `json:"menus"`
				UniqueAuth   []string    `json:"unique_auth"`
				UserInfo     userInfo    `json:"user_info"`
				Logo         string      `json:"logo"`
				LogoSquare   string      `json:"logo_square"`
				Version      string      `json:"version"`
			}{
				Token:        token,
				ExpiresTime:  expiresTime.Unix(),
				RefreshToken: refreshToken,
				Menus:        au.OnlyMenusTree,
				UniqueAuth:   au.UniqueAuthKeys,
				UserInfo: userInfo{
					ID:         au.ID,
					Username:   au.Username,
					Realname:   au.Realname,
					SuperAdmin: au.SuperAdmin,
					HeadPic:    "",
				},
				// Logo:       config.App.Fronturl + "/Logo.png",
				// LogoSquare: config.App.Fronturl + "/LogoSquare.png",
				Version: "1.0.0",
			}

			return response, true, app.CodeUserLoginSucceed

		}
		return nil, false, app.CodeUserPasswordError
	}

}

// GetAdminUser 获取管理员（包括角色、菜单、菜单树）
func (au *AdminUser) GetAdminUser() bool {
	if au.ID == 0 {
		return false
	}
	if err := global.Db.Preload("Roles", "status=?", 1).First(au).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}

	for _, role := range au.Roles {
		if role.Tag == global.SuperAdminUserTag {
			au.SuperAdmin = true
			menus := make([]*Menu, 0)
			global.Db.Where("status=?", 1).Order("sort asc").Find(&menus) // 所有菜单
			role.Menus = menus
		} else {
			role.GetRoleMenus() // 加载角色菜单
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
	au.MenusTree = GetTreeMenus(au.Menus) // 菜单包含接口

	onlyMenus := make([]*Menu, 0)
	for _, m := range au.Menus {
		if m.Type == "menu" {
			onlyMenus = append(onlyMenus, m)
		}
	}
	au.OnlyMenusTree = GetTreeMenus(onlyMenus) // 菜单不包含接口
	return true
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

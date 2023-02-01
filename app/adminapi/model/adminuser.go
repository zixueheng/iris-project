/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-18 09:34:00
 * @LastEditTime: 2023-02-01 15:58:27
 */
package model

import (
	"encoding/json"
	"errors"
	"iris-project/app"
	"iris-project/app/adminapi/validate"
	"iris-project/app/config"
	"iris-project/app/dao"
	"iris-project/global"
	"iris-project/lib/cache"
	"iris-project/lib/util"
	"time"

	"github.com/jameskeane/bcrypt"
	"gorm.io/gorm"
)

// AdminUser 模型
//
// Token使用JWT(生成后由客户端持有，请求时发送到服务器进行验证)无需保存数据库；
// RefreshToken 随机生成32位字符串保存到数据库
type AdminUser struct {
	// gorm.Model
	ID                  uint32           `gorm:"primaryKey;" json:"id"`
	CreatedAt           global.LocalTime `gorm:"type:datetime;" json:"created_at"`
	UpdatedAt           global.LocalTime `gorm:"type:datetime;" json:"updated_at"`
	RefreshToken        string           `gorm:"type:varchar(100);" json:"-" validate:"-" comment:"刷新TOKEN"`
	RefreshTokenExpired uint32           `gorm:"" json:"-" validate:"-" comment:"刷新TOKEN过期时间戳"`
	Username            string           `gorm:"type:varchar(100);unique;not null;" json:"username"`
	Realname            string           `gorm:"type:varchar(100);not null" json:"realname"`
	Password            string           `gorm:"type:varchar(100);not null" json:"-"`
	Roles               []Role           `gorm:"many2many:admin_user_role;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;association_autoupdate:false;" json:"roles"`
	SuperAdmin          bool             `gorm:"-" json:"super_admin"`                // 是否超级管理员
	Menus               []*Menu          `gorm:"-" json:"menus"`                      // 所有菜单和接口
	MenusTree           []*MenuTree      `gorm:"-" json:"-"`                          // 所有菜单和接口树
	OnlyMenusTree       []*MenuTree      `gorm:"-" json:"-"`                          // 只有菜单树
	UniqueAuthKeys      []string         `gorm:"-" json:"unique_auth_keys,omitempty"` // 所有鉴权key
	Ability             []*Ability       `gorm:"-" json:"ability"`
	Phone               string           `gorm:"unique;type:char(11);not null;" json:"phone"`
	Status              int8             `gorm:"type:tinyint(1);default:1" json:"status"`
}

// GetID 返回ID，实现接口`dao.Model`的方法
func (au *AdminUser) GetID() uint32 {
	return au.ID
}

type Ability struct {
	Action  string `json:"action"`
	Subject string `json:"subject"`
}

// CheckLogin 登录检查
func (au *AdminUser) CheckLogin(loginInfo *validate.LoginRequest) (interface{}, bool, app.Code) {
	if !au.GetOneByUsername(loginInfo.Username) {
		return nil, false, app.CodeUserNotFound
	} else if au.Status != 1 {
		return nil, false, app.CodeUserForbidden
	} else {
		if bcrypt.Match(loginInfo.Password, au.Password) {
			token, refreshToken, tokenExpired, refreshTokenExpired := app.GenTokenAndRefreshToken(global.AdminUserJWTKey, int(au.ID), global.AdminTokenMinutes, global.AdminRefreshTokenMinutes)

			au.GetOne(true)
			dao.UpdateByID(au, map[string]interface{}{"refresh_token": refreshToken, "refresh_token_expired": refreshTokenExpired.Unix()}) // 保存刷新token和过期时间至数据库

			json, _ := json.Marshal(au)
			cache.Set(config.App.Appname+global.AdminUserCacheKeyPrefix+util.ParseString(int(au.ID)), string(json), time.Minute*time.Duration(global.AdminUserCacheMinutes)) // 账号信息保存到缓存
			type userInfo struct {
				ID          uint32   `json:"id"`
				SuperAdmin  bool     `json:"super_admin"`
				Username    string   `json:"username"`
				Realname    string   `json:"realname"`
				Phone       string   `json:"phone"`
				Roles       []string `json:"roles"`
				JumpPages   []string `json:"jump_pages"`
				Avatar      string   `json:"avatar"`
				CompanyID   uint32   `json:"company_id"`
				CompanyName string   `json:"company_name"`
			}
			// type ability struct {
			// 	Action  string `json:"action"`
			// 	Subject string `json:"subject"`
			// }

			// expiresTime := time.Now().Add(time.Minute * time.Duration(global.AdminTokenMinutes)) // token 的过期时间
			// var abilities = make([]*ability, 0)
			// for _, v := range au.UniqueAuthKeys {
			// 	abilities = append(abilities, &ability{Action: "read", Subject: v})
			// }
			var (
				roles     = make([]string, 0)
				jumpPages = make([]string, 0)
			)
			for _, v := range au.Roles {
				roles = append(roles, v.Name)
				jumpPages = append(jumpPages, v.JumpPage)
			}

			response := struct {
				Token               string      `json:"token"`
				TokenExpired        int64       `json:"token_expired"`
				RefreshToken        string      `json:"refresh_token"`
				RefreshTokenExpired int64       `json:"refresh_token_expired"`
				Menus               []*MenuTree `json:"menus"`
				UniqueAuth          []string    `json:"unique_auth"`
				Ability             []*Ability  `json:"ability"`
				UserInfo            userInfo    `json:"user_info"`
				Logo                string      `json:"logo"`
				LogoSquare          string      `json:"logo_square"`
				Version             string      `json:"version"`
			}{
				Token:               token,
				TokenExpired:        tokenExpired.Unix(),
				RefreshToken:        refreshToken,
				RefreshTokenExpired: refreshTokenExpired.Unix(),
				Menus:               au.OnlyMenusTree,
				UniqueAuth:          au.UniqueAuthKeys,
				Ability:             au.Ability,
				UserInfo: userInfo{
					ID:         au.ID,
					Username:   au.Username,
					Realname:   au.Realname,
					Phone:      au.Phone,
					SuperAdmin: au.SuperAdmin,
					Roles:      roles,
					JumpPages:  jumpPages,
					Avatar:     "",
				},
				Logo:       config.App.Fronturl + "/Logo.png",
				LogoSquare: config.App.Fronturl + "/LogoSquare.png",
				Version:    "1.0.0",
			}

			return response, true, app.CodeUserLoginSucceed

		}
		return nil, false, app.CodeUserPasswordError
	}

}

// GetOne 获取管理员（包括角色，loadMenu是否加载菜单和菜单树）
func (au *AdminUser) GetOne(loadMenu bool) bool {
	if au.ID == 0 {
		return false
	}
	if err := dao.GetDB().Preload("Roles", "status=?", 1).First(au).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	if loadMenu {
		for _, role := range au.Roles {
			if role.Tag == global.SuperAdminUserTag {
				au.SuperAdmin = true
				menus := make([]*Menu, 0)
				dao.GetDB().Where("status=?", 1).Order("sort asc").Find(&menus) // 所有菜单
				role.Menus = menus
			} else {
				role.GetRoleMenus() // 加载角色菜单
			}

			for _, menu := range role.Menus {
				if !checkMenus(au.Menus, *menu) {
					au.Menus = append(au.Menus, menu)
					au.Ability = append(au.Ability, &Ability{Action: menu.Action, Subject: menu.Subject})
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
	}

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

// GetOneByUsername 根据用户名获取管理员
func (au *AdminUser) GetOneByUsername(username string) bool {
	if err := dao.GetDB().Where("username=?", username).First(au).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
	// return FindOne(au, map[string]interface{}{"username": username}, "Role")
}

// CreateUpdate 创建或更新管理员
func (au *AdminUser) CreateUpdate() error {
	if au.ID == 0 {
		var count int64
		dao.GetDB().Model(&AdminUser{}).Where("username=?", au.Username).Count(&count)
		if count > 0 {
			return errors.New("Username重复")
		}
		dao.GetDB().Model(&AdminUser{}).Where("phone=?", au.Phone).Count(&count)
		if count > 0 {
			return errors.New("手机号重复")
		}
		if err := dao.GetDB().Create(au).Error; err != nil {
			return err
		}
	} else {
		var count int64
		dao.GetDB().Model(&AdminUser{}).Where("username=? and id<>?", au.Username, au.ID).Count(&count)
		if count > 0 {
			return errors.New("Username重复")
		}
		dao.GetDB().Model(&AdminUser{}).Where("phone=? and id<>?", au.Phone, au.ID).Count(&count)
		if count > 0 {
			return errors.New("手机号重复")
		}
		// Todo
		// var role = new(Role)
		// if role.GetRoleByID(au.RoleID) && role.Tag == global.SuperAdminUserTag {
		// 	return errors.New("超级管理员禁止编辑")
		// }
		dao.GetDB().Unscoped().Where("admin_user_id=?", au.ID).Delete(&AdminUserRole{}) // 删除原来的关联
		if err := dao.GetDB().Model(au).Save(*au).Error; err != nil {
			return err
		}
	}

	return nil
}

// Delete 删除管理员
func (au *AdminUser) Delete() error {
	if au.ID == 0 {
		return errors.New("需指定ID")
	}
	if err := dao.GetDB().Unscoped().Delete(au).Error; err != nil {
		return err
	}
	return nil
}

// Updates 更新管理员状态
func (au *AdminUser) Updates(data map[string]interface{}) error {
	if au.ID == 0 {
		return errors.New("需指定ID")
	}
	if err := dao.GetDB().Model(au).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

package model

import (
	"iris-project/app"
	"iris-project/app/adminapi/validate"
	"iris-project/global"

	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/gorm"
)

// AdminUser 模型
//
// Token使用JWT(生成后由客户端持有，请求时发送到服务器进行验证)无需保存数据库；
// RefreshToken 随机生成32位字符串保存到Redis中，也无需保存到数据库
type AdminUser struct {
	// gorm.Model
	ID        uint           `gorm:"primary_key"`
	CreatedAt global.SQLTime `gorm:"type:datetime;"`
	UpdatedAt global.SQLTime `gorm:"type:datetime;"`
	Username  string         `gorm:"type:varchar(100);unique_index;not null"`
	Password  string         `gorm:"type:varchar(100);not null"`
	RoleID    uint
	Role      Role
	Phone     string `gorm:"type:char(11);not null"`
	Status    int8   `gorm:"type:tinyint(1);default:1"`
}

// FindOne 查找一个
func FindOne(obj interface{}, where map[string]interface{}, preload string) bool {
	tx := global.Db.Where(where)

	// 这里 Preload() 不起作用
	if len(preload) != 0 {
		tx.Preload(preload)
	}

	if err := tx.First(obj).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
}

// CheckLogin 登录检查
func (au *AdminUser) CheckLogin(loginInfo *validate.LoginRequest) (interface{}, bool, app.Code) {
	if !au.GetAdminUserByUsername(loginInfo.Username) {
		return nil, false, app.CodeUserNotFound
	} else if au.Status != 1 {
		return nil, false, app.CodeUserForbidden
	} else {
		if bcrypt.Match(loginInfo.Password, au.Password) {
			token, refreshToken := app.GenTokenAndRefreshToken("admin_user_id", int(au.ID), global.AdminTokenMinutes, global.AdminRefreshTokenMinutes)
			response := struct {
				Token        string `json:"token"`
				RefreshToken string `json:"refresh_token"`
			}{
				Token:        token,
				RefreshToken: refreshToken,
			}

			return response, true, app.CodeUserLoginSucceed

		}
		return nil, false, app.CodeUserPasswordError
	}

}

// GetAdminUserByID 根据ID获取管理员
func (au *AdminUser) GetAdminUserByID(id int) bool {
	if err := global.Db.Where("id=?", id).First(au).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
	// return FindOne(au, map[string]interface{}{"username": username}, "Role")
}

// GetAdminUserByUsername 根据用户名获取管理员
func (au *AdminUser) GetAdminUserByUsername(username string) bool {
	if err := global.Db.Where("username=?", username).First(au).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
	// return FindOne(au, map[string]interface{}{"username": username}, "Role")
}

// CreateAdminUser 创建管理员
func (au *AdminUser) CreateAdminUser() error {
	if global.Db.NewRecord(au) {
		if err := global.Db.Create(&au).Error; err != nil {
			return err
		}
	}

	return nil
}

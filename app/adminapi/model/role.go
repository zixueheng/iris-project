package model

import (
	"iris-project/global"
	"time"

	"github.com/jinzhu/gorm"
)

// Role 角色model
type Role struct {
	// gorm.Model
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	Name      string `gorm:"type:varchar(50);not null"`
	Tag       string `gorm:"type:varchar(50);unique"`
	Menus     []Menu `gorm:"many2many:role_menu;"`
	Status    int8   `gorm:"type:tinyint(1);default:1"`
}

// GetRoleByID 根据ID获取角色（包含菜单）
func (r *Role) GetRoleByID(id uint) bool {
	if err := global.Db.Where("id=?", id).Preload("Menus").First(r).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
}

// CreateRole 创建管理员
func (r *Role) CreateRole() error {
	if global.Db.NewRecord(r) {
		if err := global.Db.Create(&r).Error; err != nil {
			return err
		}
	}

	return nil
}

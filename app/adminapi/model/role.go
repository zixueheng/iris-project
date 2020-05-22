package model

import (
	"iris-project/global"

	"github.com/jinzhu/gorm"
)

// Role 角色model
type Role struct {
	// gorm.Model
	ID        uint           `gorm:"primary_key" json:"id"`
	CreatedAt global.SQLTime `gorm:"type:datetime;" json:"created_at"`
	Name      string         `gorm:"type:varchar(50);not null" json:"name"`
	Tag       string         `gorm:"type:varchar(50);unique" json:"tag"`
	Menus     []Menu         `gorm:"many2many:role_menu;" json:"menus"`
	Status    int8           `gorm:"type:tinyint(1);default:1" json:"status"`
}

// GetRoleByID 根据ID获取角色（包含菜单）
func (r *Role) GetRoleByID(id uint) bool {
	if err := global.Db.Where("id=?", id).Preload("Menus").First(r).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
}

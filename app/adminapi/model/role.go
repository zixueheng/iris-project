package model

import (
	"iris-project/global"

	"github.com/jinzhu/gorm"
)

// Role 角色model
type Role struct {
	gorm.Model
	Name   string `gorm:"type:varchar(50);not null"`
	Tag    string `gorm:"type:varchar(50);unique"`
	Menus  []Menu `gorm:"many2many:role_menu;"`
	Status int8   `gorm:"type:tinyint(1);default:1"`
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

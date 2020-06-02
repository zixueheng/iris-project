package model

import (
	"errors"
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
	Menus     []*Menu        `gorm:"many2many:role_menu;association_autoupdate:false" json:"-"`
	MenusTree []*MenuTree    `gorm:"-" json:"menus_tree,omitempty"`
	Status    int8           `gorm:"type:tinyint(1);default:1" json:"status"`
}

// GetRoleByID 根据ID获取角色（不包含菜单）
func (r *Role) GetRoleByID(id uint) bool {
	if err := global.Db.Where("id=?", id).First(r).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
}

// GetRoleMenusByID 根据ID获取角色（包含菜单）
func (r *Role) GetRoleMenusByID(id uint) bool {
	if err := global.Db.Where("id=?", id).Preload("Menus").First(r).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
}

// GetRoleMenusTreeByID 根据ID获取角色（包含菜单和菜单数）
func (r *Role) GetRoleMenusTreeByID(id uint) bool {
	if err := global.Db.Where("id=?", id).Preload("Menus").First(r).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	r.MenusTree = GetTreeMenus(r.Menus)
	return true
}

// CreateUpdateRole 创建或更新角色
func (r *Role) CreateUpdateRole() error {
	if r.ID == 0 {
		var count uint
		global.Db.Model(&Role{}).Where("tag=?", r.Tag).Count(&count)
		if count > 0 {
			return errors.New("Tag重复")
		}
		if err := global.Db.Create(r).Error; err != nil {
			return err
		}
	} else {
		var count uint
		global.Db.Model(&Role{}).Where("tag=? and id<>?", r.Tag, r.ID).Count(&count)
		if count > 0 {
			return errors.New("Tag重复")
		}
		if r.Tag == global.SuperAdminUserTag {
			return errors.New("超级管理员角色禁止编辑")
		}
		global.Db.Unscoped().Where("role_id=?", r.ID).Delete(&RoleMenu{}) // 删除原来的关联
		if err := global.Db.Model(r).Save(r).Error; err != nil {
			return err
		}
	}
	return nil
}

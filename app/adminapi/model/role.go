/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:28:08
 * @LastEditTime: 2022-09-29 15:29:24
 */
package model

import (
	"errors"
	"iris-project/app/dao"
	"iris-project/global"

	"gorm.io/gorm"
)

// Role 角色model
type Role struct {
	// gorm.Model
	ID        uint32           `gorm:"primaryKey;" json:"id"`
	CreatedAt global.LocalTime `gorm:"type:datetime;" json:"created_at"`
	Name      string           `gorm:"type:varchar(50);not null" json:"name"`
	JumpPage  string           `gorm:"type:varchar(255)" json:"jump_page"` // 登录后跳转页面
	Tag       string           `gorm:"type:varchar(50);unique" json:"tag"`
	Menus     []*Menu          `gorm:"many2many:role_menu;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;association_autoupdate:false" json:"-"`
	MenusTree []*MenuTree      `gorm:"-" json:"menus_tree,omitempty"`
	MenuNames string           `gorm:"-" json:"menu_names"` // 菜单名，用 逗号 分隔
	MenuIDS   []uint32         `gorm:"-" json:"menu_ids"`   // 菜单IDS
	Status    int8             `gorm:"type:tinyint(1);default:1" json:"status"`
}

// GetOne 获取角色（不包含菜单）
func (r *Role) GetOne(load bool) bool {
	if r.ID == 0 {
		return false
	}
	if err := dao.GetDB().First(r).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

// GetRoleMenus 获取角色（包含菜单）
func (r *Role) GetRoleMenus() bool {
	if r.ID == 0 {
		return false
	}
	if err := dao.GetDB().Preload("Menus").First(r).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

// GetRoleMenusTree 获取角色（包含菜单和菜单数）
func (r *Role) GetRoleMenusTree() bool {
	if r.ID == 0 {
		return false
	}
	if err := dao.GetDB().Preload("Menus").First(r).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	r.MenusTree = GetTreeMenus(r.Menus)
	return true
}

// CreateUpdate 创建或更新角色
func (r *Role) CreateUpdate() error {
	if r.ID == 0 {
		var count int64
		dao.GetDB().Model(&Role{}).Where("tag=?", r.Tag).Count(&count)
		if count > 0 {
			return errors.New("Tag重复")
		}
		if err := dao.GetDB().Create(r).Error; err != nil {
			return err
		}
	} else {
		var count int64
		dao.GetDB().Model(&Role{}).Where("tag=? and id<>?", r.Tag, r.ID).Count(&count)
		if count > 0 {
			return errors.New("Tag重复")
		}
		if r.Tag == global.SuperAdminUserTag {
			return errors.New("超级管理员角色禁止编辑")
		}
		dao.GetDB().Unscoped().Where("role_id=?", r.ID).Delete(&RoleMenu{}) // 删除原来的关联
		if err := dao.GetDB().Model(r).Save(r).Error; err != nil {
			return err
		}
	}
	return nil
}

// Delete 删除角色
func (r *Role) Delete() error {
	if r.ID == 0 {
		return errors.New("需指定ID")
	}
	if err := dao.GetDB().Unscoped().Delete(r).Error; err != nil {
		return err
	}
	return nil
}

// Updates 更新角色
func (r *Role) Updates(data map[string]interface{}) error {
	if r.ID == 0 {
		return errors.New("需指定ID")
	}
	if err := dao.GetDB().Model(r).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

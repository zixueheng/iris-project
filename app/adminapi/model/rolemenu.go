package model

// RoleMenu 角色菜单关联
type RoleMenu struct {
	RoleID uint `gorm:"" json:"role_id"`
	MenuID uint `gorm:"" json:"menu_id"`
}

package model

// AdminUserRole 管理员角色关联
type AdminUserRole struct {
	AdminUserID uint `gorm:"" json:"admin_user_id"`
	RoleID      uint `gorm:"" json:"role_id"`
}

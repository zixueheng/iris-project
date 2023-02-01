/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:28:08
 * @LastEditTime: 2021-04-28 14:59:31
 */
package model

// AdminUserRole 管理员角色关联
type AdminUserRole struct {
	AdminUserID uint32 `gorm:"type:int;" json:"admin_user_id"`
	RoleID      uint32 `gorm:"type:int;" json:"role_id"`
}

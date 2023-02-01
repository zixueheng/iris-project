/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:28:08
 * @LastEditTime: 2021-04-28 14:59:44
 */
package model

// RoleMenu 角色菜单关联
type RoleMenu struct {
	RoleID uint32 `gorm:"type:int;" json:"role_id"`
	MenuID uint32 `gorm:"type:int;" json:"menu_id"`
}

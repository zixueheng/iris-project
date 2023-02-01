/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:28:08
 * @LastEditTime: 2021-10-12 17:39:11
 */
package validate

// RoleRequest 角色请求格式
type RoleRequest struct {
	ID       uint32   `json:"id" validate:"numeric" comment:"ID"`
	Name     string   `json:"name" validate:"required" comment:"名称"`
	Tag      string   `json:"tag" validate:"required" comment:"标识"`
	JumpPage string   `json:"jump_page" validate:"-" comment:"跳转页面"` // 登录后跳转页面
	Status   int8     `json:"status" validate:"numeric" comment:"状态"`
	MenuIds  []uint32 `json:"menu_ids" validate:"required" comment:"菜单IDS"`
}

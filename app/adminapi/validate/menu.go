/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:28:08
 * @LastEditTime: 2021-08-09 16:32:30
 */
package validate

// MenuRequest 菜单请求格式
type MenuRequest struct {
	ID            uint32 `json:"id" validate:"numeric" comment:"ID"`
	PID           uint32 `json:"p_id" validate:"gte=0" comment:"父菜单"`
	Name          string `json:"name" validate:"required" comment:"名称"`
	Icon          string `json:"icon" validate:"-" comment:"图标"`
	Type          string `json:"type" validate:"required" comment:"类型"`
	MenuPath      string `json:"menu_path" validate:"-" comment:"菜单路径"`
	UniqueAuthKey string `json:"unique_auth_key" validate:"-" comment:"前端鉴权key"` // 前端鉴权key
	Action        string `json:"action" validate:"-" comment:"前端鉴权key"`          // 前端鉴权key
	Subject       string `json:"subject" validate:"-" comment:"前端鉴权key"`         // 前端鉴权key
	APIPath       string `json:"api_path" validate:"-" comment:"接口路径"`
	Method        string `json:"method" validate:"-" comment:"请求方式"`
	Header        string `json:"header" validate:"-" comment:"顶部菜单标示"`
	IsHeader      int8   `json:"is_header" validate:"numeric" comment:"是否顶部菜单1是0否"`
	Sort          uint   `json:"sort" validate:"numeric"  comment:"排序"`
	Status        int8   `json:"status" validate:"numeric" comment:"状态"`
}

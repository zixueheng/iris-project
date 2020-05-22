package validate

// MenuRequest 菜单请求格式
type MenuRequest struct {
	ID       uint   `json:"id" comment:"ID"`
	PID      uint   `json:"p_id" validate:"required" comment:"父菜单"`
	Name     string `json:"name" validate:"required" comment:"名称"`
	Icon     string `json:"icon" comment:"图标"`
	Type     string `json:"type" validate:"required" comment:"类型"`
	MenuPath string `json:"menu_path" comment:"菜单路径"`
	APIPath  string `json:"api_path" comment:"接口路径"`
	Method   string `json:"method" comment:"请求方式"`
	Sort     uint   `json:"sort" validate:"numeric"  comment:"排序"`
	Status   int8   `json:"status" comment:"状态"`
}

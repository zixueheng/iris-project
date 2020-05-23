package validate

// RoleRequest 角色请求格式
type RoleRequest struct {
	ID      uint   `json:"id" comment:"ID"`
	Name    string `json:"name" validate:"required,gte=4,lte=50" comment:"名称"`
	Tag     string `json:"tag" validate:"required" comment:"标识"`
	Status  int8   `json:"status" comment:"状态"`
	MenuIds []uint `json:"menu_ids" validate:"required" comment:"菜单IDS"`
}

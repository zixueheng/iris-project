package validate

// RoleRequest 角色请求格式
type RoleRequest struct {
	Name    string `json:"name" validate:"required,gte=4,lte=50" comment:"名称"`
	Tag     string `json:"tag" comment:"标识"`
	Status  int    `json:"status" comment:"描述"`
	MenuIds []uint `json:"menu_ids" comment:"权限"`
}

package validate

// MenuRequest 菜单请求格式
type MenuRequest struct {
	Name    string `json:"name" validate:"required,gte=4,lte=50" comment:"名称"`
	PID     uint   `json:"p_id" validate:"required" comment:"父菜单"`
	APIPath string `json:"api_path" validate:"required" comment:"路径"`
	Status  int8   `json:"act" comment:"状态"`
}

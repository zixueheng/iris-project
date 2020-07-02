package validate

// CreateUpdateAdminUserRequest 创建更新管理员请求格式
type CreateUpdateAdminUserRequest struct {
	ID       uint   `json:"id" validate:"numeric" comment:"ID"`
	Username string `json:"username" validate:"required,gte=2,lte=50" comment:"用户名"`
	Realname string `json:"realname" validate:"required,gte=2,lte=50" comment:"姓名"`
	Password string `json:"password" validate:"omitempty,gte=4,lte=50"  comment:"密码"`
	Phone    string `json:"phone" validate:"required,len=11" comment:"手机号"`
	RoleIds  []uint `json:"role_ids" validate:"required" comment:"角色IDS"`
	Status   int8   `json:"status" validate:"numeric" comment:"状态"`
}

// LoginRequest 管理员登录请求格式
type LoginRequest struct {
	Username string `json:"username" validate:"required,gte=2,lte=50" comment:"用户名"`
	Password string `json:"password" validate:"required" comment:"密码"`
}

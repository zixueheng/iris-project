package global

// 系统常量

const (
	// AdminTokenMinutes 管理员Token有效期/分钟
	AdminTokenMinutes = 6 * 60
	// AdminRefreshTokenMinutes 管理员RefreshToken 有效期/分钟
	AdminRefreshTokenMinutes = 24 * 60
	// AdminUserCacheMinutes 管理员信息缓存时间/分钟
	AdminUserCacheMinutes = 60

	// AdminUserJWTKey JWT荷载playload中要保存的键名
	AdminUserJWTKey = "admin_user_id"

	// SuperAdminUserTag 超级管理员Tag
	SuperAdminUserTag = "superadmin"
)

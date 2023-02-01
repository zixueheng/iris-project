/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-18 09:34:00
 * @LastEditTime: 2022-11-03 14:09:55
 */
package global

// 系统常量

const (
	// WapUserJWTKey Wap端JWT荷载playload中要保存的键名
	WapUserJWTKey = "wap_user_id"
	// UserTokenMinutes 会员Token有效期/分钟
	UserTokenMinutes = 7 * 24 * 60
	// UserRefreshTokenMinutes 会员RefreshToken 有效期/分钟
	UserRefreshTokenMinutes = 30 * 24 * 60

	// StaffJWTKey 员工端JWT荷载playload中要保存的键名
	StaffJWTKey = "staff_id"
	// StaffTokenMinutes 员工Token有效期/分钟
	StaffTokenMinutes = 12 * 60
	// StaffRefreshTokenMinutes 员工RefreshToken 有效期/分钟
	StaffRefreshTokenMinutes = 1 * 24 * 60
	// StaffCacheMinutes 员工信息缓存时间/分钟
	StaffCacheMinutes = 12 * 60
	// StaffCacheKeyPrefix 员工信息缓存key前缀
	StaffCacheKeyPrefix = ":vo_staff_"

	// AdminTokenMinutes 管理员Token有效期/分钟
	AdminTokenMinutes = 12 * 60
	// AdminRefreshTokenMinutes 管理员RefreshToken 有效期/分钟
	AdminRefreshTokenMinutes = 24 * 60
	// AdminUserCacheMinutes 管理员信息缓存时间/分钟
	AdminUserCacheMinutes = 60
	// AdminUserCacheKeyPrefix 管理员信息缓存key前缀
	AdminUserCacheKeyPrefix = ":vo_admin_user_"

	// AdminUserJWTKey JWT荷载playload中要保存的键名
	AdminUserJWTKey = "admin_user_id"

	// SuperAdminUserTag 超级管理员Tag
	SuperAdminUserTag = "superadmin"

	// VerifycodeLifeTime 验证码有效分钟数
	VerifycodeLifeTime = 10
)

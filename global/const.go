/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-18 09:34:00
 * @LastEditTime: 2023-02-21 11:32:53
 */
package global

// 系统常量
const (
	// WapAPI wap端路由
	WapAPI = "/wapapi"
	// AdminAPI admin端路由
	AdminAPI = "/adminapi"
	// MerchantAPI merchant端路由
	MerchantAPI = "/merchantapi"

	UserID    = "user_id" // 不可修改
	ClientKey = "client"  // 不可修改

	// WapUserJWTKey Wap端JWT荷载playload中要保存的键名
	WapUserJWTKey = UserID //"wap_user_id"
	// UserTokenMinutes 会员Token有效期/分钟
	UserTokenMinutes = 7 * 24 * 60
	// UserRefreshTokenMinutes 会员RefreshToken 有效期/分钟
	UserRefreshTokenMinutes = 30 * 24 * 60

	/*
		// StaffJWTKey 员工端JWT荷载playload中要保存的键名
		StaffJWTKey = UserID // "staff_id"
		// StaffTokenMinutes 员工Token有效期/分钟
		StaffTokenMinutes = 12 * 60
		// StaffRefreshTokenMinutes 员工RefreshToken 有效期/分钟
		StaffRefreshTokenMinutes = 1 * 24 * 60
		// StaffCacheMinutes 员工信息缓存时间/分钟
		StaffCacheMinutes = 12 * 60
		// StaffCacheKeyPrefix 员工信息缓存key前缀
		StaffCacheKeyPrefix = ":vo_staff_"
	*/

	// MerchantJWTKey 商家端JWT荷载playload中要保存的键名
	MerchantJWTKey = UserID // "merchant_id"
	// MerchantTokenMinutes 商家Token有效期/分钟
	MerchantTokenMinutes = 12 * 60
	// MerchantRefreshTokenMinutes 商家RefreshToken 有效期/分钟
	MerchantRefreshTokenMinutes = 7 * 24 * 60

	// AdminTokenMinutes 管理员Token有效期/分钟
	AdminTokenMinutes = 12 * 60
	// AdminRefreshTokenMinutes 管理员RefreshToken 有效期/分钟
	AdminRefreshTokenMinutes = 24 * 60
	// AdminUserCacheMinutes 管理员信息缓存时间/分钟
	AdminUserCacheMinutes = 60
	// AdminUserCacheKeyPrefix 管理员信息缓存key前缀
	AdminUserCacheKeyPrefix = ":vo_admin_user_"

	// AdminUserJWTKey JWT荷载playload中要保存的键名
	AdminUserJWTKey = UserID //"admin_user_id"

	// SuperAdminUserTag 超级管理员Tag
	SuperAdminUserTag = "superadmin"

	// VerifycodeLifeTime 验证码有效分钟数
	VerifycodeLifeTime = 10
)

var (
	ClientMap = map[string]string{
		WapAPI:      "wap",
		AdminAPI:    "admin",
		MerchantAPI: "merchant",
	}
)

// GetClient 获取端
func GetClient(API string) string {
	if client, ok := ClientMap[API]; ok {
		return client
	}
	return ""
}

/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-04-12 14:01:05
 * @LastEditTime: 2021-11-24 15:03:07
 */
package model

import (
	"iris-project/global"
)

// User 模型
type User struct {
	ID        uint32           `gorm:"primaryKey;" json:"id"`
	CreatedAt global.LocalTime `gorm:"type:datetime;comment:创建时间;" json:"created_at,omitempty" validate:"-"`
	Realname  string           `gorm:"type:varchar(100);comment:姓名;" json:"realname" validate:"-" comment:"姓名"`
	Nickname  string           `gorm:"type:varchar(100);comment:昵称;" json:"nickname" validate:"-" comment:"昵称"`
	Avatar    string           `gorm:"type:varchar(255);comment:头像;" json:"avatar" validate:"-" comment:"头像"`
	Phone     string           `gorm:"type:char(11);comment:手机号;" json:"phone" validate:"-" comment:"手机号"`
	Birthday  global.LocalDate `gorm:"type:date;comment:生日;" json:"birthday" validate:"-" comment:"生日"`
	Gender    int8             `gorm:"type:tinyint(1);default:0;comment:性别;" json:"gender" validate:"numeric,oneof=0 1 2" comment:"性别"` // 1男 2女 0未知
	Source    string           `gorm:"type:enum(\"miniprogram\",\"android\",\"ios\",\"\");default:\"\";comment:注册来源;" validate:"-" json:"source" comment:"注册来源"`
	Status    int8             `gorm:"type:tinyint(1);default:1;comment:状态;" json:"status" validate:"numeric,oneof=1 -1" comment:"状态"` // 1显示 -1隐藏
}

// GetID 返回ID，实现接口`dao.Model`的方法
func (m *User) GetID() uint32 {
	return m.ID
}

// LoadRelatedField 加载相关字段
func (m *User) LoadRelatedField() {
}

// UserRefreshToken 模型
type UserRefreshToken struct {
	// ID                  uint32 `gorm:"primaryKey;" json:"id"`

	UserID              uint32 `gorm:"index;" json:"user_id"`
	RefreshToken        string `gorm:"type:varchar(255);" json:"-" validate:"-" comment:"刷新TOKEN"`
	RefreshTokenExpired uint32 `gorm:"" json:"-" validate:"-" comment:"刷新TOKEN过期时间戳"`
}

// GetID 返回ID，实现接口`dao.Model`的方法
// func (m *UserRefreshToken) GetID() uint32 {
// 	return m.ID
// }

// UserWechat 模型
type UserWechat struct {
	ID     uint32 `gorm:"primaryKey;" json:"id"`
	UserID uint32 `gorm:"unique;" json:"user_id"`
	// Openid  string `gorm:"index;type:varchar(200);" json:"openid" validate:"-" comment:"Openid"`
	MinOpenid string `gorm:"index;type:varchar(200);" json:"min_openid" validate:"-" comment:"MinOpenid"`
	AppOpenid string `gorm:"index;type:varchar(200);" json:"app_openid" validate:"-" comment:"AppOpenid"`
	Unionid   string `gorm:"type:varchar(200);" json:"unionid" validate:"-" comment:"Unionid"`
	// User    User   `json:"-"`
}

// GetID 返回ID，实现接口`dao.Model`的方法
func (m *UserWechat) GetID() uint32 {
	return m.ID
}

// UserUsername 模型
type UserUsername struct {
	ID       uint32 `gorm:"primaryKey;" json:"id"`
	UserID   uint32 `gorm:"unique;" json:"user_id"`
	Username string `gorm:"unique;type:varchar(100);" json:"username" validate:"required,gte=2,lte=50" comment:"用户名"`
	Password string `gorm:"type:varchar(100);" json:"-" validate:"required" comment:"密码"`
	// User    User   `json:"-"`
}

// GetID 返回ID，实现接口`dao.Model`的方法
func (m *UserUsername) GetID() uint32 {
	return m.ID
}

// UserPhone 模型
type UserPhone struct {
	ID       uint32 `gorm:"primaryKey;" json:"id"`
	UserID   uint32 `gorm:"unique;" json:"user_id"`
	Phone    string `gorm:"unique;type:char(11);" json:"phone" validate:"required,len=11" comment:"手机号"`
	Password string `gorm:"type:varchar(100);" json:"-" validate:"required" comment:"密码"`
	// User    User   `json:"-"`
}

// GetID 返回ID，实现接口`dao.Model`的方法
func (m *UserPhone) GetID() uint32 {
	return m.ID
}

// Verifycode 验证码
type Verifycode struct {
	Phone   string           `gorm:"index;type:char(11);" json:"phone" validate:"required,len=11" comment:"手机号"`
	Code    string           `gorm:"type:varchar(10);" json:"code" validate:"required" comment:"验证码"`
	Expired global.LocalTime `gorm:"type:datetime;" json:"expired" validate:"-" comment:"过期时间"`
	Purpose string           `gorm:"type:varchar(50);" json:"purpose" validate:"-" comment:"用途"`
}

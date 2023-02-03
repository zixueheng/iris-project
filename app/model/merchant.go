/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-03-23 10:13:01
 * @LastEditTime: 2022-12-15 15:08:15
 */
package model

import (
	"encoding/json"
	"iris-project/app/dao"
	"iris-project/global"
)

// Merchant 商户模型
type Merchant struct {
	ID        uint32           `gorm:"primaryKey;" json:"id"`
	CreatedAt global.LocalTime `gorm:"type:datetime;comment:创建时间;" json:"created_at,omitempty" validate:"-"`
	Name      string           `gorm:"type:varchar(100);not null;comment:商户名称;" json:"name" validate:"required" comment:"商户名称"`
	Nickname  string           `gorm:"type:varchar(100);comment:商户别名;" json:"nickname" validate:"-" comment:"商户别名"`
	Username  string           `gorm:"type:varchar(100);not null;unique;comment:用户名;" json:"username" validate:"required" comment:"用户名"`
	Password  string           `gorm:"type:varchar(100);not null;comment:密码;" json:"-" validate:"-" comment:"密码"`
	Contacts  string           `gorm:"type:varchar(100);comment:联系人;" json:"contacts" validate:"-" comment:"联系人"`
	Phone     string           `gorm:"type:char(11);comment:手机号;" json:"phone" validate:"-" comment:"手机号"`
	Address   string           `gorm:"type:varchar(255);comment:地址;" json:"address" validate:"-" comment:"地址"`
	Remark    string           `gorm:"type:text;comment:备注;" json:"remark" validate:"-" comment:"备注"`
	Pics      string           `gorm:"type:text;comment:图片;" json:"-" validate:"-" comment:"图片"`
	PicsJSON  []string         `gorm:"-" json:"pics" validate:"required" comment:"图片"`
	Status    int8             `gorm:"type:tinyint(1);default:1;comment:状态;" json:"status" validate:"numeric,oneof=1 -1" comment:"状态"` // 1显示 -1隐藏
}

// GetID 返回ID，实现接口`dao.Model`的方法
func (m *Merchant) GetID() uint32 {
	return m.ID
}

// LoadRelatedField 加载相关字段
func (m *Merchant) LoadRelatedField() {
	if m.Pics != "" {
		json.Unmarshal([]byte(m.Pics), &m.PicsJSON)
	}
}

// MerchantUser 商户店员模型
type MerchantUser struct {
	ID           uint32           `gorm:"primaryKey;" json:"id"`
	CreatedAt    global.LocalTime `gorm:"type:datetime;comment:创建时间;" json:"created_at,omitempty" validate:"-"`
	MerchantID   uint32           `gorm:"index;comment:商户ID;" json:"merchant_id" validate:"required" comment:"商户ID"`
	MerchantName string           `gorm:"-" json:"merchant_name" validate:"-" comment:"商户名称"`
	UserID       uint32           `gorm:"index;comment:店员ID;" json:"user_id" validate:"required" comment:"店员ID"`
	UserRealname string           `gorm:"-" json:"user_realname" validate:"-" comment:"姓名"`
	UserNickname string           `gorm:"-" json:"user_nickname" validate:"-" comment:"昵称"`
	UserAvatar   string           `gorm:"-" json:"user_avatar" validate:"-" comment:"头像"`
	UserPhone    string           `gorm:"-" json:"user_phone" validate:"-" comment:"手机号"`
}

// GetID 返回ID，实现接口`dao.Model`的方法
func (m *MerchantUser) GetID() uint32 {
	return m.ID
}

// LoadRelatedField 加载相关字段
func (m *MerchantUser) LoadRelatedField() {
	if m.MerchantID != 0 {
		merchant := struct{ Name string }{}
		dao.Scan(nil, &Merchant{}, map[string]interface{}{"id": m.MerchantID}, &merchant, []string{"name"})
		m.MerchantName = merchant.Name
	}
	if m.UserID != 0 {
		user := struct {
			ID                                uint32
			Realname, Nickname, Avatar, Phone string
		}{}
		dao.Scan(nil, &User{}, map[string]interface{}{"id": m.UserID}, &user, []string{"id", "realname", "nickname", "avatar", "phone"})
		if user.ID != 0 {
			m.UserRealname = user.Realname
			m.UserNickname = user.Nickname
			m.UserAvatar = user.Avatar
			m.UserPhone = user.Phone
		}
	}
}

// MerchantPayInfo 商家收款账号
type MerchantPayInfo struct {
	// ID           uint32           `gorm:"primaryKey;" json:"id"`
	// CreatedAt    global.LocalTime `gorm:"type:datetime;comment:创建时间;" json:"created_at,omitempty" validate:"-"`
	MerchantID uint32 `gorm:"index;comment:商户ID;" json:"merchant_id" validate:"required" comment:"商户ID"`
	// MerchantName string      `gorm:"-" json:"merchant_name" validate:"-" comment:"商户名称"`
	Channel     string      `gorm:"type:varchar(100);comment:收款渠道;" json:"channel" validate:"required,oneof=wechat alipay" comment:"收款渠道"`
	PayInfo     string      `gorm:"type:text;comment:支付账号;" json:"-" validate:"-" comment:"支付账号"`
	PayInfoJSON interface{} `gorm:"-" json:"pay_info" validate:"-" comment:"支付账号"`
}

// GetID 返回ID，实现接口`dao.Model`的方法
// func (m *MerchantPayInfo) GetID() uint32 {
// 	return m.ID
// }

// LoadRelatedField 加载相关字段
func (m *MerchantPayInfo) LoadRelatedField() {
	if m.Channel == "wechat" && m.PayInfo != "" {
		json.Unmarshal([]byte(m.PayInfo), &m.PayInfoJSON)
	}
}

// WechatPayInfo 微信支付信息
type WechatPayInfo struct {
	MchId         string `json:"mch_id" validate:"required"`          // 商户号
	SerialNo      string `json:"serial_no" validate:"required"`       // 商户证书序列号
	Apiv3Key      string `json:"apiv3_key" validate:"required"`       // 商户APIv3密钥
	ApiClientCert string `json:"api_client_cert" validate:"required"` // 公钥
	ApiClientKey  string `json:"api_client_key" validate:"required"`  // 私钥
	PlatformCert  string `json:"-" validate:"-"`                      // 平台公钥
}

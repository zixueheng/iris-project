/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-10-11 15:27:55
 * @LastEditTime: 2023-10-12 14:35:28
 */
package dao

import (
	"iris-project/global"
	"testing"
)

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

	Username string `json:"username" validate:"-" comment:"用户名"` // 关联表字段不能用 gorm:"-"
	Password string `json:"password" validate:"-" comment:"密码"`  // 关联表字段不能用 gorm:"-"
}

func Test_Find(t *testing.T) {
	var user = &User{}
	// FindOne2(nil, user, map[string]interface{}{"realname": "he"}, &QueryOpts{})
	// t.Logf("%+v", user)
	// FindOne2(nil, user, map[string]interface{}{"realname": "he"}, &QueryOpts{Select: []string{"id", "realname", "nickname", "phone"}, OrderBy: []string{"id desc"}})
	// t.Logf("%+v", user)

	var opts = &QueryOpts{
		Select:  []string{"iris_user.id", "iris_user.realname", "iris_user.nickname", "iris_user.phone", "iris_user_username.username", "iris_user_username.password"},
		OrderBy: []string{"iris_user.id desc"},
		Joins:   []string{"left join iris_user_username on iris_user_username.user_id = iris_user.id"},
	}
	FindOneOpts(nil, user, map[string]interface{}{"iris_user.realname": "he"}, opts)
	t.Logf("%+v\n", user)

	ScanOpts(nil, &User{}, map[string]interface{}{"iris_user.realname": "he"}, user, opts)
	t.Logf("%+v\n", user)
}

package model

import (
	"iris-project/global"
)

// User 模型
type User struct {
	ID        uint           `gorm:"primary_key"`
	CreatedAt global.SQLTime `gorm:"type:timestamp;"`
	Name      string
	Age       uint
	Birthday  string
}

// CreateUser 创建用户
func (u *User) CreateUser() error {
	if global.Db.NewRecord(u) {
		if err := global.Db.Create(&u).Error; err != nil {
			return err
		}
	}

	return nil
}

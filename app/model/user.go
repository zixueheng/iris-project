package model

import (
	"iris-project/global"

	"github.com/jinzhu/gorm"
)

// User 模型
type User struct {
	gorm.Model
	Name     string
	Age      uint
	Birthday string
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

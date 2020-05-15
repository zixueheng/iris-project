package model

import (
	"iris-project/global"

	"github.com/jinzhu/gorm"
)

// Menu 模型
type Menu struct {
	gorm.Model
	PID     uint   `gorm:"default:0"`
	Name    string `gorm:"type:varchar(50);not null"`
	APIPath string `gorm:"not null"`
	Status  int8   `gorm:"type:tinyint(1);default:1"`
}

// CreateMenu 创建菜单
func (m *Menu) CreateMenu() error {
	if global.Db.NewRecord(m) {
		if err := global.Db.Create(&m).Error; err != nil {
			return err
		}
	}

	return nil
}

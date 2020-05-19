package model

import (
	"iris-project/global"
	"time"
)

// Menu 模型
type Menu struct {
	// gorm.Model
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	PID       uint   `gorm:"default:0"`
	Name      string `gorm:"type:varchar(50);not null"`
	Type      string `gorm:"type:enum(\"menu\",\"api\");not null"`
	APIPath   string `gorm:"not null"`
	Method    string `gorm:"type:enum(\"GET\",\"POST\",\"PUT\",\"DELETE\");not null"`
	Sort      uint   `gorm:"default:0"`
	Status    int8   `gorm:"type:tinyint(1);default:1"`
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

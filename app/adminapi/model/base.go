package model

import (
	"iris-project/global"

	"github.com/jinzhu/gorm"
)

// FindOne 查找一个
func FindOne(obj interface{}, where map[string]interface{}, preload string) bool {
	tx := global.Db.Where(where)

	// 这里 Preload() 不起作用
	if len(preload) != 0 {
		tx.Preload(preload)
	}

	if err := tx.First(obj).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
}

// GetAll 获取列表
func GetAll(where map[string]interface{}, orderBy string, offset, limit uint) *gorm.DB {
	tx := global.Db
	if len(orderBy) > 0 {
		tx.Order(orderBy + "desc")
	} else {
		tx.Order("created_at desc")
	}
	if len(where) > 0 {
		tx.Where(where)
	}
	if offset > 0 {
		tx.Offset(offset)
	}
	if limit > 0 {
		tx.Limit(limit)
	}
	return tx
}

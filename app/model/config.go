/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:28:09
 * @LastEditTime: 2022-01-06 11:32:39
 */
package model

import (
	"encoding/json"
	"iris-project/global"
	"sort"
)

// Config 配置模型
type Config struct {
	ID          uint32           `gorm:"primaryKey;" json:"id"`
	CreatedAt   global.LocalTime `gorm:"type:datetime;" json:"created_at" validate:"-"`
	Code        string           `gorm:"type:varchar(50);not null;unique" json:"code" validate:"required" comment:"编码"`
	Type        string           `gorm:"type:varchar(100);not null" json:"type" validate:"required" comment:"类型"`
	Name        string           `gorm:"type:varchar(100);not null" json:"name" validate:"required" comment:"名称"`
	Content     string           `gorm:"type:text" json:"-" validate:"-" comment:"内容"`
	ContentJSON interface{}      `gorm:"-" json:"content" validate:"-" comment:"内容"`
	Status      int8             `gorm:"type:tinyint(1);default:1" json:"status" validate:"numeric,oneof=1 -1" comment:"状态"` // 1显示 -1隐藏
}

// GetID 返回ID，实现接口`dao.Model`的方法
func (m *Config) GetID() uint32 {
	return m.ID
}

// LoadRelatedField 加载相关字段
func (m *Config) LoadRelatedField() {
	switch m.Type {
	case "banner":
		var banners []ConfigBanner
		if m.Content != "" {
			json.Unmarshal([]byte(m.Content), &banners)
		} else {
			banners = make([]ConfigBanner, 0)
		}
		sort.SliceStable(banners, func(i, j int) bool {
			return banners[i].Sort < banners[j].Sort
		})
		m.ContentJSON = banners
	case "field":
		m.ContentJSON = m.Content
	case "array":
		var arr []interface{}
		if m.Content != "" {
			json.Unmarshal([]byte(m.Content), &arr)
		} else {
			arr = make([]interface{}, 0)
		}
		m.ContentJSON = arr
	case "object":
		var obj interface{}
		if m.Content != "" {
			json.Unmarshal([]byte(m.Content), &obj)
		} else {
			obj = nil
		}
		m.ContentJSON = obj
	default:
	}
}

// ConfigBanner Banner结构
type ConfigBanner struct {
	Pic   string `json:"pic" validate:"required" comment:"图片路径"`
	Title string `json:"title" validate:"-" comment:"标题"`
	MpUrl string `json:"mp_url" validate:"-" comment:"小程序路径"`
	// AndroidUrl string `json:"android_url" validate:"-" comment:"安卓路径"`
	// IosUrl     string `json:"ios_url" validate:"-" comment:"IOS路径"`
	Sort uint32 `json:"sort" validate:"-" comment:"排序"`
}

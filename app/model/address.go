/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-03-29 14:35:33
 * @LastEditTime: 2022-03-28 17:07:57
 */
package model

import (
	"iris-project/app/dao"
	"iris-project/global"
)

// Address 会员地址
type Address struct {
	ID        uint32           `gorm:"primaryKey;" json:"id"`
	CreatedAt global.LocalTime `gorm:"type:datetime;comment:创建时间;" json:"created_at,omitempty" validate:"-"`
	UserID    uint32           `gorm:"index;comment:用户ID;" json:"user_id" validate:"-" comment:"用户ID"`
	RealName  string           `gorm:"varchar(255);comment:姓名;" json:"real_name" validate:"required" comment:"姓名"`
	Phone     string           `gorm:"char(11);comment:电话;" json:"phone" validate:"len=11" comment:"电话"`
	Province  string           `gorm:"varchar(100);comment:省;" json:"province" validate:"required" comment:"省"`
	City      string           `gorm:"varchar(100);comment:市;" json:"city" validate:"required" comment:"市"`
	// CityID    uint32           `gorm:"comment:市ID;" json:"city_id" validate:"required" comment:"市ID"`
	District  string `gorm:"varchar(100);comment:区;" json:"district" validate:"required" comment:"区"`
	Detail    string `gorm:"varchar(500);comment:详细地址;" json:"detail" validate:"required" comment:"详细地址"`
	IsDefault int8   `gorm:"type:tinyint(1);default:-1;comment:是否默认;" json:"is_default" validate:"numeric,oneof=1 -1" comment:"是否默认"`
}

// GetID 返回ID，实现接口`dao.Model`的方法
func (m *Address) GetID() uint32 {
	return m.ID

}

// LoadRelatedField 加载相关字段
func (m *Address) LoadRelatedField() {

}

// City 省市区
type City struct {
	ID       uint32 `gorm:"primaryKey;" json:"id"`
	CityID   uint32 `gorm:"comment:城市ID;" json:"city_id" validate:"required" comment:"城市ID"`
	Level    uint32 `gorm:"comment:级别;" json:"level" validate:"-" comment:"级别"`
	ParentID uint32 `gorm:"comment:父ID;" json:"parent_id" validate:"-" comment:"父ID"`
	Name     string `gorm:"varchar(100);comment:名称;" json:"name" validate:"required" comment:"名称"`
	IsShow   int8   `gorm:"type:tinyint(1);default:1;comment:显示;" json:"is_show" validate:"numeric,oneof=1 -1" comment:"显示"`
}

// CityTree 省市区树
type CityTree struct {
	*City
	Children []*CityTree `json:"children"`
}

// GetCityTree 获取省市区树
func GetCityTree() []*CityTree {
	var (
		trees    []*CityTree
		rootTree = &CityTree{
			City: &City{
				ID:   0,
				Name: "根节点",
			},
		}
		allCities []*City
	)
	dao.FindAll(nil, &allCities, nil)
	for _, v := range allCities {
		trees = append(trees, &CityTree{City: v})
	}
	makeCityTree(trees, rootTree)
	return rootTree.Children
}

func makeCityTree(list []*CityTree, node *CityTree) {
	if children, has := hasCityChild(list, node); has {
		node.Children = append(node.Children, children...)
		for _, v := range children {
			makeCityTree(list, v)
		}
	} else {
		node.Children = make([]*CityTree, 0) // 没有子节点将 nil 转成空切片，输出json 是[] 而不是 null
	}
}

func hasCityChild(list []*CityTree, node *CityTree) (children []*CityTree, has bool) {
	for _, m := range list {
		if m.ParentID == node.CityID {
			children = append(children, m)
		}
	}
	if len(children) > 0 {
		has = true
	}
	return
}

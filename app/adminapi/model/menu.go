/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:28:08
 * @LastEditTime: 2022-09-29 15:29:14
 */
package model

import (
	"errors"
	"iris-project/app/dao"
	"iris-project/global"

	"gorm.io/gorm"
)

// Menu 模型，包含菜单接口两种类型
type Menu struct {
	// gorm.Model
	ID            uint32           `gorm:"primaryKey;" json:"id"`
	CreatedAt     global.LocalTime `gorm:"type:datetime;" json:"created_at"`
	PID           uint32           `gorm:"default:0;type:int(10) unsigned;" json:"p_id"`
	Name          string           `gorm:"type:varchar(50);not null" json:"title"`
	Icon          string           `gorm:"type:varchar(50);" json:"icon"`
	Type          string           `gorm:"type:enum(\"menu\",\"api\",\"module\");not null" json:"type"`
	MenuPath      string           `gorm:"type:varchar(255);" json:"menu_path"`              // 前端菜单路径
	UniqueAuthKey string           `gorm:"unique;type:varchar(255);" json:"unique_auth_key"` // 前端鉴权key
	Action        string           `gorm:"type:varchar(255);" json:"action"`                 // 前端鉴权key
	Subject       string           `gorm:"type:varchar(255);" json:"subject"`                // 前端鉴权key
	APIPath       string           `gorm:"type:varchar(255);" json:"api_path"`               // 接口路径
	Method        string           `gorm:"type:enum(\"GET\",\"POST\",\"PUT\",\"DELETE\",\"\");default:null" json:"method"`
	Header        string           `gorm:"type:varchar(50);" json:"header"`
	IsHeader      int8             `gorm:"type:tinyint(1);default:0" json:"is_header"`
	Sort          uint             `gorm:"default:0;type:int(10);" json:"sort"`
	Status        int8             `gorm:"type:tinyint(1);default:1" json:"status"` // 1显示 -1隐藏
	HTML          string           `gorm:"-" json:"html"`                           // 用来输出层级 |----
	Level         int              `gorm:"-" json:"-"`                              // 计算层级
	Checked       bool             `gorm:"-" json:"checked"`                        // 是否选中，角色接口中用
	Selected      bool             `gorm:"-" json:"selected"`
	Expand        bool             `gorm:"-" json:"expand"`
}

// MenuTree 菜单树
type MenuTree struct {
	*Menu
	Children []*MenuTree `json:"children"`
}

// GetOne 查找菜单
func (m *Menu) GetOne(load bool) bool {
	if m.ID == 0 {
		return false
	}
	if err := dao.GetDB().First(m).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

// CreateUpdate 创建或更新菜单
func (m *Menu) CreateUpdate() error {
	if m.ID == 0 {
		var count int64
		dao.GetDB().Model(&Menu{}).Where("unique_auth_key=?", m.UniqueAuthKey).Count(&count)
		if count > 0 {
			return errors.New("鉴权Key重复")
		}
		if err := dao.GetDB().Create(m).Error; err != nil {
			return err
		}
	} else {
		var count int64
		dao.GetDB().Model(&Menu{}).Where("unique_auth_key=? and id<>?", m.UniqueAuthKey, m.ID).Count(&count)
		if count > 0 {
			return errors.New("鉴权Key重复")
		}
		if err := dao.GetDB().Model(m).Omit("created_at").Save(m).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetTreeMenus 获取树状菜单
//
// 如果参数 allMenus 不为空，则肯定包含1级菜单（即PID=0的项目）（每个角色的菜单权限是从1级菜单往下衍生的）
func GetTreeMenus(allMenus []*Menu) []*MenuTree {
	if len(allMenus) == 0 { // 不传入角色的菜单权限，则加载全部菜单权限
		dao.GetDB().Order("sort asc").Find(&allMenus)
	}

	var allMenuTree []*MenuTree

	rootMenuTree := &MenuTree{
		Menu: &Menu{ID: 0, Name: "根节点"},
	}
	// allMenuTree = append(allMenuTree, rootMenuTree)
	for _, m := range allMenus {
		mt := &MenuTree{
			Menu: m,
		}
		allMenuTree = append(allMenuTree, mt)
	}
	makeMenuTree(allMenuTree, rootMenuTree)
	return rootMenuTree.Children
}

func makeMenuTree(list []*MenuTree, menuTree *MenuTree) {
	if children, has := hasMenuChild(list, menuTree); has {
		menuTree.Children = append(menuTree.Children, children...)
		for _, v := range children {
			makeMenuTree(list, v)
		}
	} else {
		menuTree.Children = make([]*MenuTree, 0) // 没有子节点将 nil 转成空切片，输出json 是[] 而不是 null
	}
}

func hasMenuChild(list []*MenuTree, node *MenuTree) (children []*MenuTree, has bool) {
	for _, m := range list {
		if m.PID == node.ID {
			children = append(children, m)
		}
	}
	if len(children) > 0 {
		has = true
	}
	return
}

// Delete 删除菜单（包括子菜单）
func (m *Menu) Delete() error {
	if m.ID == 0 {
		return errors.New("需指定ID")
	}
	var delIDs []uint32
	delIDs = append(delIDs, m.ID)
	getChild(&delIDs, m.ID) // 查找所有子菜单

	if err := dao.GetDB().Unscoped().Where("id IN (?)", delIDs).Delete(&Menu{}).Error; err != nil {
		return err
	}
	return nil
}

func hasChild(id uint32) (ids []uint32, has bool) {
	dao.GetDB().Model(&Menu{}).Where("p_id=?", id).Pluck("id", &ids)
	if len(ids) > 0 {
		has = true
	}
	return
}

func getChild(allIDs *[]uint32, id uint32) {
	if cids, has := hasChild(id); has {
		*allIDs = append(*allIDs, cids...)
		for _, cid := range cids {
			getChild(allIDs, cid)
		}
	}
}

// Updates 更新菜单状态
func (m *Menu) Updates(data map[string]interface{}) error {
	if m.ID == 0 {
		return errors.New("需指定ID")
	}
	if err := dao.GetDB().Model(m).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

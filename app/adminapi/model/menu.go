package model

import (
	"iris-project/global"

	"github.com/jinzhu/gorm"
)

// Menu 模型，包含菜单接口两种类型
type Menu struct {
	// gorm.Model
	ID            uint           `gorm:"primary_key" json:"id"`
	CreatedAt     global.SQLTime `gorm:"type:datetime;" json:"created_at"`
	PID           uint           `gorm:"default:0" json:"p_id"`
	Name          string         `gorm:"type:varchar(50);not null" json:"name"`
	Icon          string         `gorm:"type:varchar(50);" json:"icon"`
	Type          string         `gorm:"type:enum(\"menu\",\"api\");not null" json:"type"`
	MenuPath      string         `gorm:"" json:"menu_path"`             // 前端菜单路径
	UniqueAuthKey string         `gorm:"unique" json:"unique_auth_key"` // 前端鉴权key
	APIPath       string         `gorm:"" json:"api_path"`              // 接口路径
	Method        string         `gorm:"type:enum(\"GET\",\"POST\",\"PUT\",\"DELETE\");default:null" json:"method"`
	Sort          uint           `gorm:"default:0" json:"sort"`
	Status        int8           `gorm:"type:tinyint(1);default:1" json:"status"` // 1显示 -1隐藏
}

// MenuTree 菜单树
type MenuTree struct {
	*Menu
	Children []*MenuTree `json:"children"`
}

// GetMenuByID 根据ID查找菜单
func (m *Menu) GetMenuByID(id uint) bool {
	if err := global.Db.Where("id=?", id).First(m).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
}

// CreateUpdateMenu 创建或更新菜单
func (m *Menu) CreateUpdateMenu() error {
	if m.ID == 0 {
		if err := global.Db.Create(m).Error; err != nil {
			return err
		}
	} else {
		if err := global.Db.Model(m).Save(m).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetTreeMenus 获取树状菜单
//
// 如果参数 allMenus 不为空，则肯定包含1级菜单（即PID=0的项目）（每个角色的菜单权限是从1级菜单往下衍生的）
func GetTreeMenus(allMenus []*Menu) []*MenuTree {
	if allMenus == nil || len(allMenus) == 0 { // 不传入角色的菜单权限，则加载全部菜单权限
		global.Db.Find(&allMenus)
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
	makeTree(allMenuTree, rootMenuTree)
	return rootMenuTree.Children
}

func makeTree(list []*MenuTree, menuTree *MenuTree) {
	if children, has := hasChild(list, menuTree); has {
		menuTree.Children = append(menuTree.Children, children...)
		for _, v := range children {
			makeTree(list, v)
		}
	} else {
		menuTree.Children = make([]*MenuTree, 0) // 没有子节点将 nil 转成空切片，输出json 是[] 而不是 null
	}
}

func hasChild(list []*MenuTree, node *MenuTree) (children []*MenuTree, has bool) {
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

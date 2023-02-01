/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-03-11 15:52:58
 * @LastEditTime: 2023-02-01 16:22:47
 */
package model

import (
	"iris-project/app/config"
	"iris-project/app/dao"
	"iris-project/global"
	"iris-project/lib/util"
	"strings"
)

// File 文件模型
type File struct {
	ID         uint32           `gorm:"primaryKey;" json:"id"`
	CreatedAt  global.LocalTime `gorm:"type:datetime;comment:创建时间;" json:"created_at,omitempty" validate:"-"`
	Name       string           `gorm:"index;type:varchar(255);not null;comment:文件名;" json:"name" validate:"required" comment:"文件名"`
	FileName   string           `gorm:"type:varchar(255);not null;comment:原文件名;" json:"filename" validate:"required" comment:"原文件名"`
	Path       string           `gorm:"type:varchar(500);not null;comment:路径;" json:"path" validate:"required" comment:"路径"`
	URL        string           `gorm:"-" json:"full_path" validate:"-" comment:"完整路径"`
	Size       int64            `gorm:"default:0;comment:大小;" json:"size" validate:"-" comment:"大小"`
	Type       string           `gorm:"type:enum(\"image\",\"doc\",\"audio\",\"video\",\"\");default:\"\";comment:类型;" json:"type" validate:"-" comment:"类型"`
	FileMime   string           `gorm:"type:varchar(255);comment:mime类型;" json:"file_mime" validate:"-" comment:"mime类型"`
	CategoryID uint32           `gorm:"default:0;comment:分类ID;" json:"category_id" validate:"-" comment:"分类ID"`
	Storage    string           `gorm:"type:enum(\"local\",\"qiniu\",\"tencent\",\"aliyun\",\"\");default:\"local\";comment:存储介质;" validate:"-" json:"storage" comment:"存储介质"`
	UploadBy   string           `gorm:"type:enum(\"user\",\"supplier\",\"admin\",\"\");default:\"\";comment:上传身份;" validate:"-" json:"upload_by" comment:"上传身份"`
	UploadByID uint32           `gorm:"default:0;comment:上传者ID;" json:"upload_by_id" validate:"-" comment:"上传者ID"`
	IsFavor    int8             `gorm:"-" json:"is_favor" validate:"numeric,oneof=1 -1" comment:"收藏状态"`
}

// GetID 返回ID，实现接口`dao.Model`的方法
func (m *File) GetID() uint32 {
	return m.ID
}

// FileFavor 管理员收藏文件
type FileFavor struct {
	AdminUserID uint32           `gorm:"index;comment:管理员ID;" json:"admin_user_id"`
	FileID      uint32           `gorm:"index;comment:文件ID;" json:"file_id"`
	CreatedAt   global.LocalTime `gorm:"type:datetime;comment:创建时间;" json:"created_at,omitempty" validate:"-"`
}

var (
	FileTypeImage = []string{
		"jpg", "jpeg", "jpe", "jpz", "png", "pnz", "gif", "bmp",
	}

	FileTypeDoc = []string{
		"doc", "docx", "xls", "xlsx", "pdf", "ppt", "pptx", "txt", "md",
	}

	FileTypeAudio = []string{
		"mp3", "wma", "wav",
	}

	FileTypeVideo = []string{
		"mp4", "avi", "mov", "wmv", "mpg", "mpeg",
	}

	FileTypeMap = map[string][]string{
		"image": FileTypeImage,
		"doc":   FileTypeDoc,
		"audio": FileTypeAudio,
		"video": FileTypeVideo,
	}
)

// FileType 文件类型转换
func FileType(ext string) string {
	ext = ext[strings.Index(ext, ".")+1:]
	for t, exts := range FileTypeMap {
		if util.InArray(strings.ToLower(ext), exts) {
			return t
		}
	}
	return ""
}

// FileCategory 文件分类模型
type FileCategory struct {
	ID        uint32           `gorm:"primaryKey;" json:"value"` // json:"id"
	CreatedAt global.LocalTime `gorm:"type:datetime;comment:创建时间;" json:"created_at,omitempty" validate:"-"`
	Name      string           `gorm:"type:varchar(255);not null;comment:名称;" json:"title" validate:"required" comment:"名称"` // json:"name"
	PID       uint32           `gorm:"default:0;comment:父分类ID;" json:"p_id" validate:"-" comment:"父分类ID"`
	Selected  bool             `gorm:"-" json:"selected" validate:"-" comment:"-"`
	HTML      string           `gorm:"-" json:"html" validate:"-"` // 用来输出层级 |----
	Level     int              `gorm:"-" json:"-" validate:"-"`    // 计算层级
}

// FileCategoryTree 文件分类树
type FileCategoryTree struct {
	*FileCategory
	Children []*FileCategoryTree `json:"children"`
}

// GetID 返回ID，实现接口`dao.Model`的方法
func (m *FileCategory) GetID() uint32 {
	return m.ID
}

// LoadRelatedField 加载相关字段
func (m *File) LoadRelatedField(adminUserID uint32) {
	if adminUserID > 0 {
		var favor = &FileFavor{}
		dao.FindOne(nil, favor, map[string]interface{}{"file_id": m.ID, "admin_user_id": adminUserID})
		if favor.FileID == 0 {
			m.IsFavor = -1
		} else {
			m.IsFavor = 1
		}
	} else {
		m.IsFavor = -1
	}

	if m.Path != "" {
		switch m.Storage {
		case "local":
			m.URL = config.Upload.Local.Domain + m.Path
		case "qiniu":
			m.URL = config.Upload.Qiniu.Domain + m.Path
		case "tencent":
			m.URL = config.Upload.Tencent.Domain + m.Path
		case "aliyun":
			m.URL = config.Upload.Aliyun.Domain + m.Path
		}
	}
}

// GetFileCategoryChildIDs 获取所有下级文件分类ID
func GetFileCategoryChildIDs(allIDs *[]uint32, id uint32) {
	if cids, has := hasFileCategoryChildIDs(id); has {
		*allIDs = append(*allIDs, cids...)
		for _, cid := range cids {
			GetFileCategoryChildIDs(allIDs, cid)
		}
	}
}

func hasFileCategoryChildIDs(id uint32) (ids []uint32, has bool) {
	dao.Pluck(nil, &FileCategory{}, map[string]interface{}{"p_id": id}, "id", &ids)
	if len(ids) > 0 {
		has = true
	}
	return
}

// GetFileCategoryTree 获取文件分类树
func GetFileCategoryTree() []*FileCategoryTree {
	var (
		trees    []*FileCategoryTree
		rootTree = &FileCategoryTree{
			FileCategory: &FileCategory{
				ID:   0,
				Name: "根节点",
			},
		}
		allCategories []*FileCategory
	)
	dao.FindAll(nil, &allCategories, nil)
	for _, v := range allCategories {
		trees = append(trees, &FileCategoryTree{FileCategory: v})
	}
	makeFileCategoryTree(trees, rootTree)
	return rootTree.Children
}

func makeFileCategoryTree(list []*FileCategoryTree, node *FileCategoryTree) {
	if children, has := hasFileCategoryChild(list, node); has {
		node.Children = append(node.Children, children...)
		for _, v := range children {
			makeFileCategoryTree(list, v)
		}
	} else {
		node.Children = make([]*FileCategoryTree, 0) // 没有子节点将 nil 转成空切片，输出json 是[] 而不是 null
	}
}

func hasFileCategoryChild(list []*FileCategoryTree, node *FileCategoryTree) (children []*FileCategoryTree, has bool) {
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

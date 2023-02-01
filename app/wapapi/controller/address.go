/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-03-29 15:12:04
 * @LastEditTime: 2021-05-20 16:00:07
 */
package controller

import (
	"encoding/json"
	"iris-project/app"
	"iris-project/app/dao"
	"iris-project/app/model"
	"iris-project/lib/cache"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"gorm.io/gorm"
)

// Address 控制器
type Address struct {
	Ctx       iris.Context
	WapUserID uint32
}

// BeforeActivation 前置方法
func (c *Address) BeforeActivation(b mvc.BeforeActivation) {
	b.Dependencies().Register(GetAuthUserID)
}

// GetCityTree 省市区
func (c *Address) GetCityTree() {
	var tree []*model.CityTree
	treeStr, _ := cache.Get("city_tree")
	if treeStr == "" {
		tree = model.GetCityTree()
		json, _ := json.Marshal(tree)
		cache.Set("city_tree", string(json), time.Minute*time.Duration(60))
	} else {
		json.Unmarshal([]byte(treeStr), &tree)
	}
	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", tree))
}

// GetAddressList 地址列表
func (c *Address) GetAddressList() {
	page, size := app.Pager(c.Ctx)

	var where = make(map[string]interface{})
	where["user_id"] = c.WapUserID

	var (
		list  []*model.Address
		total int64
	)

	searchList := &dao.SearchListData{
		Where:    where,
		Fields:   []string{},
		OrderBys: []string{"id desc"},
		Preloads: []string{},
		Page:     page,
		Size:     size,
	}
	if err := searchList.GetList(&list, &total); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	for _, v := range list {
		v.LoadRelatedField()
	}

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", app.List{List: list, Total: total}))
}

// GetAffairInfo 获取详情
func (c *Address) GetAddressInfo() {
	var id int
	if id = c.Ctx.URLParamIntDefault("id", 0); id == 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, "参数错误", nil))
		return
	}
	var data = &model.Address{}
	dao.FindOne(nil, data, map[string]interface{}{"id": id, "user_id": c.WapUserID})
	if data.ID == 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeNotFound, "", nil))
		return
	}
	// if data.Status != 1 {
	// 	c.Ctx.JSON(app.APIData(false, app.CodeDisabled, "", nil))
	// 	return
	// }
	data.LoadRelatedField()

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", data))
}

// GetAddressDefault 获取默认地址
func (c *Address) GetAddressDefault() {
	var data = &model.Address{}
	dao.FindOne(nil, data, map[string]interface{}{"is_default": 1, "user_id": c.WapUserID})
	if data.ID == 0 {
		data.Province = "请选择服务地址"
		// c.Ctx.JSON(app.APIData(false, app.CodeNotFound, "", nil))
		// return
	} else {
		data.LoadRelatedField()
	}
	// if data.Status != 1 {
	// 	c.Ctx.JSON(app.APIData(false, app.CodeDisabled, "", nil))
	// 	return
	// }

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", data))
}

// PostAddress 创建或更新
func (c *Address) PostAddress() {
	var data = &model.Address{}

	errmsg := app.CheckRequest(c.Ctx, data)
	if len(errmsg) != 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}

	data.UserID = c.WapUserID

	if data.IsDefault == 1 {
		if err := dao.UpdateAll(nil, &model.Address{}, map[string]interface{}{"user_id": c.WapUserID}, map[string]interface{}{"is_default": -1}); err != nil {
			c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
			return
		}
	}

	if err := dao.CreateUpdate(data); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// DeleteAddressBy 删除
func (c *Address) DeleteAddress() {
	var id int
	if id = c.Ctx.URLParamIntDefault("id", 0); id == 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, "参数错误", nil))
		return
	}
	var data dao.Model = &model.Address{ID: uint32(id)}

	if err := dao.DeleteByID(data); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// PutAddressDefaultBy 修改状态
func (c *Address) PutAddressDefault() {
	var id int
	if id = c.Ctx.URLParamIntDefault("id", 0); id == 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, "参数错误", nil))
		return
	}
	var data dao.Model = &model.Address{ID: uint32(id)}

	if !dao.GetByID(data) {
		c.Ctx.JSON(app.APIData(true, app.CodeUserNotFound, "", nil))
		return
	}

	dao.Transaction(func(tx *gorm.DB) error {
		if err := dao.UpdateAll(tx, &model.Address{}, map[string]interface{}{"user_id": c.WapUserID}, map[string]interface{}{"is_default": -1}); err != nil {
			return err
		}
		if err := dao.UpdateByID(data, map[string]interface{}{"is_default": 1}, tx); err != nil {
			return err
		}
		return nil
	})

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

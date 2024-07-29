/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-03-09 16:15:38
 * @LastEditTime: 2024-07-29 14:30:41
 */
package controller

import (
	"encoding/json"
	"iris-project/app"
	"iris-project/app/adminapi/model"
	"iris-project/app/config"
	"iris-project/app/dao"
	appmodel "iris-project/app/model"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// Config 配置
type Config struct {
	Ctx           iris.Context
	AuthAdminUser *model.AdminUser
}

// BeforeActivation 前置方法
func (c *Config) BeforeActivation(b mvc.BeforeActivation) {
	b.Dependencies().Register(GetAuthAdminUser)
}

// GetConfigImportTemplates 获取导入模板
func (c *Config) GetConfigImportTemplates() {
	var templates = map[string]string{
		"building_template":     config.App.Fronturl + "/asserts/template/building_template.xlsx",
		"house_template":        config.App.Fronturl + "/asserts/template/house_template.xlsx",
		"proprietor_template":   config.App.Fronturl + "/asserts/template/proprietor_template.xlsx",
		"parkingspace_template": config.App.Fronturl + "/asserts/template/parkingspace_template.xlsx",
		"product_template":      config.App.Fronturl + "/asserts/template/product_template.xlsx",
	}
	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", templates))
}

// @Tags        配置
// @Summary		配置列表
// @Description	配置列表
// @Accept		json
// @Produce		json
// @Param		name     query string false	"名称"
// @Param		type     query int    false	"状态"
// @Param		page     query string false	"页码"
// @Param		size     query int    false	"页大小"
// @Success		200		{object}	app.Response		""
// @Failure		200		{object}	app.Response	    ""
// @Router		/adminapi/config/list [get]
func (c *Config) GetConfigList() {
	page, size := app.Pager(c.Ctx)

	var where = make(map[string]interface{})
	if name := c.Ctx.URLParamDefault("name", ""); name != "" {
		where["name like"] = "%" + name + "%"
	}

	if configType := c.Ctx.URLParamDefault("type", ""); configType != "" {
		where["type"] = configType
	}

	var (
		list  []*appmodel.Config
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

// GetConfigInfo 详情
func (c *Config) GetConfigInfoBy(id uint32) {
	var config dao.Model = &appmodel.Config{ID: id}

	if !dao.GetByID(config) {
		c.Ctx.JSON(app.APIData(true, app.CodeNotFound, "", nil))
		return
	}
	config.(*appmodel.Config).LoadRelatedField()

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", config))
}

// GetConfigInfo 根据code获取
func (c *Config) GetConfigInfo() {
	var code string
	if code = c.Ctx.URLParamDefault("code", ""); code == "" {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, "参数错误", nil))
		return
	}
	var config = &appmodel.Config{}
	dao.FindOne(nil, config, map[string]interface{}{"code": code})
	if config.ID == 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, "没有该配置项", nil))
		return
	}
	config.LoadRelatedField()

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", config))
}

// DeleteCommunityBy 删除
func (c *Config) DeleteConfigBy(id uint32) {
	var config dao.Model = &appmodel.Config{ID: id}

	if err := dao.DeleteByID(config); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// PutCommunityStatusBy 修改状态
func (c *Config) PutConfigStatusBy(id uint32) {
	var config dao.Model = &appmodel.Config{ID: id}

	if !dao.GetByID(config) {
		c.Ctx.JSON(app.APIData(true, app.CodeNotFound, "", nil))
		return
	}
	var status int
	if config.(*appmodel.Config).Status == 1 {
		status = -1
	} else if config.(*appmodel.Config).Status == -1 {
		status = 1
	} else {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, "状态错误", nil))
		return
	}

	if err := dao.UpdateByID(config, map[string]interface{}{"status": status}); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// PostConfigBanner Banner上传
func (c *Config) PostConfigBanner() {
	type Request struct {
		Code    string                  `json:"code" validate:"required"`
		Name    string                  `json:"name" validate:"required"`
		Banners []appmodel.ConfigBanner `json:"banners" validate:"required"`
	}
	var param = &Request{}
	errmsg := app.CheckRequest(c.Ctx, param)
	if len(errmsg) != 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}

	content, _ := json.Marshal(param.Banners)
	config := &appmodel.Config{}

	dao.FindOne(nil, config, map[string]interface{}{"code": param.Code})
	if config.ID == 0 {
		config = &appmodel.Config{
			Code:    param.Code,
			Type:    "banner",
			Name:    param.Name,
			Content: string(content),
		}
	} else {
		config.Name = param.Name
		config.Content = string(content)
	}

	if err := dao.CreateUpdate(config); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// PostConfigField 添加键值对
func (c *Config) PostConfigField() {
	type Request struct {
		Code    string `json:"code" validate:"required"`
		Name    string `json:"name" validate:"required"`
		Content string `json:"content" validate:"required"`
	}
	var param = &Request{}
	errmsg := app.CheckRequest(c.Ctx, param)
	if len(errmsg) != 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}

	config := &appmodel.Config{}

	dao.FindOne(nil, config, map[string]interface{}{"code": param.Code})
	if config.ID == 0 {
		config = &appmodel.Config{
			Code:    param.Code,
			Type:    "field",
			Name:    param.Name,
			Content: param.Content,
		}
	} else {
		config.Name = param.Name
		config.Content = param.Content
	}

	if err := dao.CreateUpdate(config); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// PostConfigArray 数组
func (c *Config) PostConfigArray() {
	type Request struct {
		Code    string        `json:"code" validate:"required"`
		Name    string        `json:"name" validate:"required"`
		Content []interface{} `json:"content" validate:"required"`
	}
	var param = &Request{}
	errmsg := app.CheckRequest(c.Ctx, param)
	if len(errmsg) != 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}

	for i, v := range param.Content {
		var str, ok = v.(string)
		if !ok {
			var obj interface{}
			if err := json.Unmarshal([]byte(str), &obj); err == nil {
				param.Content[i] = obj
			} else {
				// log.Println(err.Error())
			}
		}
	}

	content, _ := json.Marshal(param.Content)
	config := &appmodel.Config{}

	dao.FindOne(nil, config, map[string]interface{}{"code": param.Code})
	if config.ID == 0 {
		config = &appmodel.Config{
			Code:    param.Code,
			Type:    "array",
			Name:    param.Name,
			Content: string(content),
		}
	} else {
		config.Name = param.Name
		config.Content = string(content)
	}

	if err := dao.CreateUpdate(config); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// PostConfigObject 对象
func (c *Config) PostConfigObject() {
	type Request struct {
		Code    string      `json:"code" validate:"required"`
		Name    string      `json:"name" validate:"required"`
		Content interface{} `json:"content" validate:"required"`
	}
	var param = &Request{}
	errmsg := app.CheckRequest(c.Ctx, param)
	if len(errmsg) != 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}

	content, _ := json.Marshal(param.Content)
	config := &appmodel.Config{}

	dao.FindOne(nil, config, map[string]interface{}{"code": param.Code})
	if config.ID == 0 {
		config = &appmodel.Config{
			Code:    param.Code,
			Type:    "object",
			Name:    param.Name,
			Content: string(content),
		}
	} else {
		config.Name = param.Name
		config.Content = string(content)
	}

	if err := dao.CreateUpdate(config); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// PutConfigIsshow 设置客户端显示或隐藏某些功能
// func (c *Config) PutConfigIsshow() {
// 	type Request struct {
// 		Key string `json:"key" validate:"required"`
// 	}
// 	var param = &Request{}
// 	errmsg := app.CheckRequest(c.Ctx, param)
// 	if len(errmsg) != 0 {
// 		c.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
// 		return
// 	}
// 	key := param.Key
// 	var show, _ = cache.Get(key)
// 	if show == "" {
// 		cache.Set(key, "1", time.Duration(0))
// 	} else {
// 		if show == "1" {
// 			cache.Set(key, "0", time.Duration(0))
// 		} else {
// 			cache.Set(key, "1", time.Duration(0))
// 		}
// 	}

// 	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
// }

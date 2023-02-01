/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-03-19 14:01:28
 * @LastEditTime: 2022-09-29 09:52:38
 */
package controller

import (
	"iris-project/app"
	"iris-project/app/dao"
	"iris-project/app/model"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// Home 控制器
type Home struct {
	Ctx iris.Context
}

// BeforeActivation 前置方法
func (c *Home) BeforeActivation(b mvc.BeforeActivation) {
}

// GetHome 首页信息
func (c *Home) GetHome() {

}

// GetHomeSystemConfig 获取系统配置
func (c *Home) GetHomeSystemConfig() {
	var code string
	if code = c.Ctx.URLParamDefault("code", ""); code == "" {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, "参数错误", nil))
		return
	}
	var config = &model.Config{}
	dao.FindOne(nil, config, map[string]interface{}{"code": code})
	if config.ID == 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, "没有该配置项", nil))
		return
	}
	config.LoadRelatedField()

	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", config))
}

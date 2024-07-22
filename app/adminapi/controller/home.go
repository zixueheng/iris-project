/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:28:08
 * @LastEditTime: 2024-01-04 14:23:59
 */
package controller

import (
	"iris-project/app"
	"iris-project/app/adminapi/model"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// Home 控制器
type Home struct {
	Ctx           iris.Context
	AuthAdminUser *model.AdminUser
}

// BeforeActivation 前置方法
func (c *Home) BeforeActivation(b mvc.BeforeActivation) {
	b.Dependencies().Register(GetAuthAdminUser)
}

// GetHomeStatistic 首页统计
func (c *Home) GetHomeStatistic() {
	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// GetHomeJnotice ...
func (c *Home) GetHomeJnotice() {
	// icon: "md-bulb", iconColor: "#87d068", title: "您还有1个待发货的订单", url: "", type: "bulb", read: 0, time: 0
	type Msg struct {
		Icon      string `json:"icon"`
		IconColor string `json:"iconColor"`
		Title     string `json:"title"`
		URL       string `json:"url"`
		Type      string `json:"type"`
		Read      int    `json:"read"`
		Time      int    `json:"time"`
	}
	var res = []Msg{
		{Icon: "md-bulb", IconColor: "#87d068", Title: "您还有1个待发货的订单", URL: "", Type: "bulb", Read: 0, Time: 0},
	}
	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", res))
}

// GetHomeLogout ...
func (c *Home) GetHomeLogout() {
	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

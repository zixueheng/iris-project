/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-04 17:11:48
 * @LastEditTime: 2021-05-17 09:56:25
 */
package controller

import (
	"iris-project/app"
	"iris-project/app/adminapi/model"
	"iris-project/app/dao"
	appmodel "iris-project/app/model"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// User 会员控制器
type User struct {
	Ctx           iris.Context
	AuthAdminUser *model.AdminUser
}

// BeforeActivation 前置方法
func (c *User) BeforeActivation(b mvc.BeforeActivation) {
	b.Dependencies().Register(GetAuthAdminUser)
}

// GetUserList 列表
func (c *User) GetUserList() {
	page, size := app.Pager(c.Ctx)

	var where = make(map[string]interface{})
	if realname := c.Ctx.URLParamDefault("realname", ""); realname != "" {
		where["realname like"] = "%" + realname + "%"
	}
	if nickname := c.Ctx.URLParamDefault("nickname", ""); nickname != "" {
		where["nickname like"] = "%" + nickname + "%"
	}
	if phone := c.Ctx.URLParamDefault("phone", ""); phone != "" {
		where["phone"] = phone
	}
	if status := c.Ctx.URLParamIntDefault("status", 0); status != 0 {
		where["status"] = status
	}

	var (
		list  []*appmodel.User
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

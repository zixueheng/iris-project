/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-18 09:34:00
 * @LastEditTime: 2021-05-31 17:53:13
 */
package controller

import (
	"iris-project/app"
	"iris-project/app/dao"
	"iris-project/app/model"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// User 控制器
type User struct {
	Ctx       iris.Context
	WapUserID uint32
}

// BeforeActivation 前置方法
func (c *User) BeforeActivation(b mvc.BeforeActivation) {
	b.Dependencies().Register(GetAuthUserID)
}

// GetUserInfo 用户信息
func (c *User) GetUserInfo() {
	user := &model.User{ID: c.WapUserID}
	if dao.GetByID(user) {
		user.LoadRelatedField()
		c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", user))
		return
	}

	c.Ctx.JSON(app.APIData(false, app.CodeCustom, "用户不存在", nil))
}

// PutUserInfo 修改用户信息
func (c *User) PutUserInfo() {
	type Request struct {
		Realname string ` json:"realname" validate:"required" comment:"姓名"`
		Nickname string `json:"nickname" validate:"required" comment:"昵称"`
		Avatar   string `json:"avatar" validate:"required" comment:"头像"`
		Birthday string `json:"birthday" validate:"required" comment:"生日"`
		Gender   int8   `json:"gender" validate:"numeric,oneof=0 1 2" comment:"生日"`
	}
	var param = &Request{}
	errmsg := app.CheckRequest(c.Ctx, param)
	if len(errmsg) != 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}
	data := map[string]interface{}{
		"realname": param.Realname,
		"nickname": param.Nickname,
		"avatar":   param.Avatar,
		"birthday": param.Birthday,
		"gender":   param.Gender,
	}
	if err := dao.UpdateByID(&model.User{ID: c.WapUserID}, data); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}
	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

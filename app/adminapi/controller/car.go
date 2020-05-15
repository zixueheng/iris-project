package controller

import (
	"iris-project/lib/util"

	"github.com/kataras/iris/v12"
)

// Car 控制器
type Car struct {
	Ctx iris.Context // IRIS框架会自动注入 Context
}

// GetCarBy 根据 carname 查找用户
func (c *Car) GetCarBy(carname string) string {
	adminUserID, _ := c.Ctx.Values().GetInt("auth_admin_user_id")
	return "adminUserID: " + util.ParseString(adminUserID)
}

package controller

import (
	"iris-project/lib/util"

	"github.com/kataras/iris/v12"
)

// Goods 控制器
type Goods struct {
	Ctx iris.Context // IRIS框架会自动注入 Context
}

// GetGoodslistBy 获取商品列表
func (g *Goods) GetGoodslistBy(page, size uint) string {
	adminUserID, _ := g.Ctx.Values().GetInt("auth_admin_user_id")
	return "adminUserID: " + util.ParseString(adminUserID)
}

// GetGoodsBy 获取商品详情
func (g *Goods) GetGoodsBy(id uint) string {
	adminUserID, _ := g.Ctx.Values().GetInt("auth_admin_user_id")
	return "adminUserID: " + util.ParseString(adminUserID)
}

// PutGoods ...
func (g *Goods) PutGoods() string {
	adminUserID, _ := g.Ctx.Values().GetInt("auth_admin_user_id")
	return "adminUserID: " + util.ParseString(adminUserID)
}

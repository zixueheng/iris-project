/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-18 09:34:00
 * @LastEditTime: 2022-09-29 09:44:38
 */
package controller

import (
	"errors"
	"iris-project/app"
	adminmodel "iris-project/app/adminapi/model"
	"iris-project/app/adminapi/validate"
	"iris-project/app/dao"
	"log"

	"iris-project/app/config"
	"iris-project/global"
	"iris-project/lib/cache"
	"iris-project/lib/util"
	"iris-project/middleware"
	"time"

	"github.com/go-redis/redis"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/jameskeane/bcrypt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/shopspring/decimal"
)

// Public 公共控制器，不需要验证token
type Public struct {
	Ctx iris.Context
}

// AfterActivation 后置方法
func (p *Public) AfterActivation(a mvc.AfterActivation) {
	// 给单独的控制器方法添加中间件
	// select the route based on the method name you want to modify.
	refreshtokenRoute := a.GetRoute("PostRefreshtoken") // 根据 方法名 获取 方法的路由
	// just prepend the handler(s) as middleware(s) you want to use. or append for "done" handlers.
	refreshtokenRoute.Handlers = append([]iris.Handler{middleware.JwtHandler().Serve}, refreshtokenRoute.Handlers...) // 将中间件 追加到 路由的 Handleers 字段中

	resetPasswordRoute := a.GetRoute("PostResetPassword")
	resetPasswordRoute.Handlers = append([]iris.Handler{middleware.JwtHandler().Serve}, resetPasswordRoute.Handlers...)
}

// PostLogin 登录
func (p *Public) PostLogin() {
	loginInfo := new(validate.LoginRequest)

	errmsg := app.CheckRequest(p.Ctx, loginInfo)
	if len(errmsg) != 0 {
		p.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}

	adminUser := &adminmodel.AdminUser{}
	response, ok, code := adminUser.CheckLogin(loginInfo)

	p.Ctx.JSON(app.APIData(ok, code, "", response))
}

func (p *Public) GetTestIp() {
	// if err := checkIPCount(p.Ctx.RemoteAddr()); err != nil {
	// 	p.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
	// 	return
	// }

	if err := checkIPCount(util.GetOutboundIP().String()); err != nil {
		p.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	p.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

func checkIPCount(remoteIP string) error {
	log.Println("IP: ", remoteIP)
	if remoteIP == "" {
		return errors.New("客户端IP无效")
	}
	var cacheKey = "verify-code-count-" + remoteIP
	countStr, err := cache.Get(cacheKey)
	log.Println("countStr", countStr)
	if countStr == "" || err == redis.Nil {
		log.Println("初始设置")
		cache.Set(cacheKey, int(1), time.Second*time.Duration(60))
		return nil
	}

	count := util.ParseInt(countStr)
	if count > 3 {
		return errors.New("操作过于频繁，请稍后再试")
	}

	ttl := cache.GetCacheInstance().TTL(cacheKey).Val().Seconds()
	left := decimal.NewFromFloat(ttl).IntPart()
	log.Println("剩余：", left)
	cache.Set(cacheKey, count+1, time.Second*time.Duration(left))

	return nil
}

// PostRefreshtoken 刷新缓存（验证token不判断是否过期）
func (p *Public) PostRefreshtoken() {
	value := p.Ctx.Values().Get("jwt").(*jwt.Token) // 这里要先给这个路由方法添加JWT的中间件才能获取到 jwt变量

	data := value.Claims.(jwt.MapClaims)

	var adminUseID string
	if value, ok := data[global.AdminUserJWTKey]; ok {
		adminUseID = value.(string)
	} else {
		app.ResponseProblemHTTPCode(p.Ctx, iris.StatusBadRequest, errors.New("Token中没有"+global.AdminUserJWTKey))
	}

	param := struct {
		RefreshToken string `json:"refresh_token"`
	}{}

	if err := p.Ctx.ReadJSON(&param); err != nil {
		p.Ctx.JSON(app.APIData(false, app.CodeRequestParamError, "", nil))
		return
	}
	if param.RefreshToken == "" {
		p.Ctx.JSON(app.APIData(false, app.CodeRequestParamError, "", nil))
		return
	}
	adminUser := &adminmodel.AdminUser{ID: uint32(util.ParseInt(adminUseID))}
	if !dao.GetByID(adminUser) {
		p.Ctx.JSON(app.APIData(false, app.CodeUserNotFound, "", nil))
		return
	}

	if param.RefreshToken != adminUser.RefreshToken {
		p.Ctx.JSON(app.APIData(false, app.CodeRefreshTokenInvalidated, "", nil))
		return
	}

	if time.Now().Unix() > int64(adminUser.RefreshTokenExpired) { // Token 超时
		p.Ctx.JSON(app.APIData(false, app.CodeRefreshTokenExpired, "", nil))
		return
	}

	token, refreshToken, tokenExpired, refreshTokenExpired := app.GenTokenAndRefreshToken(global.AdminUserJWTKey, int(adminUser.ID), global.AdminTokenMinutes, global.AdminRefreshTokenMinutes)
	dao.UpdateByID(adminUser, map[string]interface{}{"refresh_token": refreshToken, "refresh_token_expired": refreshTokenExpired.Unix()}) // 保存刷新token和过期时间至数据库
	response := struct {
		Token               string `json:"token"`
		TokenExpired        int64  `json:"token_expired"`
		RefreshToken        string `json:"refresh_token"`
		RefreshTokenExpired int64  `json:"refresh_token_expired"`
	}{
		Token:               token,
		TokenExpired:        tokenExpired.Unix(),
		RefreshToken:        refreshToken,
		RefreshTokenExpired: refreshTokenExpired.Unix(),
	}
	p.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", response))
}

// PostResetPassword 修改密码（验证token并判断是否过期）
func (p *Public) PostResetPassword() {
	value := p.Ctx.Values().Get("jwt").(*jwt.Token) // 这里要先给这个路由方法添加JWT的中间件才能获取到 jwt变量

	data := value.Claims.(jwt.MapClaims)

	var adminUseID, exp string
	if value, ok := data[global.AdminUserJWTKey]; ok {
		adminUseID = value.(string)
	} else {
		app.ResponseProblemHTTPCode(p.Ctx, iris.StatusBadRequest, errors.New("Token中没有"+global.AdminUserJWTKey))
		return
	}

	if value, ok := data["exp"]; ok {
		exp = value.(string)
	} else {
		app.ResponseProblemHTTPCode(p.Ctx, iris.StatusBadRequest, errors.New("Token中没有exp"))
		return
	}

	expObj, err := time.ParseInLocation(config.App.Timeformat, exp, time.Local)
	if err != nil { // 过期时间解析错误，返回 BadRequest
		app.ResponseProblemHTTPCode(p.Ctx, iris.StatusBadRequest, err)
	}

	if expObj.Before(time.Now()) { // Token 超时
		p.Ctx.JSON(app.APIData(false, app.CodeTokenExpired, "", nil))
		return
	}

	param := struct {
		OldPassword string `json:"old_password" validate:"required" comment:"原密码"`
		NewPassword string `json:"new_password" validate:"required" comment:"新密码"`
	}{}

	if err := p.Ctx.ReadJSON(&param); err != nil {
		p.Ctx.JSON(app.APIData(false, app.CodeRequestParamError, "", nil))
		return
	}
	if param.OldPassword == "" {
		p.Ctx.JSON(app.APIData(false, app.CodeRequestParamError, "", nil))
		return
	}
	if param.NewPassword == "" {
		p.Ctx.JSON(app.APIData(false, app.CodeRequestParamError, "", nil))
		return
	}
	password := struct{ Password string }{}
	global.Db.Model(&adminmodel.AdminUser{}).Where("id=?", util.ParseInt(adminUseID)).Select("password").Scan(&password)
	if !bcrypt.Match(param.OldPassword, password.Password) {
		p.Ctx.JSON(app.APIData(false, app.CodeUserPasswordError, "", nil))
		return
	}

	if global.Db.Model(&adminmodel.AdminUser{}).Where("id=?", util.ParseInt(adminUseID)).Update("password", util.HashPassword(param.NewPassword)).RowsAffected > 0 {
		p.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
		return
	}
	p.Ctx.JSON(app.APIData(false, app.CodeFailed, "", nil))
}

// GetCheckAuth ...
func (p *Public) GetCheckAuth() {
	p.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

// GetAuth ...
func (p *Public) GetAuth() {
	p.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", nil))
}

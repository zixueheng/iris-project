/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-18 09:34:00
 * @LastEditTime: 2021-11-02 09:37:07
 */
package controller

import (
	"errors"
	"iris-project/app"
	"iris-project/global"
	"iris-project/lib/util"
	"iris-project/middleware"

	"iris-project/app/dao"
	"iris-project/app/model"
	"iris-project/app/service"
	"iris-project/app/wapapi/validate"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// Public 控制器
type Public struct {
	Ctx iris.Context
}

// AfterActivation 后置方法
func (c *Public) AfterActivation(a mvc.AfterActivation) {
	refreshtokenRoute := a.GetRoute("PostRefreshtoken")
	refreshtokenRoute.Handlers = append([]iris.Handler{middleware.JwtHandler().Serve}, refreshtokenRoute.Handlers...)
}

func (c *Public) GetVerifycodes() {
	var where = make(map[string]interface{})
	if phone := c.Ctx.URLParamDefault("phone", ""); phone != "" {
		where["phone"] = phone
	}
	var codes []model.Verifycode
	dao.FindAll(nil, &codes, where)
	c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", codes))
}

// PostPhoneLoginregisterVerifycode 发送验证码
func (c *Public) PostPhoneLoginregisterVerifycode() {
	type Request struct {
		Phone string `json:"phone" validate:"required,len=11" comment:"手机号"`
	}
	var param = &Request{}
	errmsg := app.CheckRequest(c.Ctx, param)
	if len(errmsg) != 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}

	if _, err := service.VerifycodeGen(param.Phone, "phone_login_register"); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}
	c.Ctx.JSON(app.APIData(true, app.CodeVerifycodeSucceed, "", nil))
}

// PostUserphoneLoginregister 登录注册，没有账号生成账号并登录
func (c *Public) PostUserphoneLoginregister() {
	loginInfo := new(validate.UserphoneLoginRegister)

	errmsg := app.CheckRequest(c.Ctx, loginInfo)
	if len(errmsg) != 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}

	if err := service.VerifycodeCheck(loginInfo.Phone, loginInfo.Verifycode, "phone_login_register"); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	response, ok, code := service.UserPhoneLoginRegister(loginInfo)

	c.Ctx.JSON(app.APIData(ok, code, "", response))
}

// PostVerifycode 发送验证码
func (c *Public) PostPhoneregisterVerifycode() {
	type Request struct {
		Phone string `json:"phone" validate:"required,len=11" comment:"手机号"`
	}
	var param = &Request{}
	errmsg := app.CheckRequest(c.Ctx, param)
	if len(errmsg) != 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}

	if _, err := service.VerifycodeGen(param.Phone, "phone_register"); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}
	c.Ctx.JSON(app.APIData(true, app.CodeVerifycodeSucceed, "", nil))
}

// PostUserphoneLogin 手机号密码登录
func (c *Public) PostUserphoneLogin() {
	loginInfo := new(validate.UserphoneLogin)

	errmsg := app.CheckRequest(c.Ctx, loginInfo)
	if len(errmsg) != 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}

	response, ok, code := service.UserPhoneLogin(loginInfo)

	c.Ctx.JSON(app.APIData(ok, code, "", response))
}

// PostUsernameLogin 用户名密码登录
func (c *Public) PostUsernameLogin() {
	loginInfo := new(validate.UsernameLogin)

	errmsg := app.CheckRequest(c.Ctx, loginInfo)
	if len(errmsg) != 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}

	response, ok, code := service.UserUsernameLogin(loginInfo)

	c.Ctx.JSON(app.APIData(ok, code, "", response))
}

// PostUserphoneRegister 手机号密码注册
func (c *Public) PostUserphoneRegister() {
	registerInfo := new(validate.UserphoneRegister)

	errmsg := app.CheckRequest(c.Ctx, registerInfo)
	if len(errmsg) != 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}
	if err := service.VerifycodeCheck(registerInfo.Phone, registerInfo.Verifycode, "phone_register"); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	}

	if res, err := service.UserPhoneRegister(registerInfo); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	} else {
		c.Ctx.JSON(app.APIData(true, app.CodeUserRegisterSucceed, "", res))
	}
}

// PostUsernameRegister 用户名密码注册
func (c *Public) PostUsernameRegister() {
	registerInfo := new(validate.UsernameRegister)

	errmsg := app.CheckRequest(c.Ctx, registerInfo)
	if len(errmsg) != 0 {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, errmsg, nil))
		return
	}

	if res, err := service.UserUsernameRegister(registerInfo); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	} else {
		c.Ctx.JSON(app.APIData(true, app.CodeUserRegisterSucceed, "", res))
	}
}

// PostRefreshtoken 刷新缓存（验证token不判断是否过期）
func (c *Public) PostRefreshtoken() {
	value := c.Ctx.Values().Get("jwt").(*jwt.Token) // 这里要先给这个路由方法添加JWT的中间件才能获取到 jwt变量

	data := value.Claims.(jwt.MapClaims)

	var userID string
	if value, ok := data[global.WapUserJWTKey]; ok {
		userID = value.(string)
	} else {
		app.ResponseProblemHTTPCode(c.Ctx, iris.StatusBadRequest, errors.New("Token中没有"+global.WapUserJWTKey))
	}

	param := struct {
		RefreshToken string `json:"refresh_token"`
	}{}

	if err := c.Ctx.ReadJSON(&param); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeRequestParamError, "", nil))
		return
	}
	if param.RefreshToken == "" {
		c.Ctx.JSON(app.APIData(false, app.CodeRequestParamError, "", nil))
		return
	}

	if res, err := service.UserRefreshToken(uint32(util.ParseInt(userID)), param.RefreshToken); err != nil {
		c.Ctx.JSON(app.APIData(false, app.CodeCustom, err.Error(), nil))
		return
	} else {
		c.Ctx.JSON(app.APIData(true, app.CodeSucceed, "", res))
	}
}

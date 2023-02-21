/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-18 09:34:00
 * @LastEditTime: 2021-05-31 17:51:57
 */
package service

import (
	"errors"
	"iris-project/app"
	"iris-project/app/dao"
	"iris-project/app/model"
	"iris-project/app/wapapi/validate"
	"iris-project/global"
	"iris-project/lib/util"
	"time"

	"gorm.io/gorm"
)

// UserUsernameLogin 用户名密码登录
func UserUsernameLogin(loginInfo *validate.UsernameLogin) (interface{}, bool, app.Code) {
	username := &model.UserUsername{}
	if err := dao.FindOne(nil, username, map[string]interface{}{"username": loginInfo.Username}); err != nil {
		return nil, false, app.CodeUserNotFound
	}
	if username.ID == 0 {
		return nil, false, app.CodeUserNotFound
	}
	if util.MD5(loginInfo.Password) != username.Password {
		return nil, false, app.CodeUserPasswordError
	}
	user := &model.User{ID: username.UserID}
	if !dao.GetByID(user) {
		return nil, false, app.CodeUserNotFound
	} else if user.Status != 1 {
		return nil, false, app.CodeUserForbidden
	} else {
		response := GetUserTokenInfo(user)
		return response, true, app.CodeUserLoginSucceed
	}
}

// UserPhoneLogin 手机号密码登录
func UserPhoneLogin(loginInfo *validate.UserphoneLogin) (interface{}, bool, app.Code) {
	userphone := &model.UserPhone{}
	if err := dao.FindOne(nil, userphone, map[string]interface{}{"phone": loginInfo.Phone}); err != nil {
		return nil, false, app.CodeUserNotFound
	}
	if userphone.ID == 0 {
		return nil, false, app.CodeUserNotFound
	}
	if util.MD5(loginInfo.Password) != userphone.Password {
		return nil, false, app.CodeUserPasswordError
	}
	user := &model.User{ID: userphone.UserID}
	if !dao.GetByID(user) {
		return nil, false, app.CodeUserNotFound
	} else if user.Status != 1 {
		return nil, false, app.CodeUserForbidden
	} else {
		response := GetUserTokenInfo(user)
		return response, true, app.CodeUserLoginSucceed
	}
}

// UserPhoneLogin 手机号登录注册，没有账号生成账号并登录
func UserPhoneLoginRegister(loginInfo *validate.UserphoneLoginRegister) (interface{}, bool, app.Code) {
	userphone := &model.UserPhone{}
	if err := dao.FindOne(nil, userphone, map[string]interface{}{"phone": loginInfo.Phone}); err != nil {
		return nil, false, app.CodeUserNotFound
	}
	if userphone.ID == 0 { // 注册账号
		result, err := UserPhoneRegister(&validate.UserphoneRegister{Phone: loginInfo.Phone, Password: "123456"})
		if err != nil {
			return nil, false, app.CodeUserRegisterFailed
		}
		return result, true, app.CodeUserRegisterSucceed
	}

	user := &model.User{ID: userphone.UserID}
	if !dao.GetByID(user) {
		return nil, false, app.CodeUserNotFound
	} else if user.Status != 1 {
		return nil, false, app.CodeUserForbidden
	} else {
		response := GetUserTokenInfo(user)
		return response, true, app.CodeUserLoginSucceed
	}
}

// GetUserTokenInfo 生成token
func GetUserTokenInfo(user *model.User) interface{} {
	token, refreshToken, tokenExpired, refreshTokenExpired := app.GenTokenAndRefreshToken(global.WapUserJWTKey, int(user.ID), global.UserTokenMinutes,
		global.UserRefreshTokenMinutes,
		global.GetClient(global.WapAPI))

	dao.DeleteAll(nil, &model.UserRefreshToken{}, map[string]interface{}{"refresh_token_expired <": time.Now().Unix()})                             // 删除过期的
	dao.SaveOne(nil, &model.UserRefreshToken{UserID: user.ID, RefreshToken: refreshToken, RefreshTokenExpired: uint32(refreshTokenExpired.Unix())}) // 保存刷新token和过期时间
	user.LoadRelatedField()
	res := struct {
		Token               string      `json:"token"`
		TokenExpired        int64       `json:"token_expired"`
		RefreshToken        string      `json:"refresh_token"`
		RefreshTokenExpired int64       `json:"refresh_token_expired"`
		UserInfo            *model.User `json:"user_info"`
	}{
		Token:               token,
		TokenExpired:        tokenExpired.Unix(),
		RefreshToken:        refreshToken,
		RefreshTokenExpired: refreshTokenExpired.Unix(),
		UserInfo:            user,
	}
	return res
}

// UserUsernameRegister 用户名密码注册
func UserUsernameRegister(registerInfo *validate.UsernameRegister) (interface{}, error) {
	var username = &model.UserUsername{}
	count, _ := dao.Count(nil, username, map[string]interface{}{"username": registerInfo.Username})
	if count > 0 {
		return nil, errors.New("用户名已存在")
	}
	// birthday, err := util.TimeParse("1986-07-05 00:00:00", "")
	// if err != nil {
	// 	return err
	// }
	var user dao.Model = &model.User{Nickname: util.GetRandomString(6), Avatar: "", Realname: "", Phone: "" /*, Birthday: global.LocalDate{Time: birthday, Valid: true}*/}
	if err := dao.Transaction(func(tx *gorm.DB) error {
		if err := dao.CreateUpdate(user, tx); err != nil {
			return err
		}
		var username dao.Model = &model.UserUsername{UserID: user.GetID(), Username: registerInfo.Username, Password: util.MD5(registerInfo.Password)}
		if err := dao.CreateUpdate(username, tx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	response := GetUserTokenInfo(user.(*model.User))
	return response, nil
}

// UserPhoneRegister 手机号密码注册
func UserPhoneRegister(registerInfo *validate.UserphoneRegister) (interface{}, error) {
	var userphone = &model.UserPhone{}
	count, _ := dao.Count(nil, userphone, map[string]interface{}{"phone": registerInfo.Phone})
	if count > 0 {
		return nil, errors.New("手机号已存在")
	}
	// birthday, err := util.TimeParse("1986-07-05 00:00:00", "")
	// if err != nil {
	// 	return err
	// }
	var user dao.Model = &model.User{Nickname: util.GetRandomString(6), Avatar: "", Realname: "", Phone: registerInfo.Phone /*, Birthday: global.LocalDate{Time: birthday, Valid: true}*/}
	if err := dao.Transaction(func(tx *gorm.DB) error {
		if err := dao.CreateUpdate(user, tx); err != nil {
			return err
		}
		var userphone dao.Model = &model.UserPhone{UserID: user.GetID(), Phone: registerInfo.Phone, Password: util.MD5(registerInfo.Password)}
		if err := dao.CreateUpdate(userphone, tx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	response := GetUserTokenInfo(user.(*model.User))
	return response, nil
}

// UserRefreshToken 刷新TOKEN
func UserRefreshToken(userID uint32, refreshToken string) (interface{}, error) {
	user := &model.User{ID: userID}
	if !dao.GetByID(user) {
		return nil, errors.New("用户不存在")
	}

	m := &model.UserRefreshToken{}
	if err := dao.FindOne(nil, m, map[string]interface{}{"user_id": userID, "refresh_token": refreshToken}); err != nil {
		return nil, err
	}
	if m.UserID == 0 {
		return nil, errors.New("刷新Token不存在或已过期")
	}

	if time.Now().Unix() > int64(m.RefreshTokenExpired) { // Token 超时
		return nil, errors.New("刷新Token已过期")
	}

	response := GetUserTokenInfo(user)
	dao.DeleteAll(nil, &model.UserRefreshToken{}, map[string]interface{}{"user_id": userID, "refresh_token": refreshToken})
	return response, nil
}

// UserInfo 用户信息
func UserInfo(userID uint32) (*model.User, error) {
	user := &model.User{ID: userID}
	if !dao.GetByID(user) {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

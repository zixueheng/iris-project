/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-03-31 09:59:57
 * @LastEditTime: 2023-02-03 10:56:20
 */
package service

import (
	"encoding/json"
	"iris-project/app"
	"iris-project/app/dao"
	"iris-project/app/model"
	"iris-project/global"
	"iris-project/lib/util"
)

// MerchantLogin 商家登录
func MerchantLogin(username, password string) (interface{}, bool, app.Code) {
	merchant := &model.Merchant{}
	if err := dao.FindOne(nil, merchant, map[string]interface{}{"username": username}); err != nil {
		return nil, false, app.CodeUserNotFound
	}
	if merchant.ID == 0 {
		return nil, false, app.CodeUserNotFound
	}
	if util.MD5(password) != merchant.Password {
		return nil, false, app.CodeUserPasswordError
	}
	if merchant.Status != 1 {
		return nil, false, app.CodeUserForbidden
	} else {
		response := GetMerchantTokenInfo(merchant)
		return response, true, app.CodeUserLoginSucceed
	}
}

// GetMerchantTokenInfo 生成token
func GetMerchantTokenInfo(merchant *model.Merchant) interface{} {
	token, refreshToken, tokenExpired, refreshTokenExpired := app.GenTokenAndRefreshToken(global.MerchantJWTKey, int(merchant.ID), global.MerchantTokenMinutes,
		global.MerchantRefreshTokenMinutes,
		global.GetClient(global.MerchantAPI))
	// 不处理刷新token了
	res := struct {
		Token               string          `json:"token"`
		TokenExpired        int64           `json:"token_expired"`
		RefreshToken        string          `json:"refresh_token"`
		RefreshTokenExpired int64           `json:"refresh_token_expired"`
		MerchantInfo        *model.Merchant `json:"merchant_info"`
	}{
		Token:               token,
		TokenExpired:        tokenExpired.Unix(),
		RefreshToken:        refreshToken,
		RefreshTokenExpired: refreshTokenExpired.Unix(),
		MerchantInfo:        merchant,
	}
	return res
}

// GetMerchantWechatPayInfo 获取商户微信支付账号信息
func GetMerchantWechatPayInfo(merchantID uint32) *model.WechatPayInfo {
	if merchantID == 0 {
		return nil
	}
	var (
		wechatPayInfo   *model.WechatPayInfo
		merchantPayInfo = &model.MerchantPayInfo{}
	)

	dao.FindOne(nil, merchantPayInfo, map[string]interface{}{"channel": "wechat", "merchant_id": merchantID})
	if merchantPayInfo.MerchantID != 0 {
		// merchantPayInfo.LoadRelatedField()
		// wechatPayInfo = merchantPayInfo.PayInfoJSON.(*model.WechatPayInfo)
		wechatPayInfo = &model.WechatPayInfo{}
		json.Unmarshal([]byte(merchantPayInfo.PayInfo), wechatPayInfo)
	}
	return wechatPayInfo
}

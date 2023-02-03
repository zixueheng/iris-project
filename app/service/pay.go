/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2022-12-14 10:41:34
 * @LastEditTime: 2023-02-03 10:52:09
 */
package service

import (
	"context"
	"errors"
	"iris-project/app/config"
	"iris-project/app/model"
	fileutil "iris-project/lib/file"
	"iris-project/lib/util"
	"net/http"
	"strings"
	"time"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/xlog"
	"github.com/go-pay/gopay/wechat/v3"
	"github.com/shopspring/decimal"
)

// doPay 支付；
// attach: bill账单 / product-merchantID 产品
func doPay(payPrice float64, payType, orderSN, body, attach, openid string, wechatPayInfo *model.WechatPayInfo) (res interface{}, err error) {
	switch payType {
	case "alipay":
		return doAlipay(payPrice, orderSN, body, attach)
	case "weixinapppay":
		return doWeixinPay("app", payPrice, orderSN, body, attach, openid, wechatPayInfo)
	case "weixinminpay":
		return doWeixinPay("min", payPrice, orderSN, body, attach, openid, wechatPayInfo)
	case "unionpay":
	default:
		err = errors.New("支付方式错误")
		return
	}
	return
}

// doAlipay 支付宝支付
func doAlipay(payPrice float64, orderSN, body, attach string) (res interface{}, err error) {
	return nil, errors.New("未实现支付宝支付")
}

// 获取微信支付对象
func getWechatClient(attach string, payInfo *model.WechatPayInfo) (client *wechat.ClientV3, err error) {

	// NewClientV3 初始化微信客户端 v3
	// mchid：商户ID 或者服务商模式的 sp_mchid
	// serialNo：商户证书的证书序列号
	// apiV3Key：apiV3Key，商户平台获取
	// privateKey：私钥 apiclient_key.pem 读取后的内容

	if attach == "bill" { // 账单支付用单独的账号
		// var PrivateKey string
		// if PrivateKey, err = fileutil.ReadFileContent(config.Pay.Weixinpaybill.Apiclientkeyfile); err != nil {
		// 	return
		// }
		// client, err = wechat.NewClientV3(config.Pay.Weixinpaybill.Mchid, config.Pay.Weixinpaybill.Serialno, config.Pay.Weixinpaybill.Apiv3key, PrivateKey)
	} else {
		if payInfo == nil { // 使用系统微信支付账号
			var privateKey string
			if privateKey, err = fileutil.ReadFileContent(config.Pay.Weixinpay.Apiclientkeyfile); err != nil {
				return
			}
			client, err = wechat.NewClientV3(config.Pay.Weixinpay.Mchid, config.Pay.Weixinpay.Serialno, config.Pay.Weixinpay.Apiv3key, privateKey)
		} else { // 使用指定的微信支付账号
			// fmt.Println("使用商户支付账号：")
			// fmt.Printf("%+v \n\n", payInfo)
			client, err = wechat.NewClientV3(payInfo.MchId, payInfo.SerialNo, payInfo.Apiv3Key, payInfo.ApiClientKey)
		}
	}

	if err != nil {
		xlog.Error(err)
		return
	}

	// 设置微信平台API证书和序列号（推荐开启自动验签，无需手动设置证书公钥等信息）
	//client.SetPlatformCert([]byte(""), "")

	// 启用自动同步返回验签，并定时更新微信平台API证书（开启自动验签时，无需单独设置微信平台API证书和序列号）
	err = client.AutoVerifySign()
	if err != nil {
		xlog.Error(err)
		return
	}

	// 自定义配置http请求接收返回结果body大小，默认 10MB
	// client.SetBodySize() // 没有特殊需求，可忽略此配置

	// 打开Debug开关，输出日志，默认是关闭的
	client.DebugSwitch = gopay.DebugOn

	return
}

// doWeixinPay 执行微信支付；
// payChannel app/min
func doWeixinPay(payChannel string, payPrice float64, orderSN, body, attach string, openId string, payInfo *model.WechatPayInfo) (res interface{}, err error) {
	var (
		strs       = strings.Split(attach, "-")
		merchantID string
	)
	if len(strs) == 2 {
		attach = strs[0]
		merchantID = strs[1]
	}
	client, err := getWechatClient(attach, payInfo)
	if err != nil {
		return
	}

	var (
		notifyUrl string
		expire    = time.Now().Add(10 * time.Minute).Format(time.RFC3339)
		bm        = make(gopay.BodyMap) // 初始化 BodyMap
		wxRsp     *wechat.PrepayRsp
		ctx       = context.Background()
		prepayId  string
	)
	if attach == "bill" {
		// notifyUrl = config.App.Fronturl + "/wapapi/notify/weixinpaybill"
	} else {
		notifyUrl = config.App.Fronturl + "/wapapi/notify/weixinpay/" + merchantID
	}

	if payChannel == "app" {
		bm.
			Set("appid", config.Weixin.Weixinapp.Appid).
			Set("time_expire", expire).
			Set("description", util.Substr(body, 0, 42)).
			Set("attach", attach).
			Set("out_trade_no", orderSN).
			Set("notify_url", notifyUrl).
			SetBodyMap("amount", func(bm gopay.BodyMap) {
				bm.Set("total", decimal.NewFromFloat(payPrice).Mul(decimal.NewFromInt(100)).IntPart())
				bm.Set("currency", "CNY")
			})
		wxRsp, err = client.V3TransactionApp(ctx, bm)
	} else if payChannel == "min" { // 小程序
		bm.
			Set("appid", config.Weixin.Weixinmin.Appid).
			Set("time_expire", expire).
			Set("description", util.Substr(body, 0, 42)).
			Set("attach", attach).
			Set("out_trade_no", orderSN).
			Set("notify_url", notifyUrl).
			SetBodyMap("amount", func(bm gopay.BodyMap) {
				bm.Set("total", decimal.NewFromFloat(payPrice).Mul(decimal.NewFromInt(100)).IntPart())
				bm.Set("currency", "CNY")
			}).
			SetBodyMap("payer", func(bm gopay.BodyMap) {
				bm.Set("openid", openId)
			})
		wxRsp, err = client.V3TransactionJsapi(ctx, bm)
	} else {
		return nil, errors.New("payChannel: app/min")
	}

	if err != nil {
		xlog.Error(err)
		return
	}
	if wxRsp.Code != wechat.Success {
		xlog.Errorf("wxRsp:%s", wxRsp.Error)
		return nil, errors.New(wxRsp.Error)

	}
	// xlog.Debugf("wxRsp: %#v", wxRsp.Response)
	prepayId = wxRsp.Response.PrepayId
	if payChannel == "app" {
		var appResult *wechat.AppPayParams
		appResult, err = client.PaySignOfApp(config.Weixin.Weixinapp.Appid, prepayId)
		if err != nil {
			return
		}

		type AppPayParams2 struct {
			*wechat.AppPayParams
			PaySign string `json:"paySign"`
		}
		res = &AppPayParams2{AppPayParams: appResult, PaySign: appResult.Sign}
	} else if payChannel == "min" {
		res, err = client.PaySignOfApplet(config.Weixin.Weixinmin.Appid, prepayId)
		if err != nil {
			return
		}
	}
	return
}

// NotifyWeixinpay 微信支付回调处理
func NotifyWeixinpay(request *http.Request, merchantID uint32) (interface{}, error) {

	notifyReq, err := wechat.V3ParseNotify(request)
	if err != nil {
		xlog.Error(err)
		return nil, err
	}

	var (
		wechatPayInfo = GetMerchantWechatPayInfo(merchantID)
		apiV3Key      string
	)

	if wechatPayInfo != nil {
		apiV3Key = wechatPayInfo.Apiv3Key
	} else {
		apiV3Key = config.Pay.Weixinpay.Apiv3key
	}

	client, err := getWechatClient("", wechatPayInfo)
	if err != nil {
		return nil, err
	}

	// 获取微信平台证书
	certMap := client.WxPublicKeyMap()
	// 验证异步通知的签名
	err = notifyReq.VerifySignByPKMap(certMap)
	if err != nil {
		xlog.Error(err)
		return nil, err
	}

	result, err := notifyReq.DecryptCipherText(apiV3Key)
	if err != nil {
		return nil, err
	}

	if result.TradeState == "SUCCESS" {
		// fmt.Println("result.Amount.Total", result.Amount.Total)
		var (
			attach  = result.Attach
			orderSN = result.OutTradeNo
			// payPrice = float64(result.Amount.Total / 100)
			payPrice = getWeixinPrice(result.Amount.Total)
			paySN    = result.TransactionId
			payType  string
		)
		if result.TradeType == "APP" {
			payType = "weixinapppay"
		} else if result.TradeType == "JSAPI" {
			payType = "weixinminpay"
		}
		// fmt.Printf("转换后回调付款金额 %f", payPrice)
		if err := doNotify(attach, orderSN, paySN, payType, payPrice); err != nil {
			return nil, err
		}
	}

	return &wechat.V3NotifyRsp{Code: gopay.SUCCESS, Message: "成功"}, nil
}

func getWeixinPrice(price int) float64 {
	quantity := decimal.NewFromInt(int64(price))
	new, _ := quantity.Div(decimal.NewFromInt(100)).Round(2).Float64()
	return new
}

// doNotify 通知处理
func doNotify(attach, orderSN, paySN, payType string, payPrice float64) error {
	switch attach {
	case "product":
		return nil // orderPaid(orderSN, paySN, payType, payPrice)
	case "bill":
		return nil //billPaid(orderSN, paySN, payType, payPrice)
	default:
		return errors.New("未知订单类型")
	}
}

// NotifyWeixinpay 微信支付回调处理（账单支付）
func NotifyWeixinpaybill(request *http.Request) (interface{}, error) {
	/*
		notifyReq, err := wechat.V3ParseNotify(request)
		if err != nil {
			xlog.Error(err)
			return nil, err
		}

		client, err := getWechatClient("bill", nil)
		if err != nil {
			return nil, err
		}

		// 获取微信平台证书
		certMap := client.WxPublicKeyMap()
		// 验证异步通知的签名
		err = notifyReq.VerifySignByPKMap(certMap)
		if err != nil {
			xlog.Error(err)
			return nil, err
		}

		result, err := notifyReq.DecryptCipherText(config.Pay.Weixinpaybill.Apiv3key)
		if err != nil {
			return nil, err
		}

		if result.TradeState == "SUCCESS" {
			// fmt.Println("result.Amount.Total", result.Amount.Total)
			var (
				attach  = result.Attach
				orderSN = result.OutTradeNo
				// payPrice = float64(result.Amount.Total / 100)
				payPrice = getWeixinPrice(result.Amount.Total)
				paySN    = result.TransactionId
				payType  string
			)
			if result.TradeType == "APP" {
				payType = "weixinapppay"
			} else if result.TradeType == "JSAPI" {
				payType = "weixinminpay"
			}
			// fmt.Printf("转换后回调付款金额 %f", payPrice)
			if err := doNotify(attach, orderSN, paySN, payType, payPrice); err != nil {
				return nil, err
			}
		}
	*/
	return &wechat.V3NotifyRsp{Code: gopay.SUCCESS, Message: "成功"}, nil
}

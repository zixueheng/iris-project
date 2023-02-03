/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2022-12-14 15:14:36
 * @LastEditTime: 2023-02-03 10:57:42
 */
package service

import (
	"context"
	"errors"
	"fmt"
	"iris-project/app/config"
	"iris-project/app/model"
	"iris-project/lib/util"
	"net/http"
	"strings"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/xlog"
	"github.com/go-pay/gopay/wechat/v3"
	"github.com/shopspring/decimal"
)

// DoRefund 退款，
// refundSN: product-order.id-pay.id
func DoRefund(payPrice, refundPrice float64, payType, orderSN, paySN, refundSN string, merchantID uint32, wechatPayInfo *model.WechatPayInfo) error {
	switch payType {
	case "alipay":
		return doAlipayRefund(refundPrice, orderSN, paySN, refundSN)
	case "weixinapppay":
		return doWeixinRefund("app", payPrice, refundPrice, orderSN, paySN, refundSN, merchantID, wechatPayInfo)
	case "weixinminpay":
		return doWeixinRefund("min", payPrice, refundPrice, orderSN, paySN, refundSN, merchantID, wechatPayInfo)
	case "unionpay":
	default:
		return errors.New("支付方式错误")
	}
	return nil
}

// doAlipayRefund 支付宝退款，实时得到结果
func doAlipayRefund(refundPrice float64, orderSN, paySN, refundSN string) error {
	return errors.New("未实现支付宝退款")
}

// doWeixinRefund 微信退款（只实现商品退款，账单不退款）
// payChannel app/min
func doWeixinRefund(payChannel string, payPrice, refundPrice float64, orderSN, paySN, refundSN string, merchantID uint32, wechatPayInfo *model.WechatPayInfo) error {
	client, err := getWechatClient("", wechatPayInfo)
	if err != nil {
		return err
	}

	var (
		bm  = make(gopay.BodyMap)
		ctx = context.Background()
	)
	bm.Set("transaction_id", paySN).
		Set("out_refund_no", refundSN).
		Set("notify_url", config.App.Fronturl+"/wapapi/notify/weixinrefund/"+util.ParseString(int(merchantID))).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("refund", decimal.NewFromFloat(refundPrice).Mul(decimal.NewFromInt(100)).IntPart())
			bm.Set("total", decimal.NewFromFloat(payPrice).Mul(decimal.NewFromInt(100)).IntPart())
			bm.Set("currency", "CNY")
		})

	if res, err := client.V3Refund(ctx, bm); err != nil {
		return err
	} else if res.Code == 0 {
		return nil
	} else {
		return fmt.Errorf("微信退款失败：%s", res.Error)
	}
}

// NotifyWeixinRefund 微信退款通知
func NotifyWeixinRefund(request *http.Request, merchantID uint32) (interface{}, error) {

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

	result, err := notifyReq.DecryptRefundCipherText(apiV3Key)
	if err != nil {
		return nil, err
	}

	if result.RefundStatus == "SUCCESS" {
		// var (
		// 	orderSN     = result.OutTradeNo
		// 	refundSN    = result.OutRefundNo
		// 	refundPrice = float64(result.Amount.Refund / 100)
		// 	split       = strings.Split(refundSN, "-")
		// )
		// if len(split) != 3 {
		// 	return nil, fmt.Errorf("退单号分割错误：%s", refundSN)
		// }
		// if err := doWeixinRefund(split[0], orderSN, refundSN, refundPrice); err != nil {
		// 	return nil, err
		// }
		// fmt.Println("微信退款成功")
	} else {
		var (
			// orderSN     = result.OutTradeNo
			refundSN = result.OutRefundNo
			// refundPrice = float64(result.Amount.Refund / 100)
			split = strings.Split(refundSN, "-")
			msg   = fmt.Sprintf("微信退款处理失败：%s", result.RefundStatus)
		)
		if len(split) != 3 {
			return nil, fmt.Errorf("退单号分割错误：%s", refundSN)
		}
		if err := refundFailed(split[0], msg, uint32(util.ParseInt(split[1])), uint32(util.ParseInt(split[2]))); err != nil {
			return nil, err
		}
	}
	return &wechat.V3NotifyRsp{Code: gopay.SUCCESS, Message: "成功"}, nil
}

// 退款失败处理
func refundFailed(attach, msg string, orderID, payID uint32) error {
	switch attach {
	case "product":
		return nil // orderRefundFailed(msg, orderID, payID)
	case "bill":
	}
	return nil
}

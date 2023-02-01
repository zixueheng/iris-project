/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:28:08
 * @LastEditTime: 2021-12-14 10:33:40
 */
package service

import (
	"errors"
	"fmt"
	"iris-project/app/config"
	"iris-project/app/dao"
	"iris-project/app/model"
	"iris-project/global"
	"iris-project/lib/util"
	"math"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	tencentclouderrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711" //引入sms
	jpushclient "github.com/ylywyn/jpush-api-go-client"
)

// VerifycodeGen 验证码生成
func VerifycodeGen(phone, purpose string) (code string, err error) {
	code = util.GetRandomNumber(6)
	dao.DeleteAll(nil, &model.Verifycode{}, map[string]interface{}{"phone": phone, "purpose": purpose})
	verifycode := &model.Verifycode{
		Phone:   phone,
		Code:    code,
		Expired: global.LocalTime{Time: time.Now().Add(time.Minute * global.VerifycodeLifeTime), Valid: true},
		Purpose: purpose,
	}
	err = dao.SaveOne(nil, verifycode)
	if err != nil {
		return
	}
	dao.DeleteAll(nil, &model.Verifycode{}, map[string]interface{}{"expired <": util.TimeFormat(time.Now(), "")})

	err = SendSMS(config.Tencent.Sms.VerifycodeTemplete, []string{phone}, []string{code, util.ParseString(global.VerifycodeLifeTime)})
	if err != nil {
		return
	}
	return
}

// VerifycodeCheck 验证码检查
func VerifycodeCheck(phone, code, purpose string) error {
	verifycode := &model.Verifycode{}
	dao.FindOne(nil, verifycode, map[string]interface{}{"phone": phone, "code": code, "purpose": purpose})
	if verifycode.Code == "" {
		return errors.New("验证码不存在")
	}
	if verifycode.Expired.Time.Before(time.Now()) {
		return errors.New("验证码已过期")
	}
	return nil
}

// SendSMS 发送短信，手机号200个以内
func SendSMS(templateID string, phones []string, param []string) error {
	/* 必要步骤：
	 * 实例化一个认证对象，入参需要传入腾讯云账户密钥对 secretId 和 secretKey
	 * 本示例采用从环境变量读取的方式，需要预先在环境变量中设置这两个值
	 * 您也可以直接在代码中写入密钥对，但需谨防泄露，不要将代码复制、上传或者分享给他人
	 * CAM 密匙查询: https://console.cloud.tencent.com/cam/capi
	 */
	credential := common.NewCredential(
		// os.Getenv("TENCENTCLOUD_SECRET_ID"),
		// os.Getenv("TENCENTCLOUD_SECRET_KEY"),
		config.Tencent.SecretID,
		config.Tencent.SecretKey,
	)
	/* 非必要步骤:
	 * 实例化一个客户端配置对象，可以指定超时时间等配置 */
	cpf := profile.NewClientProfile()

	/* SDK 默认使用 POST 方法
	 * 如需使用 GET 方法，可以在此处设置，但 GET 方法无法处理较大的请求 */
	cpf.HttpProfile.ReqMethod = "POST"

	/* SDK 有默认的超时时间，非必要请不要进行调整
	 * 如有需要请在代码中查阅以获取最新的默认值 */
	//cpf.HttpProfile.ReqTimeout = 5

	/* SDK 会自动指定域名，通常无需指定域名，但访问金融区的服务时必须手动指定域名
	 * 例如 SMS 的上海金融区域名为 sms.ap-shanghai-fsi.tencentcloudapi.com */
	// cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"

	/* SDK 默认用 TC3-HMAC-SHA256 进行签名，非必要请不要修改该字段 */
	cpf.SignMethod = "HmacSHA1"

	/* 实例化 SMS 的 client 对象
	 * 第二个参数是地域信息，可以直接填写字符串 ap-guangzhou，或者引用预设的常量 */
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)

	/* 实例化一个请求对象，根据调用的接口和实际情况，可以进一步设置请求参数
	   * 您可以直接查询 SDK 源码确定接口有哪些属性可以设置
	    * 属性可能是基本类型，也可能引用了另一个数据结构
	    * 推荐使用 IDE 进行开发，可以方便地跳转查阅各个接口和数据结构的文档说明 */
	request := sms.NewSendSmsRequest()

	/* 基本类型的设置:
	 * SDK 采用的是指针风格指定参数，即使对于基本类型也需要用指针来对参数赋值。
	 * SDK 提供对基本类型的指针引用封装函数
	 * 帮助链接：
	 * 短信控制台：https://console.cloud.tencent.com/smsv2
	 * sms helper：https://cloud.tencent.com/document/product/382/3773
	 */

	/* 短信应用 ID: 在 [短信控制台] 添加应用后生成的实际 SDKAppID，例如1400006666 */
	request.SmsSdkAppid = common.StringPtr(config.Tencent.Sms.AppID)
	/* 短信签名内容: 使用 UTF-8 编码，必须填写已审核通过的签名，可登录 [短信控制台] 查看签名信息 */
	request.Sign = common.StringPtr(config.Tencent.Sms.Sign)
	/* 国际/港澳台短信 senderid: 国内短信填空，默认未开通，如需开通请联系 [sms helper] */
	// request.SenderId = common.StringPtr("xxx")
	/* 用户的 session 内容: 可以携带用户侧 ID 等上下文信息，server 会原样返回 */
	// request.SessionContext = common.StringPtr("xxx")
	/* 短信码号扩展号: 默认未开通，如需开通请联系 [sms helper] */
	// request.ExtendCode = common.StringPtr("0")
	/* 模板参数: 若无模板参数，则设置为空*/
	request.TemplateParamSet = common.StringPtrs(param)
	/* 模板 ID: 必须填写已审核通过的模板 ID，可登录 [短信控制台] 查看模板 ID */
	request.TemplateID = common.StringPtr(templateID)
	/* 下发手机号码，采用 e.164 标准，+[国家或地区码][手机号]
	 * 例如+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号，最多不要超过200个手机号*/
	for i, phone := range phones {
		phones[i] = "+86" + phone
	}
	request.PhoneNumberSet = common.StringPtrs(phones)

	// 通过 client 对象调用想要访问的接口，需要传入请求对象
	// response, err := client.SendSms(request)
	_, err := client.SendSms(request)
	// 处理异常
	if _, ok := err.(*tencentclouderrors.TencentCloudSDKError); ok {
		// fmt.Printf("An API error has returned: %s", err)
		return err
	}
	// 非 SDK 异常，直接失败。实际代码中可以加入其他的处理
	if err != nil {
		return err
	}
	// b, _ := json.Marshal(response.Response)
	// 打印返回的 JSON 字符串
	// fmt.Printf("%s", b)
	return nil
}

// SendSMSMore 多手机号发送短信
func SendSMSMore(templateID string, phones []string, param []string) error {
	var (
		maxNum = 200
		mod    = len(phones) % maxNum
		floor  = math.Floor(float64(len(phones)) / float64(maxNum))
		i      = 0
	)

	// fmt.Println(mod, floor)
	for ; i < int(floor); i++ {
		// fmt.Println(phones[i*maxNum : (i+1)*maxNum])
		SendSMS(templateID, phones[i*maxNum:(i+1)*maxNum], param)
	}
	// fmt.Println(i)
	if mod > 0 {
		// fmt.Println(phones[i*maxNum : i*maxNum+mod])
		SendSMS(templateID, phones[i*maxNum:i*maxNum+mod], param)
	}

	return nil
}

// Jpush 极光推送
func Jpush(alias []string, title string, extras map[string]interface{}) error {
	//Platform
	var pf jpushclient.Platform
	pf.Add(jpushclient.ANDROID)
	pf.Add(jpushclient.IOS)
	// pf.Add(jpushclient.WINPHONE)
	//pf.All()

	//Audience
	var ad jpushclient.Audience
	// s := []string{"1", "2", "3"}
	// ad.SetTag(s)
	ad.SetAlias(alias)
	// ad.SetID(s)
	//ad.All()

	//Notice
	var notice jpushclient.Notice
	notice.SetAlert(title)
	notice.SetAndroidNotice(&jpushclient.AndroidNotice{Alert: title, Extras: extras})
	notice.SetIOSNotice(&jpushclient.IOSNotice{Alert: title, Extras: extras})
	// notice.SetWinPhoneNotice(&jpushclient.WinPhoneNotice{Alert: "WinPhoneNotice"})

	// var msg jpushclient.Message
	// msg.Title = "Hello"
	// msg.Content = "收到吗？"

	payload := jpushclient.NewPushPayLoad()
	payload.SetPlatform(&pf)
	payload.SetAudience(&ad)
	// payload.SetMessage(&msg)
	payload.SetNotice(&notice)

	bytes, _ := payload.ToBytes()
	fmt.Printf("%s\r\n", string(bytes))

	//push
	c := jpushclient.NewPushClient(config.Jpush.Mastersecret, config.Jpush.Appkey)
	// str, err := c.Send(bytes)
	_, err := c.Send(bytes)
	if err != nil {
		return err
		// fmt.Printf("err:%s", err.Error())
	} else {
		// fmt.Printf("ok:%s", str)
	}
	return nil
}

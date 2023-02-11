/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-02-01 16:15:07
 * @LastEditTime: 2023-02-11 11:33:08
 */
package config

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/configor"
)

var (
	// App 项目配置
	App = struct {
		Appname    string
		Fronturl   string
		HTTPS      bool
		Certpath   string
		Certkey    string
		Host       string
		Port       uint
		Timeformat string
		Dateformat string
		Jwtsecret  string
		LogLevel   string
	}{}

	// DB 数据库配置
	DB = struct {
		Adapter string

		// 主库
		Host        string
		Port        uint
		Name        string
		User        string
		Password    string
		TablePrefix string
		// Debug    bool
		LogLevel string

		// 从库
		Replicas []struct {
			Host     string
			Port     uint
			Name     string
			User     string
			Password string
		}
	}{}

	// ES elasticsearch配置
	ES = struct {
		On        bool
		Addresses []string
		Username  string
		Password  string
	}{}

	// Jpush 极光推送
	Jpush = struct {
		Appkey       string
		Mastersecret string
	}{}

	// Pay支付配置
	Pay = struct {
		Debug  bool
		Alipay struct {
			Appid      string
			Publickey  string
			Privatekey string
			Isprod     bool
		}

		Weixinpay struct {
			// Appid             string
			Mchid             string
			Serialno          string
			Apiv3key          string
			Apiclientcertfile string
			Apiclientkeyfile  string
			Newcertfile       string
		}
	}{}

	// RabbitMQ 消息队列配置
	RabbitMQ = struct {
		IP       string
		Username string
		Password string
		Port     int
	}{}

	// Redis 数据库配置
	Redis = struct {
		Host     string
		Password string
		DB       int
		PoolSize int
	}{}

	// Tencent 腾讯云配置
	Tencent = struct {
		AppID        string
		SecretID     string
		SecretKey    string
		CaptchaAppId uint64
		AppSecretKey string
		Sms          struct {
			AppID              string
			Appkey             string
			Sign               string
			VerifycodeTemplete string
		}
	}{}

	// Upload 上传配置
	Upload = struct {
		Storage   string
		Maxsize   int64
		SecretKey string
		Local     struct {
			Domain string
		}
		Qiniu struct {
			Domain    string
			Accesskey string
			Secretkey string
			Bucket    string
			Area      string
		}
		Tencent struct {
			Domain    string
			Accesskey string
			Secretkey string
			Bucket    string
			Area      string
		}
		Aliyun struct {
			Domain    string
			Accesskey string
			Secretkey string
			Bucket    string
			Area      string
		}
	}{}

	// Weixin 微信配置
	Weixin = struct {
		Debug     bool
		Weixinapp struct {
			Appid     string
			Appsecret string
		}
		Weixinmin struct {
			Appid     string
			Appsecret string
		}
	}{}
)

func init() {
	rootpath, _ := os.Getwd()
	// fmt.Println(rootpath)
	if err := configor.New(&configor.Config{Debug: false}).Load(&App, filepath.Join(rootpath, "config/app.yml")); err != nil {
		panic(err)
	}

	if err := configor.New(&configor.Config{Debug: false}).Load(&DB, filepath.Join(rootpath, "config/db.yml")); err != nil {
		panic(err)
	}

	if err := configor.New(&configor.Config{Debug: false}).Load(&ES, filepath.Join(rootpath, "config/es.yml")); err != nil {
		panic(err)
	}

	if err := configor.New(&configor.Config{Debug: false}).Load(&Jpush, filepath.Join(rootpath, "config/jpush.yml")); err != nil {
		panic(err)
	}

	if err := configor.New(&configor.Config{Debug: false}).Load(&Pay, filepath.Join(rootpath, "config/pay.yml")); err != nil {
		panic(err)
	}

	if err := configor.New(&configor.Config{Debug: false}).Load(&RabbitMQ, filepath.Join(rootpath, "config/rabbitmq.yml")); err != nil {
		panic(err)
	}

	if err := configor.New(&configor.Config{Debug: false}).Load(&Redis, filepath.Join(rootpath, "config/redis.yml")); err != nil {
		panic(err)
	}

	if err := configor.New(&configor.Config{Debug: false}).Load(&Tencent, filepath.Join(rootpath, "config/tencent.yml")); err != nil {
		panic(err)
	}

	if err := configor.New(&configor.Config{Debug: false}).Load(&Upload, filepath.Join(rootpath, "config/upload.yml")); err != nil {
		panic(err)
	}

	if err := configor.New(&configor.Config{Debug: false}).Load(&Weixin, filepath.Join(rootpath, "config/weixin.yml")); err != nil {
		panic(err)
	}

}

type (
	// Version 版本
	Version struct {
		Android VersionInfo
		Ios     VersionInfo
	}
	VersionInfo struct {
		Version     string
		Title       string
		Msgs        []string
		Downloadurl string
	}
)

// GetVersion 获取APP版本
func GetVersion() *Version {
	rootpath, _ := os.Getwd()
	var version = &Version{}
	if err := configor.New(&configor.Config{Debug: false}).Load(version, filepath.Join(rootpath, "config/version.yml")); err != nil {
		panic(err)
	}
	return version
}

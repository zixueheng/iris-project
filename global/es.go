/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-02-11 10:03:54
 * @LastEditTime: 2023-02-11 11:38:23
 */
package global

import (
	"iris-project/app/config"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

// EsClient ElasticSearch客户端
var EsClient *elasticsearch.Client

func init() {
	if !config.ES.On {
		return
	}
	var (
		err error
		cfg = elasticsearch.Config{
			// 有多个节点时需要配置
			Addresses: config.ES.Addresses,
			Username:  config.ES.Username,
			Password:  config.ES.Password,
			// 配置HTTP传输对象
			// Transport: &http.Transport{
			//    //MaxIdleConnsPerHost 如果非零，控制每个主机保持的最大空闲(keep-alive)连接。如果为零，则使用默认配置2。
			//    MaxIdleConnsPerHost:   10,
			//    //ResponseHeaderTimeout 如果非零，则指定在写完请求(包括请求体，如果有)后等待服务器响应头的时间。
			//    ResponseHeaderTimeout: time.Second,
			//    //DialContext 指定拨号功能，用于创建不加密的TCP连接。如果DialContext为nil(下面已弃用的Dial也为nil)，那么传输拨号使用包网络。
			//    DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			//    // TLSClientConfig指定TLS.client使用的TLS配置。
			//    //如果为空，则使用默认配置。
			//    //如果非nil，默认情况下可能不启用HTTP/2支持。
			//    TLSClientConfig: &tls.Config{
			// 	  MaxVersion:         tls.VersionTLS11,
			// 	  //InsecureSkipVerify 控制客户端是否验证服务器的证书链和主机名。
			// 	  InsecureSkipVerify: true,
			//    },
			// },
		}
	)

	if EsClient, err = elasticsearch.NewClient(cfg); err != nil {
		log.Println(err.Error())
		for {
			time.Sleep(60 * time.Second) // 等待60秒再次连接
			log.Println("再次连接Elasticsearch")
			if EsClient, err = elasticsearch.NewClient(cfg); err != nil {
				log.Println(err.Error())
			} else {
				break
			}
		}
	}
	// res, _ := EsClient.Info()
	// defer res.Body.Close()
	// log.Println(res)

}

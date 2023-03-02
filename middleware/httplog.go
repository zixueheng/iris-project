package middleware

import (
	"encoding/json"
	"iris-project/app/config"
	"iris-project/global"
	"iris-project/lib/es"
	"log"
	"strings"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
)

var (
	AcLog            = MakeAccessLog()
	HttpLogIndexName = config.App.Appname + ".http_log"
)

type CustomerLog struct {
	*accesslog.Log
	Params    string `json:"params,omitempty"` // url 参数 ? 后面的 k=v&k=v
	Client    string `json:"client,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	Timestamp string `json:"timestamp"`
}

func init() {
	if !config.ES.On {
		return
	}
	if es.IndexExist(HttpLogIndexName) {
		return
	}
	var mapping = `{
		"mappings": {
			"properties": {
				"timestamp": {
					"type": "date",
					"format": "yyyy-MM-dd HH:mm:ss"
				},
				"latency": {
					"type": "integer"
				},
				"code": {
					"type": "keyword"
				},
				"method": {
					"type": "keyword"
				},
				"path": {
					"type": "keyword"
				},
				"params": {
					"type": "text"
				},
				"ip": {
					"type": "ip"
				},
				"client": {
					"type": "keyword"
				},
				"user_id": {
					"type": "keyword"
				},
				"request": {
					"type": "text"
				},
				"response": {
					"type": "text"
				},
				"bytes_received": {
					"type": "long"
				},
				"bytes_sent": {
					"type": "long"
				}
			}
		}
	}`
	if err := es.CreateIndex(HttpLogIndexName, mapping); err != nil {
		log.Println(err.Error())
	} else {
		log.Println("create elasticsearch index: " + HttpLogIndexName + " ok")
	}
}

// 请求日志写入Elasticsearch
type HttpLogWriter struct {
}

func (logWriter *HttpLogWriter) Write(p []byte) (n int, err error) {
	if !config.ES.On {
		return
	}

	var customLog = &CustomerLog{}
	if err = json.Unmarshal(p, customLog); err != nil {
		log.Println(err.Error())
		return 0, err
	}
	if customLog.Code == 404 {
		return 0, nil
	}
	// 当前只记录admin的请求日志
	// if !util.InArray(
	// 	customLog.Fields.GetString(global.ClientKey),
	// 	[]string{
	// 		global.GetClient(global.AdminAPI),
	// 		/*global.GetClient(global.WapAPI),*/
	// 	},
	// ) {
	// 	return 0, nil
	// }
	customLog.Client = customLog.Fields.GetString(global.ClientKey)
	customLog.UserID = customLog.Fields.GetString(global.UserID)
	customLog.Params = customLog.Fields.GetString("params")
	customLog.Fields = nil

	var bytes []byte
	if bytes, err = json.Marshal(customLog); err != nil {
		log.Println(err.Error())
		return 0, nil
	}

	if err = es.CreateUpdateDocument(HttpLogIndexName, "", string(bytes)); err != nil {
		// log.Panic(err.Error())
		return 0, err
	}
	return len(p), nil
}

// 创建访问日志
func MakeAccessLog() *accesslog.AccessLog {
	// Initialize a new access log middleware.
	// ac := accesslog.File("./access.log")

	httpLogWriter := &HttpLogWriter{}
	ac := accesslog.New(httpLogWriter)

	// Remove this line to disable logging to console:
	// ac.AddOutput(os.Stdout)

	// The default configuration:
	ac.Delim = '|'
	ac.TimeFormat = config.App.Timeformat
	ac.Async = true
	ac.IP = true
	ac.BytesReceivedBody = true
	ac.BytesSentBody = true
	ac.BytesReceived = false
	ac.BytesSent = false
	ac.BodyMinify = true
	ac.RequestBody = true
	ac.ResponseBody = true
	ac.KeepMultiLineError = true
	ac.PanicLog = accesslog.LogHandler

	ac.AddFields(func(ctx iris.Context, fields *accesslog.Fields) {
		fields.Set("params", ctx.Request().URL.RawQuery)
		fields.Set(global.ClientKey, "")
		fields.Set(global.UserID, "")

		for api, client := range global.ClientMap {
			if strings.HasPrefix(ctx.Path(), api) {
				fields.Set(global.ClientKey, client)
				break
			}
		}

		if ctx.Values().Get("jwt") != nil {
			value := ctx.Values().Get("jwt").(*jwt.Token)
			data := value.Claims.(jwt.MapClaims)

			if value, ok := data[global.UserID]; ok {
				fields.Set(global.UserID, value.(string))
			}
		}

	})

	// Default line format if formatter is missing:
	// Time|Latency|Code|Method|Path|IP|Path Params Query Fields|Bytes Received|Bytes Sent|Request|Response|
	//
	// Set Custom Formatter:
	ac.SetFormatter(&accesslog.JSON{
		Indent:    "  ",
		HumanTime: true,
	})
	// ac.SetFormatter(&accesslog.CSV{})
	// ac.SetFormatter(&accesslog.Template{Text: "{{.Code}}"})

	return ac
}

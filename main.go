/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:27:34
 * @LastEditTime: 2023-02-01 16:08:22
 */
package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"iris-project/app/config"
	"iris-project/global"
	"iris-project/lib/file"
	"iris-project/routes"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

// windows编译 go build -ldflags "-s -w -H=windowsgui" -o=iris-project-daemon.exe
func main() {
	// ac := makeAccessLog() // 请求访问日志
	// defer ac.Close() // Close the underline file.

	app := NewApp()

	app.Logger().SetLevel(config.App.LogLevel)

	file := NewLogFile()
	defer file.Close()
	// app.Logger().SetOutput(file) // 这里可以设置日志输出的地方 文件 或控制台 或 其他的地方
	app.Logger().SetOutput(io.MultiWriter(file, os.Stdout)) // 也可以同时输出到多个地方

	// Register the middleware (UseRouter to catch http errors too).
	// app.UseRouter(ac.Handler)

	// 使用日志中间件
	app.Use(logger.New(logger.Config{
		// Status displays status code
		Status: true,
		// IP displays request's remote address
		IP: true,
		// Method displays the http method
		Method: true,
		// Path displays the request path
		Path: true,
		// Query appends the url query to the Path.
		Query: true,

		// if !empty then its contents derives from `ctx.Values().Get("logger_message")
		// will be added to the logs.
		// MessageContextKeys: []string{"logger_message"},

		// if !empty then its contents derives from `ctx.GetHeader("User-Agent")
		// MessageHeaderKeys: []string{"User-Agent"},
	}))

	routes.InitRoute(app) // 加载路由
	// app.Listen(":8080")

	if config.App.HTTPS {
		host := fmt.Sprintf("%s:%d", config.App.Host, 443)
		if err := app.Run(iris.TLS(host, config.App.Certpath, config.App.Certkey)); err != nil {
			app.Logger().Errorf(fmt.Sprintf("服务启动失败: %v", err))
		}
	} else {
		if err := app.Run(
			iris.Addr(fmt.Sprintf("%s:%d", config.App.Host, config.App.Port)),
			iris.WithoutServerError(iris.ErrServerClosed),
			iris.WithOptimizations,
			// iris.WithTimeFormat(time.RFC3339),
			iris.WithTimeFormat(config.App.Timeformat),
		); err != nil {
			app.Logger().Errorf(fmt.Sprintf("服务启动失败: %v", err))
		}
	}
}

// NewApp 创建 iris 实例
func NewApp() *iris.Application {
	app := iris.New() // 返回全新的 *iris.Application 实例

	app.Use(recover.New()) // Recover 会从paincs中恢复并返回 500 错误码

	// 建议：使用Web服务器 服务静态文件
	// app.HandleDir("/static", "./assets", iris.DirOptions{ShowList: true, Gzip: true})

	iris.RegisterOnInterrupt(func() {
		sqlDb, _ := global.Db.DB()
		sqlDb.Close()
	})

	return app
}

// NewLogFile 创建日志文件
func NewLogFile() *os.File {
	path := "./logs/"
	_ = file.CreateFile(path)
	filename := path + time.Now().Format("2006-01-02") + ".log"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("日志记录出错: %v\n", err)
	}

	return f
}

// 创建访问日志
func makeAccessLog() *accesslog.AccessLog {
	// Initialize a new access log middleware.
	ac := accesslog.File("./access.log")
	// Remove this line to disable logging to console:
	ac.AddOutput(os.Stdout)

	// The default configuration:
	ac.Delim = '|'
	ac.TimeFormat = "2006-01-02 15:04:05"
	ac.Async = false
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
/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:27:34
 * @LastEditTime: 2024-07-29 15:48:24
 */
package main

import (
	// stdContext "context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"iris-project/app/config"
	"iris-project/global"
	"iris-project/lib/file"

	"iris-project/middleware"

	"iris-project/routes"

	// Swag 命令生成的目录 (swag init), 必须要导入
	// swag init --parseDependency --parseInternal
	_ "iris-project/docs"

	// "github.com/joho/godotenv"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

// 编译时可以排除swagger文档，以达到缩小包的尺寸的目的（副作用 go run 的时候不包含文档，只能使用gowatch，查看gowatch.yml）
// go build （不包含文档）
// go build -tags "docs"（包含文档）
var swaggerUI context.Handler

// @title		    项目接口文档
// @version		    1.0
// @description	    windows编译：go build -ldflags "-s -w -H=windowsgui" -o=iris-project-daemon.exe
// @termsOfService	http://swagger.io/terms/
//
// @contact.name	heyongliang
// @contact.url	    http://xxx.com
// @contact.email	356126067@qq.com
//
// @license.name	Apache 2.0
// @license.url	    http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host		    localhost:8080
// @BasePath	    /
func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	// log.Printf("appname from .env: %s\n", os.Getenv("appname"))

	app := newApp()

	/*
		file := NewLogFile()
		defer file.Close()
		app.Logger().SetOutput(io.MultiWriter(file, os.Stdout)) // 也可以同时输出到多个地方
	*/

	// 请求日志和websocket有冲突，所以在routes里面单独加载此中间件
	/*
		ac := middleware.MakeAccessLog() // 请求访问日志
		defer ac.Close()                 // Close the underline file.
		// Register the middleware (UseRouter to catch http errors too).
		app.UseRouter(ac.Handler)
	*/

	middleware.InitWebSocket(app) // 普通方式，mvc方式查看routes
	middleware.InitCron()         // 启动定时任务
	routes.InitRoute(app)         // 加载路由

	if swaggerUI != nil {
		log.Println("已加载swagger文档")
		// Register on http://localhost:8080/swagger
		app.Get("/swagger", swaggerUI)
		// And the wildcard one for index.html, *.js, *.css and e.t.c.
		app.Get("/swagger/{any:path}", swaggerUI)
	} else {
		log.Println("未加载swagger文档")
	}

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

// newApp 创建 iris 实例
func newApp() *iris.Application {
	app := iris.New() // 返回全新的 *iris.Application 实例

	app.Use(recover.New()) // Recover 会从paincs中恢复并返回 500 错误码
	// app.Use(iris.Compression) // 开启gzip压缩，会消耗一定CPU资源，一般不压缩或通过nginx压缩

	app.Logger().SetLevel(config.App.LogLevel)
	app.Logger().SetTimeFormat(config.App.Timeformat)
	app.Logger().SetOutput(io.MultiWriter(os.Stdout))

	// 使用http请求日志中间件，注意专门输出请求日志
	// 输出信息如：
	// [INFO] 2023-02-17 16:48:43 200 11.4412ms 127.0.0.1 GET /wapapi/home
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

	// 建议：使用Web服务器 服务静态文件
	// app.HandleDir("/static", "./assets", iris.DirOptions{ShowList: true, Gzip: true})

	iris.RegisterOnInterrupt(func() {
		sqlDb, _ := global.Db.DB()
		sqlDb.Close()

		middleware.AcLog.Close()
		middleware.CloseCron()

		// timeout := 10 * time.Second
		// ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		// defer cancel()
		// // close all hosts.
		// app.Shutdown(ctx)
	})

	return app
}

// newLogFile 创建日志文件
func newLogFile() *os.File {
	path := "./logs/"
	_ = file.CreateFile(path)
	filename := path + time.Now().Format("2006-01-02") + ".log"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("日志记录出错: %v\n", err)
	}

	return f
}

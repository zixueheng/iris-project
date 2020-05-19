package main

import (
	"fmt"
	"io"
	"os"
	"time"

	adminapimodel "iris-project/app/adminapi/model"
	"iris-project/config"
	"iris-project/global"
	"iris-project/lib/file"
	"iris-project/routes"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func main() {
	app := NewApp()

	app.Logger().SetLevel("debug")

	file := NewLogFile()
	defer file.Close()
	// app.Logger().SetOutput(file) // 这里可以设置日志输出的地方 文件 或控制台 或 其他的地方
	app.Logger().SetOutput(io.MultiWriter(file, os.Stdout)) // 也可以同时输出到多个地方

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

	app.HandleDir("/static", "./assets", iris.DirOptions{ShowList: true, Gzip: true})

	db := global.Db
	db.AutoMigrate(
		// &model.User{},
		&adminapimodel.Menu{},
		&adminapimodel.Role{},
		&adminapimodel.AdminUser{},
	)

	iris.RegisterOnInterrupt(func() {
		_ = global.Db.Close()
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

//go:build docs
// +build docs

// go build -tags 是一个 Go 语言的命令行选项，用于在构建过程中启用或禁用特定的编译标签。这些标签通常用于控制代码的特定行为，例如启用或禁用某些功能、优化等
// 比如：下面表示，tag1 和 tag2 是你想要启用的标签。你可以根据需要添加更多的标签，用空格分隔。
// go build -tags "tag1 tag2"
// 在你的 Go 代码中，你可以使用 // +build 注释来指定哪些文件应该在特定的标签下编译。

/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:27:34
 * @LastEditTime: 2024-07-29 15:33:07
 */
package main

import (
	"github.com/iris-contrib/swagger"
	"github.com/iris-contrib/swagger/swaggerFiles"
)

// DOC: https://github.com/iris-contrib/swagger
// https://github.com/swaggo/swag/blob/master/README_zh-CN.md
func init() {
	swaggerUI = swagger.Handler(swaggerFiles.Handler,
		swagger.URL("http://localhost:8080/swagger/swagger.json"),
		swagger.DeepLinking(true),
		swagger.Prefix("/swagger"),
		// ref: https://github.com/ostranme/swagger-ui-themes
		// current we support 7 themes
		// theme is a optional config, if you not set, it will use default theme
		swagger.SetTheme(swagger.Monokai),
	)
}

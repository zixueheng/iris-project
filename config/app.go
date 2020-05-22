package config

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/configor"
)

// App 项目配置
var App = struct {
	Appname    string
	HTTPS      bool
	Certpath   string
	Certkey    string
	Host       string
	Port       uint
	Timeformat string
	Jwtsecret  string
	LogLevel   string
}{}

func init() {
	rootpath, _ := os.Getwd()
	// fmt.Println(rootpath)
	if err := configor.New(&configor.Config{Debug: false}).Load(&App, filepath.Join(rootpath, "config/app.yml")); err != nil {
		panic(err)
	}
	// fmt.Printf("%+v\n\n", App)
}

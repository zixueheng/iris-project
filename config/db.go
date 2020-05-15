package config

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/configor"
)

// DB 数据库配置
var DB = struct {
	Adapter  string
	Host     string
	Port     uint
	Name     string
	User     string
	Password string
}{}

func init() {
	rootpath, _ := os.Getwd()
	if err := configor.New(&configor.Config{Debug: false}).Load(&DB, filepath.Join(rootpath, "config/db.yml")); err != nil {
		panic(err)
	}
}

package config

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/configor"
)

// Redis 数据库配置
var Redis = struct {
	Host     string
	Password string
	DB       int
	PoolSize int
}{}

func init() {
	rootpath, _ := os.Getwd()
	if err := configor.New(&configor.Config{Debug: false}).Load(&Redis, filepath.Join(rootpath, "config/redis.yml")); err != nil {
		panic(err)
	}
}

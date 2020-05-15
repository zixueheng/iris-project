package global

import (
	"fmt"

	"iris-project/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/mattn/go-sqlite3"
)

// Db 全局数据库连接
var Db *gorm.DB

func init() {
	connStr := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local", config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Name)
	var err error
	Db, err = gorm.Open(config.DB.Adapter, connStr)
	if err != nil {
		panic(err)
	}

	gorm.DefaultTableNameHandler = func(Db *gorm.DB, defaultTableName string) string {
		return "iris_" + defaultTableName
	}

	// 全局禁用表名复数
	Db.SingularTable(true) // 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响

	Db.DB().SetMaxIdleConns(10)
	Db.DB().SetMaxOpenConns(100)

}

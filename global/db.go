/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:27:34
 * @LastEditTime: 2022-05-05 16:07:08
 */
package global

import (
	"fmt"
	"log"
	"os"
	"time"

	"iris-project/app/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

// Db 全局数据库连接
var Db *gorm.DB

func init() {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.Name)
	var err error
	master := mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         255,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	})
	cfg := &gorm.Config{
		SkipDefaultTransaction:                   true,  // 禁用默认事务
		DisableForeignKeyConstraintWhenMigrating: false, // 迁移时是否创建外键
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.DB.TablePrefix, // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true,                  // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		// Logger: logger.Default.LogMode(logger.Silent),
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second,                  // 慢 SQL 阈值
				LogLevel:      logLevel(config.DB.LogLevel), // Log level
				Colorful:      true,                         // 彩色打印
			},
		),
		// NowFunc: func() time.Time {
		// 	return time.Now().Local()
		// },
		// NowFunc: func() time.Time {
		// 	t, _ := time.ParseInLocation(config.App.Timeformat, time.Now().Format(config.App.Timeformat), time.Local)
		// 	return t
		// },
	}
	Db, err = gorm.Open(master, cfg)
	if err != nil {
		for {
			time.Sleep(60 * time.Second) // 等待60秒再次连接
			log.Println("再次连接数据库")
			Db, err = gorm.Open(master, cfg)
			if err == nil {
				break
			} else {
				log.Println(err.Error())
			}
		}
	}

	// 从库
	var replicas []gorm.Dialector
	for _, replica := range config.DB.Replicas {
		dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", replica.User, replica.Password, replica.Host, replica.Port, replica.Name)
		replicas = append(replicas, mysql.New(mysql.Config{
			DSN:                       dsn,   // DSN data source name
			DefaultStringSize:         255,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据版本自动配置
		}))
	}

	if len(replicas) > 0 {
		Db.Use(
			dbresolver.Register(dbresolver.Config{
				// Sources:  []gorm.Dialector{mysql.Open("db2_dsn")},
				// Replicas: []gorm.Dialector{mysql.Open("db3_dsn"), mysql.Open("db4_dsn")},

				// Sources:  []gorm.Dialector{},
				Replicas: replicas,

				Policy: dbresolver.RandomPolicy{}, // sources/replicas 负载均衡策略
			}).SetConnMaxIdleTime(time.Hour).SetConnMaxLifetime(time.Hour).SetMaxIdleConns(10).SetMaxOpenConns(100),
		)
	}

	sqlDB, _ := Db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
}

// 转换日志级别
// Silent、Error、Warn、Info
func logLevel(level string) logger.LogLevel {
	switch level {
	case "Info":
		return logger.Info
	case "Warn":
		return logger.Warn
	case "Error":
		return logger.Error
	case "Silent":
		return logger.Silent
	}
	return logger.Silent
}

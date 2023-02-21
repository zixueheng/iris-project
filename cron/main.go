/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-07-21 14:51:34
 * @LastEditTime: 2023-02-21 16:34:18
 */
package main

import (
	"io/ioutil"
	"iris-project/app/config"
	"iris-project/lib/util"
	"log"
	"os"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

var lock sync.Mutex

// var i = 1

// 编译后放到根目录运行。。。
// windows编译 go build -ldflags "-s -w -H=windowsgui" -o=jmyjc-backup-daemon.exe
func main() {

	checkCron := cron.New() // 创建一个cron实例

	// 执行定时任务

	if _, err := checkCron.AddFunc("@every 1h", backDB); err != nil { // 每1小时备份数据
		log.Println(err.Error())
	}

	if _, err := checkCron.AddFunc("@daily", deleteOldFiles); err != nil { // 每天清除旧文件
		log.Println(err.Error())
	}

	// 启动/关闭
	checkCron.Start()
	defer checkCron.Stop()
	select {
	// 查询语句，保持程序运行，在这里等同于for{}
	}
}

// PathURL 路径
const PathURL = "db-backup"

func backDB() {

	lock.Lock()
	defer lock.Unlock()

	log.Println("备份数据库")

	// db := dao.GetDB()
	rootpath, _ := os.Getwd()

	_, err := util.BackupMySQLDb(config.DB.Host, util.ParseString(int(config.DB.Port)), config.DB.User, config.DB.Password, config.DB.Name, "", rootpath+"/db-backup/sql-file/")
	if err != nil {
		file := rootpath + "/db-backup/error-log/" + time.Now().Format("2006-01-02") + ".log"
		logFile, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		defer logFile.Close()
		log.SetOutput(logFile) // 将文件设置为log输出的文件
		log.SetPrefix("[错误]")
		log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)

		log.Println(err.Error())
	}
}

func deleteOldFiles() {
	rootpath, _ := os.Getwd()
	rd, err := ioutil.ReadDir(rootpath + "/db-backup/sql-file/")
	if err != nil {
		log.Println("read dir fail:", err)
		return
	}

	day, _ := util.TimeParse(util.TimeFormat(time.Now().AddDate(0, 0, -2), "2006-01-02")+" 00:00:00", "2006-01-02 15:04:05")
	// log.Println("今天", util.TimeFormat(day, "2006-01-02 15:04:05"))

	for _, fi := range rd {
		if !fi.IsDir() {
			fullName := rootpath + "/db-backup/sql-file/" + fi.Name()
			// log.Println(fullName, util.TimeFormat(fi.ModTime(), "2006-01-02 15:04:05"))

			if fi.ModTime().Before(day) {
				// log.Println("选择的文件", fullName, util.TimeFormat(fi.ModTime(), "2006-01-02 15:04:05"))
				os.Remove(fullName)
			}
			log.Println()
		}
	}
}

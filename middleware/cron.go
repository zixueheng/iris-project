package middleware

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

// 定时格式定义：https://pkg.go.dev/github.com/robfig/cron/v3@v3.0.1
var cronTasks = map[string]func(){
	// "@every 3s": func() {
	// 	log.Println("every 3 second clock")
	// },
	"@every 1h": backDB,         // 每1小时备份数据
	"@daily":    deleteOldFiles, // 每天清除旧文件
}

var checkCron *cron.Cron

// InitCron 启动定时任务
func InitCron() {
	checkCron = cron.New() // 创建一个cron实例
	for t, f := range cronTasks {
		if _, err := checkCron.AddFunc(t, f); err != nil {
			log.Println(err.Error())
		}
	}

	// 启动/关闭
	checkCron.Start()
	// defer checkCron.Stop()
	// select {
	// // 查询语句，保持程序运行，在这里等同于for{}
	// }
}

// 关闭定时任务
func CloseCron() {
	checkCron.Stop()
}

var lock sync.Mutex

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

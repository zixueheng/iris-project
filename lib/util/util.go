package util

import (
	"encoding/json"
	"io/ioutil"
	"iris-project/config"
	"math/rand"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jameskeane/bcrypt"
)

// 实用函数

// ParseInt string转换int
func ParseInt(b string) int {
	id, err := strconv.Atoi(b)
	if err != nil {
		panic(err)
	} else {
		return id
	}
}

// ParseString int转换string
func ParseString(b int) string {
	id := strconv.Itoa(b)
	return id
}

// TimeFormat 时间格式化
func TimeFormat(t time.Time, f string) string {
	if len(f) == 0 {
		f = config.App.Timeformat
	}
	return t.Format(f)
}

// Substr 从 位置pos 获取指定长度 length 的 s 的字串
func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

// HashPassword 密码加密
func HashPassword(pwd string) string {
	salt, err := bcrypt.Salt(10)
	if err != nil {
		return ""
	}
	hash, err := bcrypt.Hash(pwd, salt)
	if err != nil {
		return ""
	}

	return hash
}

// GetRandomString 生成随机字符串
func GetRandomString(n int) string {
	const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

// StructToString 结构体转成json字符串
func StructToString(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		return err.Error()
	} else {
		return string(b)
	}
}

// StructToMap 结构体转换成map对象
func StructToMap(obj interface{}) map[string]interface{} {
	k := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < k.NumField(); i++ {
		data[strings.ToLower(k.Field(i).Name)] = v.Field(i).Interface()
	}
	return data
}

// BackupMySQLDb 备份MySql数据库（确保 mysqldump 在系统PATH中）
// @param 	host: 			数据库地址: localhost
// @param 	port:			端口: 3306
// @param 	user:			用户名: root
// @param 	password:		密码: root
// @param 	databaseName:	需要被分的数据库名: test
// @param 	tableName:		需要备份的表名: user
// @param 	sqlPath:		备份SQL存储路径: D:/backup/test/
// @return 	backupPath
func BackupMySQLDb(host, port, user, password, databaseName, tableName, sqlPath string) (string, error) {
	var cmd *exec.Cmd

	if tableName == "" {
		cmd = exec.Command("mysqldump", "--opt", "-h"+host, "-P"+port, "-u"+user, "-p"+password, databaseName)
	} else {
		cmd = exec.Command("mysqldump", "--opt", "-h"+host, "-P"+port, "-u"+user, "-p"+password, databaseName, tableName)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		// log.Fatal(err)
		return "", err
	}

	if err := cmd.Start(); err != nil {
		// log.Fatal(err)
		return "", err
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		// log.Fatal(err)
		return "", err
	}
	now := TimeFormat(time.Now(), "20060102")
	var backupPath string
	if tableName == "" {
		backupPath = sqlPath + databaseName + "_" + now + ".sql"
	} else {
		backupPath = sqlPath + databaseName + "_" + tableName + "_" + now + ".sql"
	}
	err = ioutil.WriteFile(backupPath, bytes, 0644)

	if err != nil {
		// panic(err)
		return "", err
	}
	return backupPath, nil
}

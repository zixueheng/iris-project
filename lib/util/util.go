package util

import (
	"encoding/json"
	"iris-project/config"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jameskeane/bcrypt"
)

// 实用函数

// ParseInt string转换int
func ParseInt(b string, defInt int) int {
	id, err := strconv.Atoi(b)
	if err != nil {
		return defInt
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

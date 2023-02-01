/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:27:34
 * @LastEditTime: 2022-11-05 15:33:06
 */
package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"iris-project/app/config"
	"log"
	"math"
	"math/rand"
	"net"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/jameskeane/bcrypt"
	"github.com/shopspring/decimal"
)

// 实用函数

// ParseInt string转换int
func ParseInt(s string) int {
	if s == "" {
		return 0
	}
	id, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	} else {
		return id
	}
}

// PareFloat32 string 转 float32
func PareFloat32(s string) float32 {
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		panic(err)
	} else {
		return float32(f)
	}
}

// PareFloat64 string 转 float62
func PareFloat64(s string) float64 {
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	} else {
		return f
	}
}

// ParseString int转换string
func ParseString(b int) string {
	id := strconv.Itoa(b)
	return id
}

// ParseInt64String int64转换string
func ParseInt64String(b int64) string {
	id := strconv.FormatInt(b, 10)
	return id
}

// Round 对浮点数四舍五入取整
func Round(x float64) float64 {
	return math.Floor(x + 0/5)
}

// TimeFormat 时间格式化
func TimeFormat(t time.Time, f string) string {
	if len(f) == 0 {
		f = config.App.Timeformat
	}
	return t.Format(f)
}

// TimeParse 时间字符串转时间对象
func TimeParse(t string, f string) (timeObj time.Time, err error) {
	if len(f) == 0 {
		f = config.App.Timeformat
	}
	timeObj, err = time.ParseInLocation(f, t, time.Local)
	return
}

// Substr 从 位置pos 获取指定长度 length 的 s 的字串
func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	// fmt.Println(len(runes), l)
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

// UniqueStringSlice 字符串切片去重
func UniqueStringSlice(strSlice []string) []string {
	result := make([]string, 0, len(strSlice))
	temp := map[string]struct{}{}
	for _, item := range strSlice {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// UniqueIntSlice 数字切片去重
func UniqueIntSlice(strSlice []int) []int {
	result := make([]int, 0, len(strSlice))
	temp := map[int]struct{}{}
	for _, item := range strSlice {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// UniqueUint32Slice 数字切片去重
func UniqueUint32Slice(strSlice []uint32) []uint32 {
	result := make([]uint32, 0, len(strSlice))
	temp := map[uint32]struct{}{}
	for _, item := range strSlice {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// InArray 模拟PHP in_array
// 查找字符是否在数组中，要求同类型，如 InArray("a", []string{"a", "b"})，InArray(uint32(0), []uint32{0, 1, 2})
func InArray(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

// PriceFormat 价格格式化
func PriceFormat(price float64) float64 {
	v, _ := decimal.NewFromFloat(price).Round(2).Float64()
	return v
}

// NumberRound 数字保留小数点几位
func NumberRound(number float64, round int32) float64 {
	v, _ := decimal.NewFromFloat(number).Round(round).Float64()
	return v
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

// MD5 md5加密
func MD5(pwd string) string {
	m := md5.New()
	io.WriteString(m, pwd)
	return hex.EncodeToString(m.Sum(nil))
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

// GetRandomNumber 生成随机数字
func GetRandomNumber(n int) string {
	const letterBytes = "0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

// GetRandomKey 获取随机Key 20060102150405- 88888888
func GetRandomKey() string {
	s := TimeFormat(time.Now(), "20060102150405")
	return s + GetRandomNumber(8)
}

var node, _ = snowflake.NewNode(1)

// CreateSN 雪花法生产单号
func CreateSN(prefix string) string {
	// 下面代码要移到外面去，否则并发生成的单号一样的
	// var node, err := snowflake.NewNode(1)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return ""
	// }

	// Generate a snowflake ID.
	id := node.Generate()
	return prefix + id.String()
}

// CreateSNDatetime 生成日期格式单号
func CreateSNDatetime(rnd int) string {
	return TimeFormat(time.Now(), "20060102-150405") + "-" + GetRandomNumber(rnd)
}

// DistanceFormat 距离格式化显示
func DistanceFormat(meters float64) string {
	if meters < 1 {
		return "<1米"
	} else if meters < 1000 {
		v := decimal.NewFromFloat(meters).Round(0).String()
		return v + "米"
	} else if meters < 10000 {
		v := decimal.NewFromFloat(meters / 1000).Round(2).String()
		return v + "千米"
	} else {
		v := decimal.NewFromFloat(meters / 10000).Round(2).String()
		return v + "万米"
	}
}

// GenFileName 生成文件名 ext拓展名
func GenFileName(ext string) string {
	return TimeFormat(time.Now(), "20060102150405") + "-" + GetRandomNumber(6) + ext
}

// GetOutboundIP 获取对外IP
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
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

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// MAP 转键值对
func GetKV(mp map[string]string) []*KV {
	var res = make([]*KV, 0)
	for k, v := range mp {
		res = append(res, &KV{Key: k, Value: v})
	}
	return res
}

// CloneMap 复制一个Map
func CloneMap(mp map[string]interface{}) map[string]interface{} {
	clone := make(map[string]interface{})
	for k, v := range mp {
		clone[k] = v
	}
	return clone
}

// ResolveTime 将传入的“秒”解析为 1天1小时1分钟1秒
func ResolveTime(seconds int64) string {
	days := seconds / 86400   //转换天数
	seconds = seconds % 86400 //剩余秒数

	hours := seconds / 3600  //转换小时
	seconds = seconds % 3600 //剩余秒数

	minutes := seconds / 60 //转换分钟
	seconds = seconds % 60  //剩余秒数

	if days > 0 {
		return strconv.FormatInt(days, 10) + "天" + strconv.FormatInt(hours, 10) + "小时" + strconv.FormatInt(minutes, 10) + "分" + strconv.FormatInt(seconds, 10) + "秒"
	} else {
		return strconv.FormatInt(hours, 10) + "小时" + strconv.FormatInt(minutes, 10) + "分" + strconv.FormatInt(seconds, 10) + "秒"
	}
}

// GeoDistance 计算地理距离，依次为两个坐标的纬度、经度
// 未测试
func GeoDistance(lng1 float64, lat1 float64, lng2 float64, lat2 float64) float64 {
	const PI float64 = 3.141592653589793
	radlat1 := PI * lat1 / 180
	radlat2 := PI * lat2 / 180

	theta := lng1 - lng2
	radtheta := PI * theta / 180

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515
	dist = dist * 1.609344
	return dist
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
	now := TimeFormat(time.Now(), "20060102150405")
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

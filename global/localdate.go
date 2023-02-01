/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-08 14:00:17
 * @LastEditTime: 2021-04-28 14:48:03
 */
package global

import (
	"iris-project/app/config"
	"strings"
	"time"

	//"strconv"

	"database/sql/driver"
	"fmt"
)

// LocalDate Grom时间
type LocalDate struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// UnmarshalJSON ...
func (t *LocalDate) UnmarshalJSON(data []byte) error {
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	str = strings.Trim(str, "\"")

	if str == "null" || str == "" {
		return nil
	}

	t1, err := time.Parse(config.App.Dateformat, str)
	*t = LocalDate{Time: t1, Valid: true}
	return err
}

// MarshalJSON 查询时执行
func (t LocalDate) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte(fmt.Sprintf("\"%v\"", "")), nil
	}
	formatted := fmt.Sprintf("\"%v\"", time.Unix(t.Time.Unix(), 0).Format(config.App.Dateformat))

	return []byte(formatted), nil
}

// Value 写入数据库时会调用该方法将自定义时间类型转换并写入数据库；
func (t LocalDate) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	formtTime := t.Time.Format(config.App.Dateformat)
	// fmt.Println("格式化时间：", formtTime)
	// fmt.Println()
	return formtTime, nil
}

// Scan 读取数据库时会调用该方法将时间数据转换成自定义时间类型
func (t *LocalDate) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		tmp, _ := time.ParseInLocation(config.App.Dateformat, value.Format(config.App.Dateformat), time.Local)
		*t = LocalDate{Time: tmp, Valid: true}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// String ...
func (t *LocalDate) String() string {
	return fmt.Sprintf("hhh:%s", time.Time(t.Time).String())
}

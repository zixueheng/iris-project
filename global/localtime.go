/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-02 14:44:44
 * @LastEditTime: 2021-04-28 14:48:08
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

// LocalTime Grom时间
type LocalTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// UnmarshalJSON ...
func (t *LocalTime) UnmarshalJSON(data []byte) error {
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	str = strings.Trim(str, "\"")

	if str == "null" || str == "" {
		return nil
	}

	t1, err := time.Parse(config.App.Timeformat, str)
	*t = LocalTime{Time: t1, Valid: true}
	return err
}

// MarshalJSON 查询时执行
func (t LocalTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte(fmt.Sprintf("\"%v\"", "")), nil
	}
	formatted := fmt.Sprintf("\"%v\"", time.Unix(t.Time.Unix(), 0).Format(config.App.Timeformat))

	return []byte(formatted), nil
}

// Value 写入数据库时会调用该方法将自定义时间类型转换并写入数据库；
func (t LocalTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time.Format(config.App.Timeformat), nil
}

// Scan 读取数据库时会调用该方法将时间数据转换成自定义时间类型
func (t *LocalTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		tmp, _ := time.ParseInLocation(config.App.Timeformat, value.Format(config.App.Timeformat), time.Local)
		*t = LocalTime{Time: tmp, Valid: true}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

// String ...
func (t *LocalTime) String() string {
	return fmt.Sprintf("hhh:%s", time.Time(t.Time).String())
}

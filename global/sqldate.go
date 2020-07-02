package global

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"iris-project/config"
	"strings"
	"time"
)

// SQLDate 模型用自定义时间用于替换原来的 time.Time
type SQLDate time.Time

// UnmarshalJSON ...
func (t *SQLDate) UnmarshalJSON(data []byte) error {
	var err error
	//前端接收的时间字符串
	str := string(data)
	str = strings.Trim(str, "\"") // 去除接收的str收尾多余的"

	if string(data) == "null" || str == "" {
		return nil
	}

	t1, err := time.Parse(config.App.Dateformat, str)
	if err != nil {
		fmt.Println("ERR: ", err)
	}
	*t = SQLDate(t1)
	return err
}

// MarshalJSON ...
func (t SQLDate) MarshalJSON() ([]byte, error) {
	// fmt.Println("AAAAA")
	// if time.Time(t).IsZero() {
	// 	return []byte(""), nil
	// }
	formatted := fmt.Sprintf("\"%v\"", time.Time(t).Format(config.App.Dateformat))
	return []byte(formatted), nil
}

// Value ...
func (t SQLDate) Value() (driver.Value, error) {
	// fmt.Println("CCCCC")
	// if time.Time(t).IsZero() {
	// 	return "", nil
	// }
	// SQLDate 转换成 time.Time 类型
	tTime := time.Time(t)
	return tTime.Format(config.App.Dateformat), nil
}

// Scan ...
func (t *SQLDate) Scan(v interface{}) error {
	// fmt.Println("DDDDDD")

	switch vt := v.(type) {
	case time.Time:
		// 字符串转成 time.Time 类型
		*t = SQLDate(vt)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

// String ...
func (t *SQLDate) String() string {
	return fmt.Sprintf("hhh:%s", time.Time(*t).String())
}

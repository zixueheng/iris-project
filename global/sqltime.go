package global

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"iris-project/config"
	"strings"
	"time"
)

// SQLTime 模型用自定义时间用于替换原来的 time.Time
type SQLTime time.Time

// UnmarshalJSON ...
func (t *SQLTime) UnmarshalJSON(data []byte) error {
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	str = strings.Trim(str, "\"")

	if str == "null" || str == "" {
		return nil
	}

	t1, err := time.Parse(config.App.Timeformat, str)
	*t = SQLTime(t1)
	return err
}

// MarshalJSON ...
func (t SQLTime) MarshalJSON() ([]byte, error) {
	// fmt.Println("AAA")
	// if time.Time(t).IsZero() {
	// 	return []byte(""), nil
	// }
	formatted := fmt.Sprintf("\"%v\"", time.Time(t).Format(config.App.Timeformat))
	return []byte(formatted), nil
}

// Value ...
func (t SQLTime) Value() (driver.Value, error) {
	// if time.Time(t).IsZero() {
	// 	return "", nil
	// }
	// SQLTime 转换成 time.Time 类型
	tTime := time.Time(t)
	return tTime.Format(config.App.Timeformat), nil
}

// Scan ...
func (t *SQLTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		// 字符串转成 time.Time 类型
		*t = SQLTime(vt)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

// String ...
func (t *SQLTime) String() string {
	return fmt.Sprintf("hhh:%s", time.Time(*t).String())
}

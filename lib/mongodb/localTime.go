/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-04-26 14:35:53
 * @LastEditTime: 2023-04-26 15:31:26
 */
package mongodb

import (
	"encoding/json"
	"iris-project/app/config"
	"time"
)

// LocalTime represents the BSON datetime value.
// Copy from `primitive.DateTime` and modify Marshal format to `2006-01-02 15:04:05`
type LocalTime int64

var _ json.Marshaler = LocalTime(0)
var _ json.Unmarshaler = (*LocalTime)(nil)

// MarshalJSON marshal to time type.
// Return format 2006-01-02 15:04:05
func (d LocalTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON creates a primitive.DateTime from a JSON string.
func (d *LocalTime) UnmarshalJSON(data []byte) error {
	// Ignore "null" to keep parity with the time.Time type and the standard library. Decoding "null" into a non-pointer
	// DateTime field will leave the field unchanged. For pointer values, the encoding/json will set the pointer to nil
	// and will not defer to the UnmarshalJSON hook.
	if string(data) == "null" {
		return nil
	}

	var tempTime time.Time
	if err := json.Unmarshal(data, &tempTime); err != nil {
		return err
	}

	*d = NewLocalTime(tempTime)
	return nil
}

// Time returns the date as a time type.
func (d LocalTime) Time() time.Time {
	return time.Unix(int64(d)/1000, int64(d)%1000*1000000)
}

// String format time to 2006-01-02 15:04:05
func (d LocalTime) String() string {
	return d.Time().Format(config.App.Timeformat)
}

// NewLocalTime creates a new DateTime from a Time.
func NewLocalTime(t time.Time) LocalTime {
	return LocalTime(t.Unix()*1e3 + int64(t.Nanosecond())/1e6)
}

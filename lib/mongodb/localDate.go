/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-04-26 14:35:53
 * @LastEditTime: 2023-04-26 15:36:15
 */
package mongodb

import (
	"encoding/json"
	"iris-project/app/config"
	"time"
)

// LocalDate represents the BSON datetime value.
// Copy from `primitive.DateTime` and modify Marshal format to `2006-01-02`
type LocalDate int64

var _ json.Marshaler = LocalDate(0)
var _ json.Unmarshaler = (*LocalDate)(nil)

// MarshalJSON marshal to time type.
// Return format 2006-01-02
func (d LocalDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON creates a primitive.DateTime from a JSON string.
func (d *LocalDate) UnmarshalJSON(data []byte) error {
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

	*d = NewLocalDate(tempTime)
	return nil
}

// Time returns the date as a time type.
func (d LocalDate) Time() time.Time {
	return time.Unix(int64(d)/1000, 0)
}

// String format time to 2006-01-02
func (d LocalDate) String() string {
	return d.Time().Format(config.App.Dateformat)
}

// NewLocalDate creates a new LocalDate from a Time.
func NewLocalDate(t time.Time) LocalDate {
	return LocalDate(t.Unix()*1e3 + int64(t.Nanosecond())/1e6)
}

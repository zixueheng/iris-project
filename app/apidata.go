/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:28:08
 * @LastEditTime: 2021-04-28 14:50:18
 */
package app

// Response 数据响应格式
type Response struct {
	Success bool        `json:"success"`
	Code    Code        `json:"code"`
	Msg     interface{} `json:"msg"`
	Data    interface{} `json:"data"`
}

// List 分页数据格式
type List struct {
	List      interface{} `json:"list"`                // 列表
	Total     int64       `json:"total"`               // 总数
	Statistic interface{} `json:"statistic,omitempty"` // 统计数据字段
}

// APIData 返回响应数据
//
// success 是否成功，code 响应码，customMsg自定义（错误码code=999有效），objects响应数据
func APIData(success bool, code Code, customMsg string, objects interface{}) *Response {
	return &Response{
		Success: success,
		Code:    code,
		Msg:     code.String(customMsg),
		Data:    objects,
	}
}

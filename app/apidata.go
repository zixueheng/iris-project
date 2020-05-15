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
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
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

package request

import (
	"testing"
)

func Test_Get(t *testing.T) {
	h := NewHTTPRequest(URLBuild("http://127.0.0.1/test.php", map[string]string{"name": "xiaochuan"}))
	_, err := h.Get()
	if err != nil {
		t.Error("请求错误:", err)
	} else {
		t.Log("正常返回")
	}
}

func Test_Post(t *testing.T) {
	h := NewHTTPRequest("http://127.0.0.1/test.php")
	h.SetBody(map[string]string{"name": "xiaochuan"})
	_, err := h.Post()
	if err != nil {
		t.Error("请求错误:", err)
	} else {
		t.Log("正常返回")
	}
}

func Test_Json(t *testing.T) {
	h := NewHTTPRequest("http://127.0.0.1/test.php")
	h.SetContentType(JSONType)
	h.SetBody(map[string]string{"name": "xiaochuan"})
	_, err := h.Post()
	if err != nil {
		t.Error("请求错误:", err)
	} else {
		t.Log("正常返回")
	}
}

func Benchmark_GET(b *testing.B) {
	for i := 0; i < b.N; i++ {
		h := NewHTTPRequest(URLBuild("http://127.0.0.1/test.php", map[string]string{"name": "xiaochuan"}))
		_, err := h.Get()
		if err != nil {
			b.Error("请求错误:", err)
		} else {
			b.Log("正常返回")
		}
	}
}

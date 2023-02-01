/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2021-02-01 11:27:34
 * @LastEditTime: 2022-10-31 16:20:04
 */
package request

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// 常量
const (
	// 请求类型
	GetMethod    = "GET"
	PostMethod   = "POST"
	DeleteMethod = "DELETE"
	PutMethod    = "PUT"

	// 请求 content-type
	FormType = "form"
	JSONType = "json"
)

// HTTPRequest HTTP请求结构体
type HTTPRequest struct {
	Link        string
	ContentType string
	Header      map[string]string
	Body        interface{}
	sync.RWMutex
}

// NewHTTPRequest 新建一个HTTP请求
func NewHTTPRequest(link string) *HTTPRequest {
	return &HTTPRequest{
		Link:        link,
		ContentType: FormType,
	}
}

// SetBody 设置请求Body
func (h *HTTPRequest) SetBody(body interface{}) {
	h.Lock()
	defer h.Unlock()
	h.Body = body
}

// SetHeader 设置请求Header
func (h *HTTPRequest) SetHeader(header map[string]string) {
	h.Lock()
	defer h.Unlock()
	h.Header = header
}

// SetContentType 设置请求Content类型
func (h *HTTPRequest) SetContentType(contentType string) {
	h.Lock()
	defer h.Unlock()
	h.ContentType = contentType
}

// Get 发送GET请求
func (h *HTTPRequest) Get() ([]byte, error) {
	return h.send(GetMethod)
}

// Post 发送POST请求
func (h *HTTPRequest) Post() ([]byte, error) {
	return h.send(PostMethod)
}

// Put 发送Put请求
func (h *HTTPRequest) Put() ([]byte, error) {
	return h.send(PutMethod)
}

// Delete 发送Delete请求
func (h *HTTPRequest) Delete() ([]byte, error) {
	return h.send(DeleteMethod)
}

// URLBuild 构建URL
func URLBuild(link string, data map[string]string) string {
	u, _ := url.Parse(link)
	q := u.Query()
	for k, v := range data {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// send 执行请求
func (h *HTTPRequest) send(method string) ([]byte, error) {
	var (
		req      *http.Request
		resp     *http.Response
		client   http.Client
		sendData string
		err      error
	)

	if h.Body != nil {
		if strings.ToLower(h.ContentType) == JSONType {
			sendBody, jsonErr := json.Marshal(h.Body)
			if jsonErr != nil {
				return nil, jsonErr
			}
			sendData = string(sendBody)
		} else {
			sendBody := http.Request{}
			sendBody.ParseForm()
			for k, v := range h.Body.(map[string]string) {
				sendBody.Form.Add(k, v)
			}
			sendData = sendBody.Form.Encode()
		}
	}

	//忽略https的证书
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client.Timeout = 30 * time.Second

	req, err = http.NewRequest(method, h.Link, strings.NewReader(sendData))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	//设置默认header
	if len(h.Header) == 0 {
		//json
		if strings.ToLower(h.ContentType) == JSONType {
			h.Header = map[string]string{
				"Content-Type": "application/json;charset=utf-8",
			}
		} else { //form
			h.Header = map[string]string{
				"Content-Type": "application/x-www-form-urlencoded;charset=utf-8",
			}
		}
	}

	for k, v := range h.Header {
		if strings.ToLower(k) == "host" {
			req.Host = v
		} else {
			req.Header.Add(k, v)
		}
	}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error http code :%d", resp.StatusCode) //errors.New(fmt.Sprintf("error http code :%d", resp.StatusCode))
	}

	return ioutil.ReadAll(resp.Body)
}

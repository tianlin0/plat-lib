package curl

import (
	"fmt"
	"net/http"
	"time"
)

// Request 请求变量
type Request struct {
	Url     string        `json:"url"`
	Data    interface{}   `json:"data"`
	Method  string        `json:"method"`
	Header  http.Header   `json:"header"`
	Timeout time.Duration `json:"timeout,omitempty"`
	Cache   time.Duration `json:"cache,omitempty"`
}

func NewRequest(r *Request) *genRequest {
	newRequest := new(genRequest)
	newRequest.Url = r.Url
	newRequest.Data = r.Data
	newRequest.Method = r.Method
	newRequest.Header = r.Header
	newRequest.Timeout = r.Timeout
	newRequest.Cache = r.Cache
	return newRequest
}

type client struct {
	beforeCallback InjectBeforeCallback
	afterCallback  InjectAfterCallback
}

// NewClient 客户端
func NewClient() *client {
	return new(client)
}

// SetCallback 设置执行前方法
func (c *client) SetCallback(beforeCallback InjectBeforeCallback, afterCallback InjectAfterCallback) *client {
	c.beforeCallback = beforeCallback
	c.afterCallback = afterCallback
	return c
}

func (c *client) NewRequest(r *Request) *genRequest {
	gen := NewRequest(r)
	gen.SetCallback(c.beforeCallback, c.afterCallback)
	return gen
}

// RetrieveError 错误类型
type RetrieveError struct {
	Response         *http.Response
	Body             []byte
	ErrorCode        string
	ErrorDescription string
	ErrorURI         string
}

func (r *RetrieveError) Error() string {
	if r.ErrorCode != "" {
		s := fmt.Sprintf("curl: %q", r.ErrorCode)
		if r.ErrorDescription != "" {
			s += fmt.Sprintf(" %q", r.ErrorDescription)
		}
		if r.ErrorURI != "" {
			s += fmt.Sprintf(" %q", r.ErrorURI)
		}
		return s
	}
	return fmt.Sprintf("curl: Status: %v\nResponse: %s", r.Response.Status, r.Body)
}

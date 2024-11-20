package curl

import (
	"context"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/tianlin0/plat-lib/logs"
	"io"
	"net"
	"net/http"
	"time"
)

var (
	//默认超时
	defaultConnectTimeoutSecond  = 20 * time.Second
	defaultTimeoutSecond         = 30 * time.Second
	defaultMaxCoons              = 100
	defaultIdleConnTimeoutSecond = 600 * time.Second
)

// InjectBeforeCallback 发送前的方法
type InjectBeforeCallback func(ctx context.Context, rs *Request, httpReq *http.Request) error

// InjectAfterCallback 发送后的方法
type InjectAfterCallback func(ctx context.Context, rp *Response) error

//默认使用短连接，长连接没有弄好的前提下
func createHTTPClient(timeout time.Duration, trans *http.Transport) *http.Client {
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   defaultConnectTimeoutSecond, // 执行连接超时
			KeepAlive: defaultTimeoutSecond,        // 长连接保持时间
		}).DialContext,
		DisableKeepAlives:   true,                         //关闭长连接
		MaxIdleConnsPerHost: -1,                           //关闭长连接
		MaxIdleConns:        defaultMaxCoons,              // 长连接个数
		IdleConnTimeout:     defaultIdleConnTimeoutSecond, // 长连接保持时间
		DisableCompression:  true,
	}
	if trans != nil {
		if trans.MaxIdleConns == 0 {
			trans.MaxIdleConns = tr.MaxIdleConns
		}
		if trans.IdleConnTimeout == 0 {
			trans.IdleConnTimeout = tr.IdleConnTimeout
		}
		if trans.DialContext == nil {
			trans.DialContext = tr.DialContext
		}
		tr = trans
	}
	client := &http.Client{Transport: tr, Timeout: timeout}
	return client
}

// doRequest 发起实际的HTTP请求
func doRequest(req *http.Request, reqTemp *Request, timeout time.Duration, trans *http.Transport) (status int,
	resData string, resHeader http.Header, err error) {
	time1 := time.Now().UnixNano()

	httpClient := createHTTPClient(timeout, trans)
	req.Close = true

	var resp *http.Response
	resp, err = httpClient.Do(req)
	if err != nil {
		return
	}

	status = resp.StatusCode
	if len(resp.Header) > 0 {
		resHeader = resp.Header
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resData = string(body)

	if status != http.StatusOK {
		time2 := time.Now().UnixNano()
		curlTime := int((time2 - time1) / 1000000)
		msg := fmt.Sprintf("[comm-request do-request][%dms] url=%s, StatusCode=%d, resData=%s, resHeader=%v",
			curlTime, reqTemp.Url, status, resData, resHeader)
		logs.DefaultLogger().Error(msg) //业务逻辑的错，不能算错误
	}
	return status, resData, resHeader, err
}

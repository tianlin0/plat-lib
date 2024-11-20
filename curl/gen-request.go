package curl

import (
	"bytes"
	"context"
	"fmt"
	"github.com/tianlin0/plat-lib/conv"
	"github.com/tianlin0/plat-lib/logs"
	"net/http"
	"strings"
	"time"
)

type genRequest struct {
	Request
	retryPolicy *RetryPolicy

	defaultPrintLogInt int //0 表示默认，只打印一条，1表示完全打开所有信息，-1 表示完全关闭
	transportConfig    *http.Transport
	logger             logs.ILogger
	beforeCallback     InjectBeforeCallback
	afterCallback      InjectAfterCallback
	ctx                context.Context
}

func (p *genRequest) buildGenRequest() {
	if p.Data == nil {
		p.Data = ""
	}
	p.Method = getMethod(p.Method)

	if p.Timeout <= 0 {
		p.Timeout = defaultTimeoutSecond
	}

	if p.Cache > 0 {
		if p.Cache > defaultMaxCacheTime {
			p.Cache = defaultMaxCacheTime
		}
	}

	p.Url = strings.TrimSpace(p.Url)
	p.Header = getHeaders(p.Header, p.Method, p.Data)

	if p.beforeCallback == nil {
		p.beforeCallback = tracePreCallback
	}
	if p.afterCallback == nil {
		p.afterCallback = traceSufCallback
	}

	if p.retryPolicy == nil {
		p.retryPolicy = new(RetryPolicy)
	}
	p.retryPolicy.buildRetryable()

}

func (p *genRequest) getRequest() *Request {
	req := new(Request)
	req.Url = p.Url
	req.Data = p.Data
	req.Method = p.Method
	req.Header = p.Header
	req.Timeout = p.Timeout
	req.Cache = p.Cache
	return req
}

func (p *genRequest) Submit(ctx context.Context) *Response {
	p.buildGenRequest()

	resp := NewResponse(p.getRequest())

	dataString, err := getDataString(p.Data)
	if err != nil {
		resp.Error = err
		return resp
	}

	if p.Url == "" {
		resp.Error = fmt.Errorf("url请求地址为空")
		return resp
	}

	_, err = url.Parse(p.Url)
	if err != nil {
		resp.Error = fmt.Errorf("url格式错误：%s, %v", p.Url, err)
		return resp
	}

	if ctx == nil {
		ctx = p.ctx
	}

	if ctx == nil {
		ctx = context.Background()
	}

	startTime := time.Now()
	if p.Cache > 0 {
		respTxt := getDataFromCache(p)
		if respTxt != "" {
			resp.Response = respTxt
			resp.fromCache = true
			resp.SetDuration(startTime)

			simpleCurlStrData := resp.getLoggerResponse(startTime)
			logStr := fmt.Sprintf("[comm-request cache return]id:%s, data:%s",
				resp.Id, conv.String(simpleCurlStrData))
			printLog(ctx, p.logger, p.defaultPrintLogInt, logStr)

			return resp
		}
	}

	postUrl := getNewPostUrl(p.Url, p.Method, dataString)

	logStr := fmt.Sprintf("[comm-request request] url:%s", postUrl)
	printLog(ctx, p.logger, p.defaultPrintLogInt, logStr)

	startTime = time.Now()
	allResp := p.httpRequest(ctx, dataString, resp)
	allResp.SetDuration(startTime)

	//对返回值进行检查
	if p.retryPolicy != nil {
		isRetry, err := p.retryPolicy.onlyCheckCondition(allResp.Response)
		if err != nil {
			allResp.Error = err
		} else {
			if isRetry {
				allResp.Error = fmt.Errorf(allResp.Response)
			}
		}
	}

	simpleCurlStrData := resp.getLoggerResponse(startTime)
	if p.defaultPrintLogInt >= 0 {
		simpleCurlStrData.printLoggerResponse(ctx)
	}

	if allResp.HttpStatus == http.StatusOK &&
		allResp.Error == nil &&
		allResp.Response != "" &&
		p.Cache > 0 {
		setDataToCache(allResp, p.Cache)
	}

	return allResp
}

func (p *genRequest) getHttpRequest(ctx context.Context, dataString string) (*http.Request, error) {
	postUrl := getNewPostUrl(p.Url, p.Method, dataString)

	httpReq, err := http.NewRequest(p.Method, postUrl, bytes.NewBufferString(dataString))
	if err != nil {
		logStr := fmt.Sprintf("[comm-request request] url:%s, error: %s", postUrl, err.Error())
		printLog(ctx, p.logger, p.defaultPrintLogInt, logStr)
		return nil, err
	}

	if len(p.Header) > 0 {
		for k, v := range p.Header {
			httpReq.Header = setHeaderValues(httpReq.Header, k, v...)
		}
	}

	return httpReq, nil
}

// 递归使用
func (p *genRequest) httpRequest(ctx context.Context, dataString string, resp *Response) *Response {
	httpReq, err := p.getHttpRequest(ctx, dataString)
	if err != nil {
		resp.Error = err
		return resp
	}
	resp.Error = nil

	newRequest := p.getRequest()
	if p.beforeCallback != nil {
		err = p.beforeCallback(ctx, newRequest, httpReq)
		if err != nil {
			resp.Error = err
			return resp
		}
	}

	bodyIsJson := true
	if p.retryPolicy != nil {
		if p.retryPolicy.RespDateType == "string" {
			bodyIsJson = false
		}
	}

	startTime := time.Now()
	{ //直接请求
		status, resData, resHeader, err := doRequest(httpReq, newRequest, p.Timeout, p.transportConfig)
		resp.Response = resData
		resp.Request = newRequest
		resp.HttpStatus = status
		resp.Header = resHeader
		resp.Error = err
		resp.SetDuration(startTime)
	}

	if bodyIsJson {
		if resp.Error == nil {
			var obj interface{}
			err = json.Unmarshal([]byte(resp.Response), &obj)
			if err != nil {
				//返回的不是json格式
				resp.Error = fmt.Errorf("url: %s, response not json: %s，Request=>SetRetryPolicy=>RespDateType=string", newRequest.Url, resp.Response)
			}
		}
	}

	simpleCurlStrData := resp.getLoggerResponse(startTime)
	logStr := fmt.Sprintf("[comm-request http-request return]id:%s, data:%s, error:%v", simpleCurlStrData.Id,
		conv.String(simpleCurlStrData), resp.Error)
	printLog(ctx, p.logger, p.defaultPrintLogInt, logStr)

	if p.afterCallback != nil {
		err = p.afterCallback(ctx, resp)
		if err != nil {
			resp.Error = err
			return resp
		}
	}

	//如果有返回，则判断是否成功。
	isRetry := canRetry(p.retryPolicy, resp.Response)
	if isRetry {
		return p.httpRequest(ctx, dataString, resp)
	}
	return resp
}

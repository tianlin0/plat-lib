package curl

import (
	"context"
	"fmt"
	"github.com/tianlin0/plat-lib/conv"
	"github.com/tianlin0/plat-lib/logs"
	"github.com/tianlin0/plat-lib/utils"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"net/http"
	"time"
)

// Response 方法返回的变量，因为外部方法
type Response struct {
	Id         string        `json:"id"`
	Request    *Request      `json:"request"`
	Response   string        `json:"response"`
	Header     http.Header   `json:"header"`
	HttpStatus int           `json:"status"`
	Duration   time.Duration `json:"duration"` //请求间隔时间
	Error      error         `json:"error"`
	fromCache  bool
}

func (p *Response) SetDuration(startTime time.Time) {
	p.Duration = time.Now().Sub(startTime)
}

func NewResponse(req *Request) *Response {
	reqId := getRequestId(req)
	return &Response{
		Id:      reqId,
		Request: req,
	}
}

type loggerResponse struct {
	Id      string `json:"id"`
	Request struct {
		Url    string      `json:"url"`
		Data   interface{} `json:"data"`
		Method string      `json:"method"`
		Header http.Header `json:"header"`
	} `json:"request"`
	Response   string `json:"response"`
	HttpStatus int    `json:"http_status"`
	Error      error  `json:"error"`
	Times      int64  `json:"times"`
}

func (p *Response) getLoggerResponse(startTime time.Time) *loggerResponse {
	if p.Duration == 0 {
		p.SetDuration(startTime)
	}
	logResp := new(loggerResponse)
	logResp.Id = p.Id
	logResp.Response = p.Response
	logResp.Request.Url = p.Request.Url
	logResp.Request.Data = p.Request.Data

	if len(logResp.Response) > defaultPrintLogDataLength {
		logResp.Response = utils.SubStr(logResp.Response, 0, defaultPrintLogDataLength)
	}
	dataString, err := getDataString(p.Request.Data)
	if err == nil {
		if len(dataString) > defaultPrintLogDataLength {
			logResp.Request.Data = utils.SubStr(dataString, 0, defaultPrintLogDataLength)
		}
	}

	logResp.Request.Method = p.Request.Method
	logResp.Request.Header = p.Request.Header
	logResp.HttpStatus = p.HttpStatus
	logResp.Error = p.Error
	tempTime := time.Now().Sub(startTime)
	logResp.Times = tempTime.Milliseconds()
	return logResp
}

func (p *loggerResponse) printLoggerResponse(ctx context.Context) {
	logLevel := logs.WARNING
	if p.Error != nil {
		logLevel = logs.ERROR
	} else {
		if p.HttpStatus != http.StatusOK && p.HttpStatus != 0 {
			logLevel = logs.ERROR
		}
	}

	returnData := conv.String(p)
	//这里默认打上日志，方便查问题，需要将数据量减少，避免默认内容太多了
	//rData := gjson.Get(returnData, "request.data").String()
	//rHeader := gjson.Get(returnData, "request.header").String()
	//repData := gjson.Get(returnData, "response").String()
	//maxLen := defaultPrintLogDataLength
	//if len(rData) > maxLen {
	//	returnData, _ = sjson.Set(returnData, "request.data", utils.SubStr(rData, 0, maxLen))
	//}
	//if len(rHeader) > maxLen {
	//	returnData, _ = sjson.Set(returnData, "request.headers", utils.SubStr(rHeader, 0, maxLen))
	//}
	//if len(repData) > maxLen {
	//	returnData, _ = sjson.Set(returnData, "response", utils.SubStr(repData, 0, maxLen))
	//}
	logStrTemp := fmt.Sprintf("[comm-request print return]id:%s, data:%s, error: %v", p.Id, returnData, p.Error)
	if logLevel == logs.WARNING {
		logs.CtxLogger(ctx).Warn(logStrTemp)
	} else if logLevel == logs.ERROR {
		logs.CtxLogger(ctx).Error(logStrTemp)
	}
}

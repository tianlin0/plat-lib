package httputil

import (
	"github.com/tianlin0/plat-lib/conv"
	"net/http"
)

// CommResponse 接口返回值
type CommResponse struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Now     string      `json:"now,omitempty"`
	Env     string      `json:"env,omitempty"` //环境
	Time    int64       `json:"time,omitempty"`
	LogId   string      `json:"logid,omitempty"`
	TraceId string      `json:"traceid,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Debug   interface{} `json:"debug,omitempty"`
	Data    interface{} `json:"data"`
}

// PageModel 分页结构输出
type PageModel struct {
	Count     int64       `json:"count"`
	PageNow   int         `json:"pagenow,omitempty"`
	PageStart int         `json:"pagestart,omitempty"`
	PageSize  int         `json:"pagesize,omitempty"`
	PageTotal int         `json:"pagetotal,omitempty"`
	DataList  interface{} `json:"datalist"`
}

// GetCommResponse 获取通用的返回格式
func GetCommResponse(comm *CommResponse) *CommResponse {
	outputMap := &CommResponse{}
	if comm != nil {
		outputMap = comm
	}
	outputMap.Now = conv.FormatFromUnixTime() //当前时间
	return outputMap
}

// WriteCommResponse 将通用返回设置到response，输出到客户端
func WriteCommResponse(respWriter http.ResponseWriter, comm *CommResponse, statusCode ...int) error {
	response := GetCommResponse(comm)

	contentType := "Content-Type"
	respWriter.Header().Set(contentType, "application/json; charset=utf-8")

	respStr := conv.String(response)
	respByte := []byte(respStr)

	oneStatusCode := http.StatusOK
	if len(statusCode) > 0 {
		oneStatusCode = statusCode[0]
	}
	respWriter.WriteHeader(oneStatusCode)

	_, err := respWriter.Write(respByte)

	return err
}

// GetErrorResponse 系统获取错误码和错误信息
func GetErrorResponse(allErrorMap map[int64]string, errorCode int64, err ...error) *CommResponse {
	respError := &CommResponse{}

	respError.Code = errorCode

	if len(err) > 0 {
		respError.Message = err[0].Error()
	}

	if allErrorMap != nil {
		if errorMsg, ok := allErrorMap[errorCode]; ok {
			if respError.Message == "" {
				respError.Message = conv.String(errorMsg)
			}
			return respError
		}
	}

	if respError.Message == "" {
		respError.Message = "系统错误"
	}

	return respError
}

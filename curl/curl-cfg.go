package curl

import (
	"github.com/json-iterator/go"
	"net/http"
	"time"
)

var (
	tracePreCallback InjectBeforeCallback
	traceSufCallback InjectAfterCallback

	defaultMaxCacheTime       = 3600 * 24 * 2 * time.Second //最大用来存2天
	defaultMethod             = http.MethodPost
	defaultReturnKey          = "code"
	defaultReturnVal          = "0"
	defaultPrintLogDataLength = 200 //默认打印日志的时候，数据最长，避免显示太多了

	jsonApi = jsoniter.Config{
		SortMapKeys: true,
	}.Froze()
)

// SetBeforeCallback 设置通用的trace方法
func SetBeforeCallback(injectCallbackTemp InjectBeforeCallback) {
	tracePreCallback = injectCallbackTemp
}

// SetAfterCallback 设置通用的trace方法
func SetAfterCallback(injectCallbackTemp InjectAfterCallback) {
	traceSufCallback = injectCallbackTemp
}

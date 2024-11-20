package curl

import (
	"context"
	"fmt"
	"github.com/ChengjinWu/gojson"
	"github.com/json-iterator/go"
	"github.com/tianlin0/plat-lib/logs"
	"github.com/tianlin0/plat-lib/utils"
	"net/http"
	"net/textproto"
	"net/url"
	"reflect"
	"sort"
	"strings"
)

// setHeaderValues 为headers添加value
func setHeaderValues(h http.Header, key string, values ...string) http.Header {
	if key == "" {
		return h
	}
	for _, v := range values {
		if v == "" {
			continue
		}
		hasValues := h.Values(key)
		isFind := false
		for _, one := range hasValues {
			if one == v {
				isFind = true
			}
		}
		if !isFind {
			h.Add(key, v)
		}
	}
	return h
}

// 需要执行header格式的，如果没有，则直接使用，如果有且value是相同的话，则直接覆盖，避免提交两份数据
func beautifulHeader(headers http.Header) http.Header {
	if headers == nil {
		return nil
	}
	//0,复制一个headers
	oldHeaders := headers.Clone()
	newHeaders := http.Header{}

	// 不写一起的原因是可能后面有key更满足要求的情况。

	//1、首先将header标准key的值取出来
	hasStoreKeyList := make([]string, 0)
	for key, val := range oldHeaders {
		newKey := textproto.CanonicalMIMEHeaderKey(key)
		if key == newKey {
			newHeaders = setHeaderValues(newHeaders, newKey, val...)
			hasStoreKeyList = append(hasStoreKeyList, newKey)
		}
	}

	//2、将已经存储的key删除掉
	for _, key := range hasStoreKeyList {
		oldHeaders.Del(key)
	}

	//3、剩下的看value是否相同，不同的就存储下来
	for key, val := range oldHeaders {
		newKey := textproto.CanonicalMIMEHeaderKey(key)
		allNewValues := newHeaders.Values(newKey)
		if len(allNewValues) == 0 {
			//如果完全不存在，则用new进行存储
			newHeaders = setHeaderValues(newHeaders, newKey, val...)
			continue
		}
		newHeaders = setHeaderValues(newHeaders, key, val...)
	}
	return newHeaders
}

// createParamStrOrder 对参数进行排序，然后拼接成URL的字符串
func createParamStrOrder(params map[string]interface{}) string {
	aParams := make([]string, 0)
	for k, v := range params {
		val := fmt.Sprintf("%v", v)
		aParams = append(aParams, k+"="+url.QueryEscape(val))
	}
	sort.Strings(aParams)
	return strings.Join(aParams, "&")
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	vi := reflect.ValueOf(i)
	kind := vi.Kind()
	if kind == reflect.Ptr ||
		kind == reflect.Chan ||
		kind == reflect.Func ||
		kind == reflect.UnsafePointer ||
		kind == reflect.Map ||
		kind == reflect.Interface ||
		kind == reflect.Slice {
		return vi.IsNil()
	}
	return false
}

func printLog(ctx context.Context, loggers logs.ILogger, printLogInt int, logStr string) {
	//系统外的打印日志
	if !IsNil(loggers) {
		loggers.Debug(logStr)
		return
	}
	// 默认的打印日志
	if printLogInt > 0 {
		if ctx != nil {
			logs.CtxLogger(ctx).Debug(logStr)
		} else {
			logs.DefaultLogger().Debug(logStr)
		}
	}
}

func getNewPostUrl(url, method string, dataString string) string {
	err2 := gojson.CheckValid([]byte(dataString))

	var postUrl = url
	if method == http.MethodGet {
		if dataString != "" {
			param := dataString
			if err2 == nil {
				param2 := make(map[string]interface{})
				err3 := jsoniter.Unmarshal([]byte(dataString), &param2)
				if err3 == nil {
					param = createParamStrOrder(param2)
				}
			}
			if strings.Index(url, "?") > 0 {
				postUrl = url + "&" + param
			} else {
				postUrl = url + "?" + param
			}
		}
	}
	return postUrl
}

func getMethod(method string) string {
	method = strings.ToUpper(method)

	newMethod := defaultMethod
	methodList := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodDelete,
		http.MethodHead, http.MethodPut, http.MethodPatch}
	for _, one := range methodList {
		if one == method {
			newMethod = method
			break
		}
	}
	return newMethod
}

func getHeaders(headers http.Header, method string, data interface{}) http.Header {
	if headers == nil {
		headers = http.Header{}
	}

	ct := headers.Get("Content-Type")
	if ct == "" {
		isSetType := false
		dataString, err := getDataString(data)
		if err == nil {
			if method == http.MethodGet {
				isSetType = true
				headers.Set("Content-Type", "application/x-www-form-urlencoded")
			}
		}
		if !isSetType {
			if dataString != "" {
				checkJsonError := gojson.CheckValid([]byte(dataString))
				if checkJsonError == nil {
					//表示数据是json格式
					headers.Set("Content-Type", "application/json; charset=utf-8")
				}
			}
		}
	} else {
		//如果data数据是json，并且不是get的话，则不能是 x-www-form-urlencoded
		if method != http.MethodGet {
			if strings.Contains(ct, "x-www-form-urlencoded") {
				dataString, err := getDataString(data)
				if err == nil {
					checkJsonError := gojson.CheckValid([]byte(dataString))
					if checkJsonError == nil {
						//表示数据是json格式
						headers.Set("Content-Type", "application/json; charset=utf-8")
					}
				}
			}
		}
	}

	headers = beautifulHeader(headers)

	return headers
}
func getDataString(data interface{}) (string, error) {
	var paramDataStr string
	typeData := fmt.Sprintf("%T", data)
	if typeData != "string" {
		paramDataByte, err2 := jsonApi.Marshal(data)
		if err2 != nil {
			return "", fmt.Errorf("data 格式目前不支持:%w", err2)
		}
		paramDataStr = string(paramDataByte)
	} else {
		paramDataStr = data.(string)
	}
	return paramDataStr, nil
}

func getRequestId(p *Request) string {
	paramDataOnlyStr := ""
	{
		paramDataStr, err := getDataString(p.Data)
		if err == nil {
			paramDataOnlyStr = utils.GetJsonOnlyKey(paramDataStr)
		}
	}

	headerDataOnlyStr := ""
	{
		if p.Header != nil {
			var keys []string
			for k := range p.Header {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			headerArray := make([]string, 0)
			for _, m := range keys {
				headerArray = append(headerArray, m+"="+p.Header.Get(m))
			}
			headerDataOnlyStr = utils.GetJsonOnlyKey(headerArray)
		}
	}

	return utils.GetUUID(fmt.Sprintf("[%s][%s][%s][%s]", p.Url,
		paramDataOnlyStr, p.Method, headerDataOnlyStr))
}

func getCacheKey(p *Request) string {
	return fmt.Sprintf("{comm-request}%s", getRequestId(p))
}

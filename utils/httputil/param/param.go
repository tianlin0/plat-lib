// Package param 获取参数
package param

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tianlin0/plat-lib/conv"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	LocationHeader Location = "header"
	LocationCookie Location = "cookie"
	LocationQuery  Location = "query"
	LocationBody   Location = "body"
)

type Location string

type paramStruct struct {
	defaultBodyKeyName string
	querySplit         string
	locationOrder      []Location
}

// NewParam 新建
func NewParam() *paramStruct {
	return &paramStruct{
		querySplit:         ",",
		locationOrder:      []Location{LocationBody, LocationQuery},
		defaultBodyKeyName: "request___body___key_name",
	}
}

// SetQuerySplit 设置query传数组以后需要连接起来的字符串分隔符，因为一般不会有多个
func (p *paramStruct) SetQuerySplit(splitStr string) *paramStruct {
	p.querySplit = splitStr
	return p
}

// SetLocationOrder 设置获取参数的位置
func (p *paramStruct) SetLocationOrder(list []Location) *paramStruct {
	p.locationOrder = list
	return p
}

func getMapFromHeaderCookieQuery(l Location, r *http.Request, querySplit string) (map[string]interface{}, interface{}, bool) {
	allRetParamMap := make(map[string]interface{})
	if l == LocationHeader {
		headerMap := getParamFromHeader(r)
		for k, v := range headerMap {
			allRetParamMap[k] = v
		}
		return allRetParamMap, headerMap, true
	}
	if l == LocationCookie {
		headerMap := getParamFromCookie(r)
		for k, v := range headerMap {
			allRetParamMap[k] = v
		}
		return allRetParamMap, headerMap, true
	}
	if l == LocationQuery {
		//如果有多个，则为数组，否则为字符串
		headerMap, _ := getParamFromQuery(r, querySplit)
		for k, v := range headerMap {
			allRetParamMap[k] = v
		}
		return allRetParamMap, headerMap, true
	}

	return nil, nil, false
}
func getMapFromBody(r *http.Request, defaultBodyKeyName string, valueToString bool, querySplit string) (map[string]interface{}, interface{}, error) {
	allRetParamMap, bodyDataStr, err := getMapFromBodyForm(r, defaultBodyKeyName)

	if json.Valid([]byte(bodyDataStr)) {
		allRetParamMap1, paasParamMap, err1 := getMapFromBodyJsonString(bodyDataStr, valueToString)
		for key, one := range allRetParamMap1 {
			allRetParamMap[key] = one
		}
		if err1 == nil {
			return allRetParamMap, paasParamMap, nil
		} else {
			err = err1
		}
	} else {
		allRetParamMap1, err1 := getMapFromBodyQueryString(bodyDataStr, valueToString, querySplit)
		for key, one := range allRetParamMap1 {
			allRetParamMap[key] = one
		}
		if err1 == nil {
			return allRetParamMap, bodyDataStr, nil
		} else {
			err = err1
		}
	}
	return allRetParamMap, bodyDataStr, err
}

func getMapFromBodyJsonString(bodyDataStr string, valueToString bool) (map[string]interface{}, interface{}, error) {
	allRetParamMap := make(map[string]interface{})
	paasParamMap := make(map[string]interface{})
	err := conv.Unmarshal(bodyDataStr, &paasParamMap)
	if err == nil && len(paasParamMap) > 0 {
		//表示是map格式
		for key, one := range paasParamMap {
			allRetParamMap[key] = one
		}
		if valueToString {
			for key, one := range allRetParamMap {
				allRetParamMap[key] = conv.String(one)
			}
		}
		return allRetParamMap, paasParamMap, nil
	}
	//array
	paasParamArray := make([]interface{}, 0)
	err = conv.Unmarshal(bodyDataStr, &paasParamArray)
	if err == nil && len(paasParamArray) > 0 {
		return allRetParamMap, paasParamArray, nil
	}

	if err == nil {
		err = fmt.Errorf("bodyDataStr is not a valid json")
	}

	return allRetParamMap, bodyDataStr, err
}
func getMapFromBodyQueryString(bodyDataStr string, valueToString bool, querySplit string) (map[string]interface{}, error) {
	allRetParamMap := make(map[string]interface{})

	// 如果body中有aaa=bbb&ccc=ddd的格式的话，则直接转换过来
	paramMap, err := getMapByQueryString(bodyDataStr)

	for key, one := range paramMap {
		if valueToString {
			tempVal := ""
			if len(one) == 1 {
				tempVal = one[0]
			} else if len(one) > 1 {
				if querySplit == "" {
					tempVal = conv.String(one)
				} else {
					tempVal = strings.Join(one, querySplit)
				}
			}
			allRetParamMap[key] = tempVal
		} else {
			allRetParamMap[key] = one
		}
	}

	return allRetParamMap, err
}

func getMapFromBodyForm(r *http.Request, defaultBodyKeyName string) (map[string]interface{}, string, error) {
	bodyForm := make(map[string]interface{})

	//先取form的值
	forms := getParamFromForm(r)
	for key, val := range forms {
		bodyForm[key] = val
	}

	var bodyDataStr string
	var err error
	bodyDataStr, err = getParamFromBody(r)
	if bodyDataStr == "" {
		//如果返回为空，则可能是因为header中没有添加 Content-Type: application/json，造成获取不到的情况
		if len(forms) == 1 {
			bodyContent := ""
			for key := range bodyForm {
				if key != "" {
					bodyContent = key
				}
				break
			}
			if bodyContent != "" {
				//为json串才记录，否则与后面的bodyForm重复了
				if json.Valid([]byte(bodyContent)) {
					bodyDataStr = bodyContent
				}
			}
		}
	}

	allRetParamMap := make(map[string]interface{})
	for key, val := range bodyForm {
		if key != "" {
			allRetParamMap[key] = val
		}
	}

	if defaultBodyKeyName != "" {
		allRetParamMap[defaultBodyKeyName] = bodyDataStr
	}

	return allRetParamMap, bodyDataStr, err
}

func getAllByLocation(r *http.Request, l Location, defaultBodyKeyName string, valueToString bool, querySplit string) (map[string]interface{}, interface{}, error) {
	if l == "" {
		return nil, nil, fmt.Errorf("location is null")
	}

	if l == LocationHeader || l == LocationCookie || l == LocationQuery {
		allRetParam, headerMap, found := getMapFromHeaderCookieQuery(l, r, querySplit)
		if found {
			return allRetParam, headerMap, nil
		}
	}

	if l == LocationBody {
		return getMapFromBody(r, defaultBodyKeyName, valueToString, querySplit)
	}

	allRetParamMap := make(map[string]interface{})
	return allRetParamMap, nil, fmt.Errorf("location error: %s", l)
}

// GetAll 转成interface{}
func (p *paramStruct) GetAll(r *http.Request) interface{} {
	allParamMap := make(map[string]interface{})
	for _, one := range p.locationOrder {
		oneParamMap, oneParamRet, err := getAllByLocation(r, one, p.defaultBodyKeyName, false, p.querySplit)
		for key, val := range oneParamMap {
			allParamMap[key] = val
		}
		if err == nil && oneParamRet != nil {
			if oneParamMapTemp, ok := oneParamRet.(map[string]interface{}); ok {
				for key, val := range oneParamMapTemp {
					allParamMap[key] = val
				}
			} else if oneParamList, ok := oneParamRet.([]interface{}); ok {
				return oneParamList
			}
		}
	}

	//为空，则不添加key
	if bodyData, ok := allParamMap[p.defaultBodyKeyName]; ok {
		if conv.String(bodyData) == "" || conv.String(bodyData) == "null" {
			delete(allParamMap, p.defaultBodyKeyName)
		}
	}

	return allParamMap
}

// GetAllMap 转成map hasAll 包含所有的变量，false则去掉默认的data返回值，简化内容
func (p *paramStruct) GetAllMap(r *http.Request, hasAll bool) map[string]interface{} {
	paramMap := p.getAllBasic(r, false)
	if hasAll {
		return paramMap
	}
	//重新复制，不破坏原来的值
	newParamMap := make(map[string]interface{})
	err := conv.Unmarshal(paramMap, &newParamMap)
	if err != nil {
		return paramMap
	}
	if bodyAll, ok := newParamMap[p.defaultBodyKeyName]; ok {
		keyAll := conv.String(bodyAll)
		if valAll, ok := newParamMap[keyAll]; ok {
			if conv.String(valAll) == "[\"\"]" {
				delete(newParamMap, keyAll)
			}
		}
		delete(newParamMap, p.defaultBodyKeyName)
	}
	return newParamMap
}

// GetAllString 转成字符串，主要解决[]string问题
func (p *paramStruct) GetAllString(r *http.Request) map[string]string {
	allParamMap := p.getAllBasic(r, true)
	allParamStr := make(map[string]string)
	for k, v := range allParamMap {
		allParamStr[k] = conv.String(v)
	}
	return allParamStr
}

func (p *paramStruct) getAllBasic(r *http.Request, valueToString bool) map[string]interface{} {
	allParamMap := make(map[string]interface{})
	for _, one := range p.locationOrder {
		oneParamMap, _, err := getAllByLocation(r, one, p.defaultBodyKeyName, valueToString, p.querySplit)
		if err == nil && oneParamMap != nil {
			for key, val := range oneParamMap {
				allParamMap[key] = val
			}
		}
	}
	return allParamMap
}

// GetAllHeaders 获取所有Headers
func (p *paramStruct) GetAllHeaders(r *http.Request) http.Header {
	return getAllHeaders(r)
}

// GetAllCookies 获取所有Cookies
func (p *paramStruct) GetAllCookies(r *http.Request) map[string]*http.Cookie {
	return getAllCookies(r)
}

// GetAllQuery 获取所有Query
func (p *paramStruct) GetAllQuery(r *http.Request) map[string]string {
	allParamMap, _ := getParamFromQuery(r, p.querySplit)
	return allParamMap
}

// GetAllBody 获取所有Body内容
func (p *paramStruct) GetAllBody(r *http.Request) string {
	allParamBody, _ := getParamFromBody(r)
	if allParamBody != "" {
		return allParamBody
	}
	//form里的数据
	forms := getParamFromForm(r)
	if len(forms) == 0 {
		return ""
	}

	retMap := make(map[string]interface{})
	for key, one := range forms {
		retMap[key] = one
	}

	//如果是query的，则直接覆盖，避免有数组的格式，读取数据不正确
	rawQuery := r.URL.RawQuery
	queryMap, err := url.ParseQuery(rawQuery)
	if err == nil {
		for key, one := range queryMap {
			if old, ok := forms[key]; ok {
				if len(one) == 1 && len(old) == 1 && one[0] == old[0] {
					retMap[key] = one[0] //进行替换
				}
			}
		}
	}

	return conv.String(retMap)
}

// GetGinParamFromUrl 取得地址栏的参数
func (p *paramStruct) GetGinParamFromUrl(c *gin.Context, paramName ...string) map[string]string {
	retMap := make(map[string]string)
	if len(paramName) == 0 {
		for _, entry := range c.Params {
			retMap[entry.Key] = entry.Value
		}
		return retMap
	}
	for _, key := range paramName {
		if key == "" {
			continue
		}
		if va, ok := c.Params.Get(key); ok {
			retMap[key] = va
		}
	}
	return retMap
}

func getParamFromHeader(req *http.Request) map[string]string {
	allParamMap := make(map[string]string)
	headers := getAllHeaders(req)
	for i, _ := range headers {
		allParamMap[i] = headers.Get(i)
	}
	return allParamMap
}

func getParamFromCookie(req *http.Request) map[string]string {
	allParamMap := make(map[string]string)
	cookies := getAllCookies(req)
	for i, one := range cookies {
		allParamMap[i] = one.Value
	}
	return allParamMap
}

func getParamFromQuery(req *http.Request, splitString string) (map[string]string, error) {
	allParamMap := make(map[string]string)
	if req == nil {
		return allParamMap, nil
	}

	query, err := getMapByQueryString(req.URL.RawQuery)
	if len(query) > 0 {
		for i, one := range query {
			tempVal := ""
			if len(one) == 0 {
				tempVal = ""
			} else if len(one) == 1 {
				tempVal = one[0]
			} else if len(one) > 1 {
				if splitString == "" {
					tempVal = conv.String(one)
				} else {
					tempVal = strings.Join(one, splitString)
				}
			}
			allParamMap[i] = tempVal
		}
	}

	return allParamMap, err
}

func getParamFromBody(req *http.Request) (string, error) {
	if req == nil || req.Body == nil {
		return "", nil
	}

	b, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(b))

	return string(b), nil
}
func getParamFromForm(req *http.Request) url.Values {
	queryMap := make(url.Values)
	if req == nil {
		return queryMap
	}

	var formValue url.Values

	err := req.ParseForm()
	if err == nil {
		formValue = req.Form
	}
	postValue := req.PostForm
	if formValue == nil {
		formValue = postValue
	} else {
		if postValue != nil {
			for keyName, one := range postValue {
				for _, oneVal := range one {
					oneVal = strings.TrimSpace(oneVal)
					if oneVal == "" {
						continue
					}
					if !formValue.Has(keyName) {
						formValue.Add(keyName, oneVal)
						continue
					}
					isFind := false
					for _, oneVal1 := range formValue[keyName] {
						if oneVal1 == oneVal {
							isFind = true
							break
						}
					}
					//不相同才进行添加，避免相同的情况产生，重复了
					if !isFind {
						formValue.Add(keyName, oneVal)
					}
				}
			}
		}
	}

	if formValue != nil && len(formValue) > 0 {
		return formValue
	}

	return queryMap
}

func getMapByQueryString(query string) (url.Values, error) {
	queryMap := make(url.Values)

	q, err := url.ParseQuery(query)
	if err == nil {
		if q == nil {
			return queryMap, nil
		}

		for k, v := range q {
			newV := make([]string, 0)
			for _, one := range v {
				one = strings.TrimSpace(one)
				if one != "" {
					newV = append(newV, one)
				}
			}
			queryMap[k] = newV
		}
		return queryMap, nil
	}

	p := make(url.Values)
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "&"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if strings.Contains(key, ";") {
			err = fmt.Errorf("invalid semicolon separator in query")
			continue
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		oldValue := value
		value, err1 = url.QueryUnescape(value)
		if err1 != nil {
			value = oldValue
		}
		p[key] = append(p[key], value)
	}
	if p != nil && len(p) > 0 {
		for k, v := range p {
			newV := make([]string, 0)
			for _, one := range v {
				if one != "" {
					newV = append(newV, one)
				}
			}
			queryMap[k] = newV
		}
	}
	return queryMap, err
}

package conv

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/tianlin0/plat-lib/cond"
	"github.com/tianlin0/plat-lib/internal"
	"net/url"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// String 转换为string
func String(src interface{}) string {
	if src == nil {
		return ""
	}

	strType := reflect.TypeOf(src)
	strValue := reflect.ValueOf(src)
	if strType.Kind() == reflect.Ptr {
		if strValue.IsNil() {
			return ""
		}
		return String(strValue.Elem().Interface())
	}

	// 常用特殊类型
	if strValue.Type().String() == "sync.Map" {
		retStr := ""
		if synMap, ok := src.(sync.Map); ok {
			retStr = String(getBySyncMap(&synMap))
		}
		return retStr
	}

	if strType.Kind() == reflect.Map {
		if strValue.IsNil() {
			return ""
		}
		retStr, newMap, err := getByMap(src)
		if err == nil {
			return retStr
		}
		src = newMap
	}

	if strType.Kind() == reflect.Slice {
		if strValue.IsNil() {
			return ""
		}
		retStr, newList, err := getBySlice(src)
		if err == nil {
			return retStr
		}
		src = newList
	}

	retStr, err := getByType(src)
	if err == nil {
		return retStr
	}

	retStr, err = getByTypeString(src)
	if err == nil {
		return retStr
	}

	retStr, err = getByCopy(src) //concurrent map read and map write
	if err == nil {
		return retStr
	}

	fmt.Printf("jsoniter.Marshal error:%s", err.Error())
	return fmt.Sprintf("%v", src)
}

func getBySyncMap(synMap *sync.Map) map[interface{}]interface{} {
	newMap := make(map[interface{}]interface{})
	defer func() {
		if err := recover(); interface{}(err) != nil {
			fmt.Println("getBySyncMap error:", err)
			return
		}
	}()
	fmt.Println("getBySyncMap 1:")
	synMap.Range(func(key, value interface{}) bool {
		fmt.Println("getBySyncMap 2:")
		//newMap[key] = value
		return true
	})
	fmt.Println("getBySyncMap 3:")
	return newMap
}
func getByMap(src interface{}) (string, map[interface{}]interface{}, error) {
	retStr, err := getStringFromJson(src)
	if err == nil {
		return retStr, nil, nil
	}

	strValue := reflect.ValueOf(src)

	newMap := make(map[interface{}]interface{})
	iter := strValue.MapRange()
	for iter.Next() {
		newMap[iter.Key().Interface()] = iter.Value().Interface()
	}

	retStr, err = getStringFromJson(newMap)
	if err == nil {
		return retStr, newMap, nil
	}

	return "", newMap, err
}
func getBySlice(src interface{}) (string, []interface{}, error) {
	//如果是[]byte，则直接转为string
	if strByte, ok := src.([]byte); ok {
		return string(strByte), nil, nil
	}

	json, err := getStringFromJson(src)
	if err == nil {
		return json, nil, nil
	}
	strValue := reflect.ValueOf(src)

	newMap := make([]interface{}, 0)
	for i := 0; i < strValue.Len(); i++ {
		oneItem := strValue.Index(i).Interface()
		newMap = append(newMap, oneItem)
	}

	retStr, err := getStringFromJson(newMap)
	if err == nil {
		return retStr, newMap, nil
	}

	return "", newMap, err
}
func getByType(src interface{}) (string, error) {
	switch src.(type) {
	case []byte:
		return string(src.([]byte)), nil
	case byte:
		return string(src.(byte)), nil
	case string:
		return src.(string), nil
	case int:
		return strconv.Itoa(src.(int)), nil
	case int64:
		return strconv.FormatInt(src.(int64), 10), nil
	case error:
		err, _ := src.(error)
		return err.Error(), nil
	//case float64:
	//	return strconv.FormatFloat(str.(float64), 'g', -1, 64)
	case time.Time:
		{
			oneTime := src.(time.Time)
			//如果为空时间，则返回空字符串
			if cond.IsTimeEmpty(oneTime) {
				//return "", nil
			}
			if cst := GetTimeLocation(); cst != nil {
				return oneTime.In(cst).Format(fullTimeForm), nil
			}
			return oneTime.Format(fullTimeForm), nil
		}
	}
	return "", fmt.Errorf("type error")
}

func getByTypeString(src interface{}) (string, error) {
	strType := fmt.Sprintf("%T", src)
	if strType == "errors.errorString" {
		errTemp := fmt.Sprintf("%v", src)
		if len(errTemp) <= 2 {
			return "", nil
		}
		return errTemp[1 : len(errTemp)-1], nil
	}

	//看看是否是数组类型
	if len(strType) >= 2 {
		subTemp := internal.SubStr(strType, 0, 2)
		if subTemp == "[]" && strType != "[]string" {
			arrTemp := reflect.ValueOf(src)
			newArrTemp := make([]interface{}, 0)
			for i := 0; i < arrTemp.Len(); i++ {
				oneTemp := arrTemp.Index(i).Interface()
				newArrTemp = append(newArrTemp, oneTemp)
			}
			retStr, _, err := getBySlice(newArrTemp)
			return retStr, err
		}
	}

	return "", fmt.Errorf("typeString error")
}
func getByCopy(src interface{}) (string, error) {
	newStrTemp := mapDeepCopy(src) //concurrent map read and map write

	retStr, err := getStringFromJson(newStrTemp)
	if err == nil {
		return retStr, nil
	}
	return "", fmt.Errorf("copy error")
}

func getStringFromJson(src interface{}) (string, error) {
	json, err := jsoniter.MarshalToString(src)
	if err == nil {
		if len(json) >= 2 { //解决返回字符串首位带"的问题
			match, errTemp := regexp.MatchString(`^".*"$`, json)
			if errTemp == nil {
				if match {
					json = json[1 : len(json)-1]
				}
			}
		}
		//解决 & 会转换成 \u0026 的问题
		return strFix(json), nil
	}
	return "", fmt.Errorf("getStringFromJsoniter error:" + err.Error())
}

func strFix(s string) string {
	// https://stackoverflow.com/questions/28595664/how-to-stop-json-marshal-from-escaping-and/28596225
	if strings.Contains(s, "\\u0026") {
		s = strings.Replace(s, "\\u0026", "&", -1)
	}
	if strings.Contains(s, "\\u003c") {
		s = strings.Replace(s, "\\u003c", "<", -1)
	}
	if strings.Contains(s, "\\u003e") {
		s = strings.Replace(s, "\\u003e", ">", -1)
	}
	return s
}

func mapDeepCopy(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for k, v := range v {
			newMap[k] = mapDeepCopy(v)
		}
		return newMap
	case []interface{}:
		newSlice := make([]interface{}, len(v))
		for k, v := range v {
			newSlice[k] = mapDeepCopy(v)
		}
		return newSlice
	default:
		return value
	}
}

// HttpBuildQuery 将map转换为a=1&b=2
func HttpBuildQuery(paramData map[string]interface{}) string {
	params := url.Values{}
	keyList := make([]string, 0)
	for key := range paramData {
		keyList = append(keyList, key)
	}
	sort.Strings(keyList)
	for _, key := range keyList {
		params.Set(key, String(paramData[key]))
	}
	return params.Encode()
}

// ToKeyListFromMap 参数传字符串，避免不是map结构
// {"app":{"mm":1}} ==> "app.mm" : 1
func ToKeyListFromMap(keyMapJsonObject interface{}) map[string]interface{} {
	keyMapJson := String(keyMapJsonObject)

	allMap := make(map[string]interface{})
	if keyMapJson == "" {
		return allMap
	}
	keyMap := make(map[string]interface{})
	keyList := make([]interface{}, 0)
	var err1, err2 error
	err1 = jsoniter.UnmarshalFromString(keyMapJson, &keyMap)
	if err1 != nil {
		err2 = jsoniter.UnmarshalFromString(keyMapJson, &keyList)
		if err2 == nil {
			toStringFromList(keyList, "", nil, 0, allMap)
		}
	} else {
		toStringFromMap(keyMap, nil, 0, allMap)
	}
	if err1 != nil && err2 != nil {
		allMap["."] = keyMapJson
	}
	return allMap
}

func toStringFromList(oneList []interface{}, lastKey string, keyList []string, index int,
	allMap map[string]interface{}) {
	if keyList == nil {
		keyList = make([]string, 0)
	}
	for i, one := range oneList {
		newKey := fmt.Sprintf("%s[%d]", lastKey, i)
		if target2, ok := one.(map[string]interface{}); ok {
			keyList = append(keyList, newKey)
			index = index + 1
			toStringFromMap(target2, keyList, index, allMap)
			index = index - 1
			keyList = append(keyList[:index])
		} else if target3, ok := one.([]interface{}); ok {
			toStringFromList(target3, newKey, keyList, index, allMap)
		} else {
			keyList = append(keyList, newKey)
			keyStr := strings.Join(keyList, ".")
			allMap[keyStr] = one
			keyList = append(keyList[:index])
		}
	}
}

func toStringFromMap(oneMap map[string]interface{}, keyList []string, index int, allMap map[string]interface{}) {
	if keyList == nil {
		keyList = make([]string, 0)
	}
	for key, val := range oneMap {
		if target, ok := val.(map[string]interface{}); ok {
			keyList = append(keyList, key)
			index = index + 1
			toStringFromMap(target, keyList, index, allMap)
			index = index - 1
			keyList = append(keyList[:index])
		} else {
			if list, ok := val.([]interface{}); ok {
				toStringFromList(list, key, keyList, index, allMap)
				continue
			}
			keyList = append(keyList, key)
			keyStr := strings.Join(keyList, ".")
			allMap[keyStr] = val
			keyList = append(keyList[:index])
		}
	}
	return
}

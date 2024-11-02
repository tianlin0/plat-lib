package cond

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// IsNil 判断是否为空
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

// IsTime 是否是时间格式
func IsTime(dateTime string) bool {
	regPattern := "((([0-9]{3}[1-9]|[0-9]{2}[1-9][0-9]{1}|[0-9]{1}[1-9][0-9]{2}|[1-9][0-9]{3})-(((0[13578]|1[02])-"
	regPattern += "(0[1-9]|[12][0-9]|3[01]))|((0[469]|11)-(0[1-9]|[12][0-9]|30))|(02-(0[1-9]|[1][0-9]|2[0-8]))))|"
	regPattern += "((([0-9]{2})(0[48]|[2468][048]|[13579][26])|((0[48]|[2468][048]|[3579][26])00))-02-29))\\s"
	regPattern += "([0-1][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9])$"
	matched, err := regexp.Match(regPattern, []byte(dateTime))
	if err == nil {
		return matched
	}
	return false
}

// IsTimeEmpty 是否为空时间
func IsTimeEmpty(timeParam time.Time) bool {
	nilTime := time.Time{}      //赋零值
	return timeParam == nilTime //此处即为零值
}

// IsNumeric 是否是数字
func IsNumeric(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
	case float32, float64, complex64, complex128:
		return true
	case string:
		str := val.(string)
		if str == "" {
			return false
		}
		// Trim any whitespace
		str = strings.Trim(str, " \\t\\n\\r")
		if str == "" {
			return false
		}
		return toNumberFromString(str)
	}

	//其它类型的全部转换为字符串来判断
	return IsNumeric(fmt.Sprintf("%v", val))
}

func toNumberFromString(str string) bool {
	if str[0] == '-' || str[0] == '+' {
		if len(str) == 1 {
			return false
		}
		str = str[1:]
	}
	// hex
	if len(str) > 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X') {
		for _, h := range str[2:] {
			if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'f') || (h >= 'A' && h <= 'F')) {
				return false
			}
		}
		return true
	}
	// 0-9,Point,Scientific
	p, s, l := 0, 0, len(str)
	for i, v := range str {
		if v == '.' { // Point
			if p > 0 || s > 0 || i+1 == l {
				return false
			}
			p = i
		} else if v == 'e' || v == 'E' { // Scientific
			if i == 0 || s > 0 || i+1 == l {
				return false
			}
			s = i
		} else if v < '0' || v > '9' {
			return false
		}
	}
	return true
}

func Contains[T comparable](s []T, e T) (bool, int) {
	for i, a := range s {
		if a == e {
			return true, i
		}
	}
	return false, -1
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	return false, err
}

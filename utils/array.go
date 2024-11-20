package utils

import (
	"github.com/tianlin0/plat-lib/cond"
	"github.com/tianlin0/plat-lib/conv"
	"reflect"
	"strconv"
	"strings"
)

// Split 字符串，通过多个分隔符进行分隔
func Split(s string, sep []string) []string {
	if s == "" {
		return []string{}
	}
	if len(sep) == 0 {
		return []string{s}
	}
	sepStr := sep[0]
	for i, one := range sep {
		if i == 0 {
			continue
		}
		s = strings.Replace(s, one, sepStr, -1)
	}
	return strings.Split(s, sepStr)
}

//func getKind(kind reflect.Kind) reflect.Kind {
//	switch {
//	case kind >= reflect.Int && kind < reflect.Int64:
//		return reflect.Int
//	case kind >= reflect.Uint && kind < reflect.Uint64:
//		return reflect.Uint
//	case kind >= reflect.Float32 && kind <= reflect.Float64:
//		return reflect.Float64
//	default:
//		return kind
//	}
//}

// Join 任意数组类型，通过分隔符连接成字符串
func Join(i interface{}, s string) string {
	if i == nil {
		return ""
	}
	i = reflect.Indirect(reflect.ValueOf(i)).Interface()
	r := make([]string, 0)
	switch i.(type) {
	case [][]byte:
		temp := i.([][]byte)
		for _, one := range temp {
			r = append(r, string(one))
		}
		return strings.Join(r, s)
	case []string:
		temp := i.([]string)
		return strings.Join(temp, s)
	case []int:
		temp := i.([]int)
		for _, one := range temp {
			r = append(r, strconv.Itoa(one))
		}
		return strings.Join(r, s)
	case []int64:
		temp := i.([]int64)
		for _, one := range temp {
			r = append(r, strconv.FormatInt(one, 10))
		}
		return strings.Join(r, s)
	case []interface{}:
		temp := i.([]interface{})
		for _, one := range temp {
			r = append(r, conv.String(one))
		}
		return strings.Join(r, s)
	case string:
		return conv.String(i)
	}
	return conv.String(i)

}

// AppendUnique 给数组添加不重复的对象
func AppendUnique[T comparable](slice []T, elems ...T) []T {
	allSlice := make([]T, 0)

	if len(slice) > 0 {
		allSlice = append(allSlice, slice...)
	}
	if len(elems) > 0 {
		allSlice = append(allSlice, elems...)
	}

	if len(allSlice) == 0 {
		return allSlice
	}
	newSlice := make([]T, 0)

	for _, elem := range allSlice {
		if ok, _ := cond.Contains(newSlice, elem); ok {
			continue
		}
		newSlice = append(newSlice, elem)
	}

	return newSlice
}

// RemoveItem 移除一个元素
func RemoveItem[T comparable](slice []T, oneElem T) []T {
	allSlice := make([]T, 0)

	if len(slice) > 0 {
		allSlice = append(allSlice, slice...)
	}

	if len(allSlice) == 0 {
		return allSlice
	}
	newSlice := make([]T, 0)

	for _, elem := range allSlice {
		if elem == oneElem {
			continue
		}
		newSlice = append(newSlice, elem)
	}

	return newSlice
}

// SliceDiff find elements that in slice1 but not in slice2
func SliceDiff[T comparable](slice1 []T, slice2 []T) []T {
	diff := make([]T, 0)
	for _, s1 := range slice1 {
		found := false
		for _, s2 := range slice2 {
			if reflect.TypeOf(s1).Name() == reflect.TypeOf(s2).Name() &&
				reflect.ValueOf(s1).Interface() == reflect.ValueOf(s2).Interface() {
				found = true
				break
			}
		}
		if !found {
			diff = append(diff, s1)
		}
	}
	return diff
}

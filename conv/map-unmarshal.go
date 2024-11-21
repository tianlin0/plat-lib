package conv

import (
	"fmt"
	"github.com/jinzhu/copier"
	jsoniter "github.com/tianlin0/plat-lib/internal/jsoniter/go"
	"log"
	"reflect"
)

var Marshal = String

// Unmarshal 将前一个的对象填充到后一个对象中，字段名相同的覆盖值，
// 返回 interface 的作用是如果toPoint为nil的时候，也能正常返回对象.
func Unmarshal(srcStruct interface{}, dstPoint interface{}) error {
	if srcStruct == nil {
		return nil
	}

	isString := false
	oldString := ""
	ok := false
	if oldString, ok = srcStruct.(string); ok {
		isString = true
	}
	if isString {
		//字符串为空，则直接返回
		if oldString == "" {
			return nil
		}
	}

	if dstPoint == nil {
		return fmt.Errorf("unmarshal DstPoint is nil")
	}

	// 1、首先看能否直接赋值
	srcType := reflect.TypeOf(srcStruct)
	dstType := reflect.TypeOf(dstPoint)
	if srcType == dstType {
		if srcType.Kind() != reflect.Ptr &&
			srcType.Kind() != reflect.Struct &&
			srcType.Kind() != reflect.Map {

			err := copier.CopyWithOption(dstPoint, srcStruct, copier.Option{IgnoreEmpty: true, DeepCopy: true})
			if err == nil {
				return nil
			}

			t := new(toolsService)
			err = t.AssignTo(reflect.ValueOf(srcStruct), dstPoint)
			if err == nil {
				return nil
			}

		}
	}

	// 2、不行则用json方法
	t := new(toolsService)
	srcStruct, dstPoint = t.GetNewSrcAndDst(srcStruct, dstPoint)

	//2.2 Unmarshal
	b, err := jsoniter.Marshal(srcStruct)
	if err != nil {
		return err
	}
	errJson := jsoniter.Unmarshal(b, dstPoint)
	if errJson == nil {
		return nil
	}

	// 3、用转换一一覆盖
	//表示有格式不能兼容，出现错误，所以需要进行特殊处理
	srcType = reflect.TypeOf(srcStruct)
	if srcType.Kind() != reflect.Struct &&
		srcType.Kind() != reflect.Map &&
		srcType.Kind() != reflect.Slice {
		//如果是字符串，则需要保证是json格式的
		if isString {
			return fmt.Errorf("Unmarshal error:" + oldString)
		}
		return errJson
	}

	err = AssignTo(srcStruct, dstPoint)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return errJson
	}
	return nil
}

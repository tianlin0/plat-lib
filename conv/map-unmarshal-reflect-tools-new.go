package conv

import (
	"fmt"
	"github.com/tianlin0/plat-lib/cond"
	jsoniter "github.com/tianlin0/plat-lib/internal/jsoniter/go"
	"reflect"
	"strings"
)

type toolsService struct {
}

func (c *toolsService) appendSliceValue(s reflect.Value, x reflect.Value) reflect.Value {
	if !x.IsValid() {
		return s
	}

	if s.Type().Elem() == x.Type() {
		return reflect.Append(s, x)
	}

	if s.Type().Elem().Kind() == reflect.Struct && x.Type().Kind() == reflect.Ptr {
		return c.appendSliceValue(s, x.Elem())
	}

	if s.Type().Elem().Kind() == reflect.Ptr && x.Type().Kind() == reflect.Struct {
		return c.appendSliceValue(s, x.Addr())
	}

	return s
}

// 根据指针获取最终元素的类型
func (c *toolsService) getDirectTypeByPtr(dstType reflect.Type) reflect.Type {
	if dstType.Kind() != reflect.Ptr {
		return dstType
	}
	return c.getDirectTypeByPtr(dstType.Elem())
}

func (c *toolsService) getDirectValueByPtr(dstType reflect.Type, data reflect.Value) reflect.Value {
	if !data.IsValid() {
		return reflect.Value{}
	}

	if dstType.Kind() == reflect.Ptr {
		var newPtr = reflect.New(dstType.Elem())
		if dstType.Elem().Kind() == reflect.Ptr {
			retValue := c.getDirectValueByPtr(dstType.Elem(), data)
			if retValue.IsValid() {
				newPtr.Elem().Set(retValue)
				return newPtr
			}
		} else {
			newPtr.Elem().Set(data)
			return newPtr
		}
	}

	return reflect.Value{}
}

func (c *toolsService) canSetStructColumn(columnName string, columnValue reflect.Value) bool {
	//首字母大写的才需要进行设置
	if columnName == "" {
		return false
	}
	//firstNum := columnName[0:1]
	//// 公共的属性才能设置
	//if firstNum != strings.ToUpper(firstNum) {
	//	return false
	//}

	// 是否能设置
	if !columnValue.CanSet() {
		return false
	}
	return true
}

// 取得所有的column，当type有继承关系时
func (c *toolsService) getAllStructColumn(srcType reflect.Type, srcValue reflect.Value) (
	[]reflect.StructField, []reflect.Value) {

	allStructList := make([]reflect.StructField, 0)
	allValueList := make([]reflect.Value, 0)

	if srcType.Kind() != reflect.Struct {
		return allStructList, allValueList
	}

	for j := 0; j < srcType.NumField(); j++ {
		s := srcType.Field(j)
		v := srcValue.Field(j)
		if s.Name == s.Type.Name() {
			sonTList, sonVList := c.getAllStructColumn(s.Type, v)
			allStructList = append(allStructList, sonTList...)
			allValueList = append(allValueList, sonVList...)
			continue
		}
		allStructList = append(allStructList, s)
		allValueList = append(allValueList, v)
	}

	return allStructList, allValueList
}

func (c *toolsService) assignTo(srcValue reflect.Value, dstPoint interface{}) (err error) {
	if dstPoint == nil {
		return fmt.Errorf("assignTo dstPoint is nil")
	}

	dstType := reflect.TypeOf(dstPoint)
	dstValue := reflect.ValueOf(dstPoint)

	if dstType.Kind() != reflect.Ptr {
		return fmt.Errorf("assignTo dstPoint is not pointer:" + dstType.String())
	}

	g := new(getNewService)
	retData, err := g.GetByDstAll(srcValue.Interface(), dstType)
	if err != nil {
		return err
	}

	if retData.IsValid() {
		if retData.Type().Kind() == reflect.Ptr || retData.Type().Kind() == reflect.Interface {
			if dstValue.Type() == retData.Type() &&
				dstValue.Elem().CanSet() {
				dstValue.Elem().Set(retData.Elem())
				return nil
			}
		}
	}

	return fmt.Errorf("assignTo can not set dstPoint")
}

func (c *toolsService) getTagJsonName(sf reflect.StructField) string {
	tag := sf.Tag.Get("json")
	tag = strings.ReplaceAll(tag, ",omitempty", "")
	if tag == "-" {
		return ""
	}
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag
}

func (c *toolsService) getAllMapNameByField(dstColumn reflect.StructField) []string {
	t := new(toolsService)

	dstColumnTypeName := dstColumn.Name

	dstColumnJsonNameList := make([]string, 0) //后面的覆盖前面的
	//完全一致的情况
	dstColumnJsonNameList = append(dstColumnJsonNameList, dstColumnTypeName)

	dstColumnJsonName := t.getTagJsonName(dstColumn)
	if dstColumnJsonName != "" {
		//如果有设置json，则以json为准
		dstColumnJsonNameList = append(dstColumnJsonNameList, dstColumnJsonName)
	} else {
		// 默认为先snake 后 camel，后者优先
		snakeName := ChangeVariableName(dstColumnTypeName, "snake")
		camelName := ChangeVariableName(dstColumnTypeName, "camel")
		dstColumnJsonNameList = append(dstColumnJsonNameList, snakeName)
		if snakeName != camelName {
			dstColumnJsonNameList = append(dstColumnJsonNameList, camelName)
		}
	}

	return dstColumnJsonNameList
}

func (c *toolsService) split(s string, sep []string) []string {
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

func (c *toolsService) AssignTo(srcValue reflect.Value, dstPoint interface{}) (err error) {
	if dstPoint == nil {
		return fmt.Errorf("continueAssignTo dstPoint is nil")
	}

	dstType := reflect.TypeOf(dstPoint)
	dstValue := reflect.ValueOf(dstPoint)

	if dstType.Kind() != reflect.Ptr {
		return fmt.Errorf("continueAssignTo dstPoint is not pointer:" + dstType.String())
	}

	newData := reflect.New(dstType.Elem())

	defer func() {
		errTemp := recover()
		if !cond.IsNil(errTemp) {
			err = fmt.Errorf(fmt.Sprintf("continueAssignTo error: %v", errTemp))
		}
	}()

	hasSet := false
	if newData.Elem().Type() == srcValue.Type() {
		newData.Elem().Set(srcValue)
		hasSet = true
	} else {
		if srcValue.Kind() == reflect.Interface || srcValue.Kind() == reflect.Ptr {
			if newData.Elem().Type() == srcValue.Elem().Type() {
				newData.Elem().Set(srcValue.Elem())
				hasSet = true
			}
		}
	}
	if !hasSet {
		if srcValue.CanConvert(newData.Elem().Type()) {
			newV := srcValue.Convert(newData.Elem().Type())
			if newData.Elem().Type() == newV.Type() {
				newData.Elem().Set(newV)
				hasSet = true
			}
		}
	}

	if hasSet {
		if dstValue.Elem().Type() == newData.Elem().Type() {
			dstValue.Elem().Set(newData.Elem())
			return nil
		}
	}

	return fmt.Errorf("AssignTo can not set dstPoint")
}

func (c *toolsService) getSrcStruct(srcStruct interface{}) interface{} {
	srcValue := reflect.ValueOf(srcStruct)
	srcType := reflect.TypeOf(srcStruct)

	if srcType.Kind() == reflect.Ptr { //指针
		srcStruct = srcValue.Elem().Interface()
		return c.getSrcStruct(srcStruct) //指针的指针嵌套
	}

	//如果是byte数组类型
	if srcType.Kind() == reflect.Slice {
		if strByte, ok := srcStruct.([]byte); ok {
			srcStruct = string(strByte)
			srcValue = reflect.ValueOf(srcStruct)
			srcType = reflect.TypeOf(srcStruct)
		}
	}

	if srcType.Kind() == reflect.String { //字符串类型
		newStruct := make(map[string]interface{})
		srcStructString := String(srcStruct)
		err := jsoniter.UnmarshalFromString(srcStructString, &newStruct)
		// 数组
		if err != nil {
			newList := make([]interface{}, 0)
			err = jsoniter.UnmarshalFromString(srcStructString, &newList)
			if err == nil {
				srcStruct = newList
			}
		} else {
			srcStruct = newStruct
		}
	}

	return srcStruct
}

func (c *toolsService) getDstPointType(dstPoint interface{}) (newDstPoint interface{}, dstType reflect.Type) {
	if dstPoint == nil {
		dstPoint = new(map[string]interface{})
	}
	dstValue := reflect.ValueOf(dstPoint)
	if dstValue.Kind() != reflect.Ptr {
		dstPoint = new(map[string]interface{})
		dstValue = reflect.ValueOf(dstPoint)
	}
	dstStruct := dstValue.Elem().Interface()
	dstType = reflect.TypeOf(dstStruct)
	return dstPoint, dstType
}
func (c *toolsService) getNewValueByType(dstType reflect.Type) (newDstPoint reflect.Value) {
	if dstType.Kind() == reflect.Ptr {
		return reflect.New(c.getNewValueByType(dstType.Elem()).Type())

	}
	if dstType.Kind() == reflect.Struct {
		return reflect.New(dstType)
	}
	if dstType.Kind() == reflect.Slice {
		sonSliceValue := reflect.MakeSlice(dstType, 0, 0)
		return sonSliceValue
	}

	dstPoint := new(map[string]interface{})

	return reflect.ValueOf(dstPoint)
}

func (c *toolsService) GetNewSrcAndDst(srcStruct interface{}, dstPoint interface{}) (
	newSrcStruct interface{}, newDstPoint interface{}) {
	if srcStruct == nil {
		return nil, dstPoint
	}
	srcStruct = c.getSrcStruct(srcStruct)
	newDstPoint, _ = c.getDstPointType(dstPoint)

	return srcStruct, newDstPoint
}

func (c *toolsService) extendPartDst(srcByte []byte, srcType reflect.Type, dstPoint interface{}) (interface{}, error) {
	if srcType.Kind() == reflect.Slice {
		toPointList := make([]interface{}, 0)
		err2 := jsoniter.Unmarshal(srcByte, &toPointList)
		if err2 == nil {
			srcByte, err2 = jsoniter.Marshal(toPointList)
			if err2 == nil {
				err := jsoniter.Unmarshal(srcByte, dstPoint)
				return dstPoint, err
			}
		}
	} else {
		toPointMap := make(map[string]interface{})
		err2 := jsoniter.Unmarshal(srcByte, &toPointMap)
		if err2 == nil {
			toPointMap2 := make(map[string]string)
			for key, val := range toPointMap {
				toPointMap2[key] = String(val)
			}
			srcByte, err2 = jsoniter.Marshal(toPointMap2)
			if err2 == nil {
				err := jsoniter.Unmarshal(srcByte, dstPoint)
				return dstPoint, err
			}
		}
	}
	return dstPoint, nil
}

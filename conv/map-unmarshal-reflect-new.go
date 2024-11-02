package conv

import (
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/tianlin0/plat-lib/cond"
	jsoniter "github.com/tianlin0/plat-lib/internal/jsoniter/go"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

/*
*
1、目前不能解决继承为小写的情况
2、已有值了，填充没有的情况
*/
func assignTo(srcStruct interface{}, dstPoint interface{}) error {
	//对 srcStruct 和 dstPoint 进行处理
	fill := new(getNewService)
	dstValue, err := fill.GetByDstAll(srcStruct, reflect.TypeOf(dstPoint))
	if err != nil {
		return err
	}
	if !dstValue.IsValid() {
		err = copier.CopyWithOption(dstPoint, srcStruct, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		if err == nil {
			return nil
		}
		return fmt.Errorf("UnmarshalByReflect error: %s, %s", reflect.TypeOf(dstPoint).String(), err.Error())
	}

	//fmt.Println("UnmarshalByReflect getData:", dstValue.Interface())
	t := new(toolsService)
	dstStruct, _ := t.GetNewSrcAndDst(dstValue.Interface(), dstPoint)

	b, err := jsoniter.Marshal(dstStruct)
	if err != nil {
		return err
	}
	errJson := jsoniter.Unmarshal(b, dstPoint)
	if errJson == nil {
		return nil
	}

	errTemp := t.assignTo(dstValue, dstPoint)
	if errTemp != nil {
		err = copier.CopyWithOption(dstPoint, dstStruct, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		if err == nil {
			return nil
		}
		return errTemp
	}

	return nil
}

func NewPtrByType(dstType reflect.Type) interface{} {
	if dstType.Kind() == reflect.Slice {
		newType := reflect.New(dstType)
		dstSliceValue := reflect.MakeSlice(dstType, 0, 0)
		newType.Elem().Set(dstSliceValue)
		return newType.Interface()
	} else if dstType.Kind() == reflect.Ptr {
		ts := new(toolsService)
		dstDirectType := ts.getDirectTypeByPtr(dstType)
		if dstDirectType.Kind() != reflect.Ptr {
			return NewPtrByType(dstDirectType)
		}
		return reflect.New(dstDirectType).Interface()
	} else if dstType.Kind() == reflect.Struct {
		if dstType == reflect.TypeOf(time.Time{}) {
			return &time.Time{}
		}
		return reflect.New(dstType).Interface()
	} else if dstType.Kind() == reflect.Map {
		newType := reflect.New(dstType)
		keyType := dstType.Key()
		valueType := dstType.Elem()
		datMapValue := reflect.MakeMap(reflect.MapOf(keyType, valueType))
		newType.Elem().Set(datMapValue)
		return newType.Interface()
	} else if dstType.Kind() == reflect.String {
		oneStr := ""
		return &oneStr
	} else if dstType.Kind() == reflect.Int64 ||
		dstType.Kind() == reflect.Int {
		oneInt := 0
		return &oneInt
	} else if dstType.Kind() == reflect.Bool {
		oneBool := false
		return &oneBool
	} else if dstType.Kind() == reflect.Interface {
		if dstType == reflect.TypeOf((*error)(nil)).Elem() {
			return fmt.Errorf("")
		}
	}

	return reflect.New(dstType).Interface()
}

type getNewService struct {
}

// GetByDstAll 根据Dst的类型，获取srcInterface的值
func (c *getNewService) GetByDstAll(srcInterface interface{}, dstType reflect.Type) (newDstValue reflect.Value, err error) {
	//fmt.Println("GetByDstAll param:", srcInterface, dstType.String())

	srcType := reflect.TypeOf(srcInterface)
	if srcType == dstType {
		//直接返回
		return reflect.ValueOf(srcInterface), nil
	}

	var newDstList reflect.Value
	var found bool

	if dstType.Kind() == reflect.Slice {
		found = true
		newDstList, err = c.getByDstSlice(srcInterface, dstType)
	} else if dstType.Kind() == reflect.Ptr {
		found = true
		newDstList, err = c.getByDstPtr(srcInterface, dstType)
	} else if dstType.Kind() == reflect.Struct {
		found = true
		newDstList, err = c.getByDstStruct(srcInterface, dstType)
	} else if dstType.Kind() == reflect.Map {
		found = true
		newDstList, err = c.getByDstMap(srcInterface, dstType)
	}

	// 完成
	if found {
		if newDstList.IsValid() && err == nil {
			return newDstList, nil
		}
		newDstList2, err2 := c.getByDstOther(srcInterface, dstType)
		if err2 == nil && newDstList2.IsValid() {
			return newDstList2, nil
		}
		return newDstList, err
	}

	//未找到的情况用默认的方法
	newDstList2, err2 := c.getByDstOther(srcInterface, dstType)
	if err2 == nil && newDstList2.IsValid() {
		return newDstList2, nil
	}

	return newDstList2, err2
}

// getByDstSlice 根据DstSlice获取列表
func (c *getNewService) getByDstSlice(srcSlice interface{}, dstType reflect.Type) (newDstList reflect.Value, err error) {
	if dstType.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("getByDstSlice is not slice:" + dstType.String())
	}

	toPointList := make([]interface{}, 0)

	srcByte, err2 := jsoniter.Marshal(srcSlice)
	if err2 != nil {
		return reflect.Value{}, err2
	}
	err2 = jsoniter.Unmarshal(srcByte, &toPointList)
	if err2 != nil {
		return reflect.Value{}, err2
	}

	dstSliceValue := reflect.MakeSlice(dstType, 0, 0)
	elemType := dstSliceValue.Type().Elem()

	t := new(toolsService)
	for m := 0; m < len(toPointList); m++ {
		oneElem := toPointList[m]
		newDataValue, errTemp := c.GetByDstAll(oneElem, elemType)
		if errTemp == nil && newDataValue.IsValid() {
			dstSliceValue = t.appendSliceValue(dstSliceValue, newDataValue)
		}
		if errTemp != nil {
			err = errTemp
		}
	}

	return dstSliceValue, err
}

// getByDstPtr 根据Ptr获得一个指针对象
func (c *getNewService) getByDstPtr(srcInterface interface{}, dstType reflect.Type) (newDstPtr reflect.Value, err error) {
	if dstType.Kind() != reflect.Ptr {
		return reflect.Value{}, fmt.Errorf("getByDstPtr is not ptr:" + dstType.String())
	}

	t := new(toolsService)

	dstDataType := t.getDirectTypeByPtr(dstType)

	dstDataInterface, err := c.GetByDstAll(srcInterface, dstDataType)
	if err != nil || !dstDataInterface.IsValid() {
		return reflect.Value{}, err
	}

	//fmt.Println("getByDstPtr data:", dstDataInterface.Interface(), dstDataInterface.Type())

	newRetVal := t.getDirectValueByPtr(dstType, dstDataInterface)
	if newRetVal.IsValid() {
		//fmt.Println("getByDstPtr:", newRetVal.Interface())
		//fmt.Println("getByDstPtr:", newRetVal.Type().String())
		return newRetVal, nil
	}
	return reflect.Value{}, nil
}

// getByDstStruct 根据Struct获得一个
func (c *getNewService) getByDstStruct(srcStruct interface{}, dstType reflect.Type) (newDstStruct reflect.Value, err error) {
	if dstType.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("getByDstStruct is not Struct:" + dstType.String())
	}

	//屏蔽意外的错误
	defer func() {
		errTemp := recover()
		if !cond.IsNil(errTemp) {
			log.Println("Unmarshal error:", errTemp)
			err = fmt.Errorf(fmt.Sprintf("getByDstStruct error: %v", errTemp))
		}
	}()

	isSetStruct := false

	t := new(toolsService)

	//查找每一个字段
	dstIns := reflect.New(dstType)

	dstStructValue := dstIns.Elem()
	columnNum := dstType.NumField()
	for i := 0; i < columnNum; i++ {
		dstColumnField := dstType.Field(i)
		dstColumnValue := dstStructValue.Field(i)

		//fmt.Println("dstColumnField index:", columnNum, dstColumnField.Name, i)

		//继承
		if dstColumnField.Name == dstColumnField.Type.Name() {
			//fmt.Println("dstColumnField Type:", dstColumnField.Type.Name())

			newDataValue, errTemp := c.GetByDstAll(srcStruct, dstColumnField.Type)

			//fmt.Println(dstColumnField.Name, dstColumnField.Type.String(), newDataValue.Interface())

			if errTemp == nil && newDataValue.IsValid() {
				//fmt.Println(dstColumnValue.Type().String(), newDataValue.Interface())
				if dstColumnValue.CanSet() {
					dstColumnValue.Set(newDataValue)
					isSetStruct = true
				} else {
					//如果是继承，则需要递归设置
					//if dstColumnValue.Type() == newDataValue.Type() {
					//	newDst := reflect.New(newDataValue.Type())
					//	newDst.Elem().Set(newDataValue)
					//	dstColumnValue.Set(newDst.Elem())
					//}
				}
			}
			continue
		}

		// 当前字段是否能设置，放后面解决类为小写的情况
		if canSet := t.canSetStructColumn(dstColumnField.Name, dstColumnValue); !canSet {
			continue
		}

		//fmt.Println("getByDstStruct: ", dstColumnField.Name, srcStruct)

		//从src获取每一个目标的值,src 是一个整体，需要一一读取
		valueTemp := c.GetSrcFromStructField(srcStruct, dstColumnField)
		if cond.IsNil(valueTemp) {
			//源数据为nil，则不用设置
			continue
		}

		newDataValue, errTemp := c.GetByDstAll(valueTemp, dstColumnField.Type)
		if errTemp == nil && newDataValue.IsValid() {
			//fmt.Println(dstColumnField.Name, dstColumnValue.Type().String(), newDataValue.Interface())
			dstColumnValue.Set(newDataValue)
			isSetStruct = true
		}
	}

	//正常设置
	if isSetStruct {
		//fmt.Println("getByDstStruct: ", dstStructValue.Interface(), dstStructValue.Type().String())
		return dstStructValue, err
	}

	//fmt.Println("getByDstStruct error:", srcStruct)

	return reflect.Value{}, err
}

// getByDstMap 根据map获得一个指针对象
func (c *getNewService) getByDstMap(srcStruct interface{}, dstType reflect.Type) (newDstStruct reflect.Value, err error) {
	if dstType.Kind() != reflect.Map {
		return reflect.Value{}, fmt.Errorf("getByDstMap is not Map:" + dstType.String())
	}
	srcByte, err2 := jsoniter.Marshal(srcStruct)
	if err2 != nil {
		return reflect.Value{}, err2
	}

	keyType := dstType.Key()
	if keyType.Kind() != reflect.String {
		return reflect.Value{}, fmt.Errorf("getByDstMap is not string:" + keyType.String())
	}

	toMap := make(map[string]interface{})
	err2 = jsoniter.Unmarshal(srcByte, &toMap)
	if err2 != nil {
		return reflect.Value{}, err2
	}

	valueType := dstType.Elem()
	datMapValue := reflect.MakeMap(reflect.MapOf(keyType, valueType))

	isSetMap := false
	for key, val := range toMap {
		tempKey, err1 := c.GetByDstAll(key, keyType)
		tempVal, err2 := c.GetByDstAll(val, valueType)
		if err1 == nil &&
			err2 == nil &&
			tempKey.IsValid() &&
			tempVal.IsValid() {
			datMapValue.SetMapIndex(tempKey, tempVal)
			isSetMap = true
		}
	}
	if isSetMap {
		return datMapValue, nil
	}

	return reflect.Value{}, err
}

// getByDstOther 根据map获得一个指针对象
func (c *getNewService) getByDstOther(srcOther interface{}, dstType reflect.Type) (newDstOther reflect.Value, err error) {
	newPtr := reflect.New(dstType)
	srcValue := reflect.ValueOf(srcOther)

	defer func() {
		errTemp := recover()
		if !cond.IsNil(errTemp) {
			err = fmt.Errorf(fmt.Sprintf("getByDstOther error: %v", errTemp))
		}
	}()

	hasSet := false
	if newPtr.Elem().Type() == srcValue.Type() {
		newPtr.Elem().Set(srcValue)
		hasSet = true
	} else {
		if srcValue.Kind() == reflect.Interface || srcValue.Kind() == reflect.Ptr {
			if newPtr.Elem().Type() == srcValue.Elem().Type() {
				newPtr.Elem().Set(srcValue.Elem())
				hasSet = true
			}
		}
	}
	if !hasSet {
		if srcValue.CanConvert(dstType) {
			newV := srcValue.Convert(dstType)
			if newPtr.Elem().Type() == newV.Type() {
				newPtr.Elem().Set(newV)
				hasSet = true
			}
		}
		if !hasSet {
			//自定义转换
			newDst, err := c.getByDstDefault(srcOther, dstType)
			if err == nil {
				if newPtr.Elem().Type() == newDst.Type() {
					newPtr.Elem().Set(newDst)
					hasSet = true
				}
			}
		}
	}
	if hasSet {
		return newPtr.Elem(), err
	}

	return reflect.Value{}, err
}

// getByDstDefault 自定义的格式转换
func (c *getNewService) getByDstDefault(srcDefault interface{}, dstType reflect.Type) (newDstOther reflect.Value, err error) {
	//fmt.Println("getByDstDefault:", srcDefault, dstType.String())

	srcValue := reflect.ValueOf(srcDefault)
	if !srcValue.IsValid() {
		return reflect.Value{}, nil
	}
	if dstType == srcValue.Type() {
		return srcValue, nil
	}

	if retData, ok := c.changeValueToDstByDstType(srcValue, dstType); ok {
		if retData != nil {
			return reflect.ValueOf(retData), nil
		}
		return reflect.Value{}, nil
	}

	if retData, ok := c.changeValueToDstBySrcType(srcValue, dstType); ok {
		if retData != nil {
			return reflect.ValueOf(retData), nil
		}
		return reflect.Value{}, nil
	}

	return reflect.Value{}, err
}

func (c *getNewService) changeValueToDstByDstType(srcValue reflect.Value, dstType reflect.Type) (interface{}, bool) {
	dstTypeName := dstType.Name()
	dstTypeString := dstType.String()

	//fmt.Println("changeValueToDstByDstType:", dstTypeName, dstTypeString, srcValue.Type().String())

	if dstType.Kind() == reflect.String {
		tempTime, ok := c.changeValueToString(srcValue)
		if ok {
			return tempTime, ok
		}
		return nil, true
	}

	if dstType.Kind() == reflect.Int64 {
		tempTime, ok := c.changeValueToInt64(srcValue)
		if ok {
			return tempTime, ok
		}
		return nil, true
	}

	if dstType.Kind() == reflect.Bool {
		tempTime, ok := c.changeValueToBool(srcValue)
		if ok {
			return tempTime, ok
		}
		return nil, true
	}

	if dstType.Kind() == reflect.Interface {
		if dstTypeName == "error" {
			tempTime, ok := c.changeValueToError(srcValue)
			if ok {
				return tempTime, ok
			}
			return nil, true
		}
	}
	if dstType.Kind() == reflect.Struct {
		if dstType == reflect.TypeOf(time.Time{}) {
			tempTime, ok := c.changeValueToTime(srcValue)
			if ok {
				return tempTime, ok
			}
			return nil, true
		}
	}
	if dstType.Kind() == reflect.Slice {
		if dstTypeString == "[]string" {
			tempTime, ok := c.changeValueStringToStringList(srcValue)
			if ok {
				return tempTime, ok
			}
			return nil, true
		}
	}
	return nil, false
}

func (c *getNewService) changeValueToString(srcValue reflect.Value) (string, bool) {
	srcTypeName := srcValue.Type().Name()
	sStr := ""
	if srcTypeName == "int64" {
		sStr = strconv.FormatInt(srcValue.Int(), 10)
		return sStr, true
	} else if srcTypeName == "int" {
		int64Num := srcValue.Int()
		intNum := *(*int)(unsafe.Pointer(&int64Num))
		sStr = strconv.Itoa(intNum)
		return sStr, true
	} else if srcTypeName == "Time" {
		temp := srcValue.Interface().(time.Time)
		return String(temp), true
	} else if strings.Contains(srcTypeName, "byte") {
		temp := srcValue.Interface().(string)
		sStr = String(temp)
		return sStr, true
	} else if srcValue.Interface() != nil {
		sStr = String(srcValue.Interface())
		return sStr, true
	}
	return sStr, false
}
func (c *getNewService) changeValueToTime(srcValue reflect.Value) (time.Time, bool) {
	tempTime, ok := Time(srcValue.Interface())
	if ok {
		return tempTime, true
	}
	return time.Time{}, false
}

func (c *getNewService) changeValueToInt64(srcValue reflect.Value) (int64, bool) {
	srcValueTypeName := srcValue.Type().Name()
	if srcValueTypeName == "int" {
		return srcValue.Int(), true
	}

	srcInterface := srcValue.Interface()
	intTemp, _ := Int64(srcInterface)
	if intTemp != 0 {
		return intTemp, true
	}

	sStr := fmt.Sprintf("%v", srcValue)
	sStrInt, err := strconv.ParseInt(sStr, 10, 64)
	if err == nil {
		return sStrInt, true
	}
	return 0, false
}

func (c *getNewService) changeValueToBool(srcValue reflect.Value) (bool, bool) {
	srcTypeName := reflect.TypeOf(srcValue.Interface()).String()
	if srcTypeName == "string" {
		srcColumnValueString := srcValue.String()
		newColumnString := strings.ToLower(srcColumnValueString)
		if newColumnString == "true" {
			return true, true
		} else if newColumnString == "false" {
			return false, true
		} else {
			sInt, err := strconv.Atoi(srcColumnValueString)
			if err == nil {
				if sInt == 1 {
					return true, true
				} else if sInt == 0 {
					return false, true
				}
			}
		}
	}

	if srcTypeName == "int" ||
		srcTypeName == "float64" ||
		srcTypeName == "int64" {
		srcInterface := srcValue.Interface()
		boolRet, _ := Int64(srcInterface)
		if boolRet == 1 {
			return true, true
		} else if boolRet == 0 {
			return false, true
		}
	}
	return false, false
}

func (c *getNewService) changeValueToError(srcValue reflect.Value) (error, bool) {
	srcIns := srcValue.Interface()
	if srcInsErr, ok := srcIns.(error); ok {
		return srcInsErr, true
	}
	return nil, false
}

func (c *getNewService) changeValueStringToStringList(srcValue reflect.Value) ([]string, bool) {
	srcTypeName := srcValue.Type().Name()
	if srcTypeName == "string" {
		srcColumnValueString := srcValue.String()
		arrList := make([]string, 0)
		err := jsoniter.UnmarshalFromString(srcColumnValueString, &arrList)
		if err == nil {
			return arrList, true
		} else {
			t := new(toolsService)
			arrList = t.split(srcColumnValueString, []string{"|", ";", ","})
			return arrList, true
		}
	}

	return []string{}, false
}

func (c *getNewService) changeFromString(srcValue reflect.Value, dstTypeName string) (interface{}, bool) {
	srcColumnValueString := srcValue.String()
	if dstTypeName == "float32" {
		sFloat, err := strconv.ParseFloat(srcColumnValueString, 32)
		if err == nil {
			return float32(sFloat), true
		}
		return nil, true
	}
	if dstTypeName == "float64" {
		sFloat, err := strconv.ParseFloat(srcColumnValueString, 64)
		if err == nil {
			return sFloat, true
		}
		return nil, true
	}
	if dstTypeName == "int" {
		sInt, err := strconv.Atoi(srcColumnValueString)
		if err == nil {
			return sInt, true
		}
		return nil, true
	}
	if dstTypeName == "Time" {
		sTime, ok := Time(srcColumnValueString)
		if ok {
			return sTime, true
		}
		return nil, true
	}

	return nil, false
}

func (c *getNewService) changeFromByte(srcValue reflect.Value, dstTypeName string) (interface{}, bool) {
	srcInterface := srcValue.Interface()
	if byteTemp, ok := srcInterface.([]byte); ok {
		srcColumnValueString := string(byteTemp)
		if dstTypeName == "string" {
			return srcColumnValueString, true
		} else if dstTypeName == "int64" {
			sStr, _ := Int64(srcColumnValueString)
			if sStr != 0 {
				return sStr, true
			}
		} else if dstTypeName == "int" {
			sInt, err := strconv.Atoi(srcColumnValueString)
			if err == nil {
				return sInt, true
			}
		} else if dstTypeName == "float64" {
			sFloat, err := strconv.ParseFloat(srcColumnValueString, 64)
			if err == nil {
				return sFloat, true
			}
		} else if dstTypeName == "Time" {
			STime, ok := Time(srcColumnValueString)
			if ok {
				return STime, true
			}
		}
	}
	return nil, false
}
func (c *getNewService) changeFromUint8(srcValue reflect.Value, dstTypeName string) (interface{}, bool) {
	srcInterface := srcValue.Interface()
	srcString := String(srcInterface)
	if dstTypeName == "int" {
		one, _ := Int64(srcString)
		return int(one), true
	}
	if dstTypeName == "int64" {
		return Int64(srcString)
	}
	if dstTypeName == "Time" {
		timeTemp, ok := Time(srcString)
		if ok {
			return timeTemp, true
		}
	}

	return srcString, true
}
func (c *getNewService) changeFromFloat64(srcValue reflect.Value, dstTypeName string) (interface{}, bool) {
	float64Num := srcValue.Float()
	if dstTypeName == "int" {
		intNum := int(float64Num)
		return intNum, true
	}
	if dstTypeName == "int64" {
		intNum := int64(float64Num)
		return intNum, true
	}
	if dstTypeName == "string" {
		intStr := strconv.FormatFloat(float64Num, 'E', -1, 64)
		return intStr, true
	}
	return nil, false
}

func (c *getNewService) changeValueToDstBySrcType(srcValue reflect.Value, dstType reflect.Type) (interface{}, bool) {
	dstTypeName := dstType.Name()
	srcTypeString := srcValue.Type().String()

	srcTypeKind := srcValue.Type().Kind()

	//fmt.Println("changeValueToDstBySrcType:")
	//fmt.Println(srcValue.Interface())
	//fmt.Println(srcValue.String())
	//fmt.Println(dstType.String())

	if srcTypeKind == reflect.String {
		if retData, found := c.changeFromString(srcValue, dstTypeName); found {
			return retData, found
		}
	}

	if srcTypeKind == reflect.Int64 {
		if dstType.Kind() == reflect.Int {
			int64Num := srcValue.Int()
			intNum := *(*int)(unsafe.Pointer(&int64Num))
			return intNum, true
		}
	}

	if srcTypeKind == reflect.Float64 {
		if retData, found := c.changeFromFloat64(srcValue, dstTypeName); found {
			return retData, found
		}
	}

	if srcTypeKind == reflect.Slice {
		sonType := srcValue.Type().Elem()
		if sonType.Kind() == reflect.Uint8 {

			if srcTypeString == "[]byte" { //不是普通类型，是复合类型，看看能否再次调用了
				if retData, found := c.changeFromByte(srcValue, dstTypeName); found {
					return retData, found
				}
			}

			if srcTypeString == "[]uint8" {
				if retData, found := c.changeFromUint8(srcValue, dstTypeName); found {
					return retData, found
				}
			}
		}
	}

	if srcTypeString == "map[string]interface {}" {
		srcInterface := srcValue.Interface()
		if _, ok := srcInterface.(map[string]interface{}); ok {
			//for kkk, vvv := range mapTemp {
			//
			//}
			//fmt.Println(mapTemp)
		}

	}

	return nil, false
}

// GetSrcFromStructField 取下一级的数据
func (c *getNewService) GetSrcFromStructField(srcInterface interface{}, dstColumn reflect.StructField) interface{} {
	//1、如果是struct，则首先从struct中进行匹配，名称完全一样的进行匹配
	srcType := reflect.TypeOf(srcInterface)
	if srcType.Kind() == reflect.Struct {
		return c.getColumnValueFromStruct(srcInterface, dstColumn)
	}

	if srcType.Kind() == reflect.Map {
		return c.getColumnValueFromMap(srcInterface, dstColumn)
	}

	return c.getColumnValueFromType(srcInterface, dstColumn)
}

func (c *getNewService) getColumnValueFromStruct(srcStruct interface{}, dstColumn reflect.StructField) interface{} {
	//1、如果是struct，则首先从struct中进行匹配，名称完全一样的进行匹配
	srcType := reflect.TypeOf(srcStruct)
	srcValue := reflect.ValueOf(srcStruct)

	var srcColumnField reflect.StructField
	var srcColumnValue reflect.Value

	findFromSrc := false //名字一模一样，

	t := new(toolsService)
	allSrcTypeList, allSrcValueList := t.getAllStructColumn(srcType, srcValue)

	for j := 0; j < len(allSrcTypeList); j++ {
		s := allSrcTypeList[j]
		//fmt.Println(dstColumn.Name, s.Name)
		if s.Name == dstColumn.Name {
			srcColumnField = s
			srcColumnValue = allSrcValueList[j]
			findFromSrc = true
			break
		}
	}

	//没有找到
	if !findFromSrc {
		return nil
	}
	if !srcColumnValue.IsValid() {
		return nil
	}

	if srcColumnField.Type == dstColumn.Type {
		return srcColumnValue.Interface()
	}

	retVal, err := c.GetByDstAll(srcColumnValue.Interface(), dstColumn.Type)
	if err != nil || !retVal.IsValid() {
		return nil
	}
	if !retVal.IsValid() {
		return nil
	}

	return retVal.Interface()
}

func (c *getNewService) getColumnValueFromMap(srcMap interface{}, dstColumn reflect.StructField) interface{} {
	//1、如果是map，则先用原来的名字，再用json的名字
	srcValue := reflect.ValueOf(srcMap)

	t := new(toolsService)
	dstColumnJsonNameList := t.getAllMapNameByField(dstColumn)

	var srcColumnKey reflect.Value
	findFromSrc := false
	for _, oneName := range dstColumnJsonNameList {
		for _, k := range srcValue.MapKeys() {
			if k.String() == oneName {
				srcColumnKey = k
				findFromSrc = true
				break
			}
		}
	}

	if !findFromSrc {
		return nil
	}

	srcColumnValue := srcValue.MapIndex(srcColumnKey)

	if !srcColumnValue.IsValid() {
		return nil
	}

	if srcColumnValue.Type() == dstColumn.Type {
		return srcColumnValue.Interface()
	}

	retVal, err := c.GetByDstAll(srcColumnValue.Interface(), dstColumn.Type)
	if err != nil || !retVal.IsValid() {
		return nil
	}
	if !retVal.IsValid() {
		return nil
	}

	//fmt.Println("getColumnValueFromMap return:", retVal.Interface())

	return retVal.Interface()
}

// getColumnValueFromType 从里面拿一个值，而不是取本身
func (c *getNewService) getColumnValueFromType(srcInterface interface{}, dstColumn reflect.StructField) interface{} {
	srcByte, err2 := jsoniter.Marshal(srcInterface)
	if err2 != nil {
		return nil
	}

	toMap := make(map[string]interface{})
	err2 = jsoniter.Unmarshal(srcByte, &toMap)
	if err2 != nil {
		return nil
	}

	if len(toMap) == 0 {
		return nil
	}

	return c.getColumnValueFromMap(toMap, dstColumn)
}

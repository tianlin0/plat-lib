// Package conv 转换方法
package conv

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/tianlin0/plat-lib/cond"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Bool convert to bool
func Bool(s interface{}) (bool, bool) {
	if s == nil {
		return false, false
	}
	if b, ok := s.(bool); ok {
		return b, true
	}
	if b, ok := Int64(s); ok {
		if b == 0 {
			return false, true
		}
		return true, true
	}
	bs := strings.ToLower(String(s))
	if bs == "true" {
		return true, true
	}
	if bs == "false" {
		return false, true
	}
	return false, false
}

// Float64 convert any to float64
func Float64(val interface{}) (float64, bool) {
	if val == nil {
		return 0, false
	}
	reValue := reflect.ValueOf(val)
	for reValue.Kind() == reflect.Ptr {
		reValue = reValue.Elem()
		if !reValue.IsValid() {
			return 0, false
		}
		val = reValue.Interface()
		if val == nil {
			return 0, false
		}
		reValue = reflect.ValueOf(val)
	}

	switch v := val.(type) {
	case bool:
		if v {
			return 1, true
		}
		return 0, true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case []byte:
		t, err := strconv.ParseFloat(string(v), 64)
		if err == nil {
			return t, true
		}
		return 0, false
	case json.Number:
		i, err := v.Float64()
		if err != nil {
			return 0, false
		}
		return i, true
	case string:
		if len(v) > 15 {
			return 0, false
		}
		t, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return t, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// Int64 转换为int64，最大：9223372036854775807
func Int64(i interface{}) (int64, bool) {
	if i == nil {
		return 0, false
	}
	intTemp, retBool := forceToInt64(i)
	if retBool {
		return intTemp, true
	}

	switch i.(type) {
	case string:
		intStr := strings.TrimSpace(i.(string))
		num, err := strconv.Atoi(intStr)
		if err != nil {
			//转数字错误，判断是否小数点后面全部是0，则可以去掉 3.0000
			strList := strings.Split(intStr, ".")
			if len(strList) == 2 {
				num1, err := strconv.Atoi(strList[1])
				if err == nil && num1 == 0 {
					num2, err := strconv.Atoi(strList[0])
					if err == nil {
						return int64(num2), true
					}
				}
			}

			return 0, false
		}
		return int64(num), true
	case []byte:
		bits := i.([]byte)
		if len(bits) == 8 {
			return int64(binary.LittleEndian.Uint64(bits)), true
		} else if len(bits) <= 4 {
			num, err := strconv.Atoi(string(bits))
			if err != nil {
				return 0, false
			}
			return int64(num), true
		}
	}
	return Int64(String(i))
}

func forceToInt64(i interface{}) (int64, bool) {
	switch i.(type) {
	case uint:
		return int64(i.(uint)), true
	case uint8:
		return int64(i.(uint8)), true
	case uint16:
		return int64(i.(uint16)), true
	case uint32:
		return int64(i.(uint32)), true
	case uint64:
		return int64(i.(uint64)), true
	case int:
		return int64(i.(int)), true
	case int8:
		return int64(i.(int8)), true
	case int16:
		return int64(i.(int16)), true
	case int32:
		return int64(i.(int32)), true
	case int64:
		return i.(int64), true
	case float32:
		return int64(i.(float32)), true
	case float64:
		return int64(i.(float64)), true
	}
	return 0, false
}

// Time 转换为Time
func Time(val interface{}) (time.Time, bool) {
	timeRet := time.Time{}
	if val == nil {
		return timeRet, true
	}
	reValue := reflect.ValueOf(val)
	for reValue.Kind() == reflect.Ptr {
		reValue = reValue.Elem()
		if !reValue.IsValid() {
			return timeRet, true
		}
		val = reValue.Interface()
		if val == nil {
			return timeRet, true
		}
		reValue = reflect.ValueOf(val)
	}
	if val == "" {
		return timeRet, true
	}

	if v, ok := val.(time.Time); ok {
		return v, true
	}

	valTemp := String(val)
	if timeTemp, ok := toTimeFromString(valTemp); ok {
		return timeTemp, ok
	}

	return timeRet, true
}

func milliTime() int64 {
	return time.Now().UnixMilli()
}

func toTimeFromNormal(v string) (time.Time, error) {
	tLen := len(v)
	if tLen == 0 {
		return time.Time{}, nil
	} else if tLen == 8 {
		return time.ParseInLocation(ShortDateForm08, v, time.Local)
	} else if tLen == len(time.ANSIC) {
		return time.Parse(time.ANSIC, v)
	} else if tLen == len(time.UnixDate) {
		return time.Parse(time.UnixDate, v)
	} else if tLen == len(time.RubyDate) {
		t, err := time.Parse(time.RFC850, v)
		if err != nil {
			t, err = time.Parse(time.RubyDate, v)
		}
		return t, err
	} else if tLen == len(time.RFC822Z) {
		return time.Parse(time.RFC822Z, v)
	} else if tLen == len(time.RFC1123) {
		return time.Parse(time.RFC1123, v)
	} else if tLen == len(time.RFC1123Z) {
		return time.Parse(time.RFC1123Z, v)
	} else if tLen == len(time.RFC3339) {
		return time.Parse(time.RFC3339, v)
	} else if tLen == len(time.RFC3339Nano) {
		return time.Parse(time.RFC3339Nano, v)
	}

	return time.Time{}, fmt.Errorf("no found")
}

func toTimeFromString(v string) (time.Time, bool) {
	t, err := toTimeFromNormal(v)
	if err == nil {
		return t, true
	}

	tLen := len(v)

	if tLen == 10 {
		if cond.IsNumeric(v) {
			mcInt, _ := Int64(v)
			t = time.Unix(mcInt, 0)
			err = nil
			return t, true
		}
		t, err = time.ParseInLocation(FullDateForm, v, time.Local)
	} else if tLen == len(String(milliTime())) { //毫秒
		if cond.IsNumeric(v) {
			mcTempStr := v[0 : len(v)-3]
			mcInt, _ := Int64(mcTempStr)
			t = time.Unix(mcInt, 0)
			err = nil
			return t, true
		}
	} else if tLen == 19 { //毫秒
		t, err = time.ParseInLocation(fullTimeForm, v, time.Local)
		if err != nil {
			t, err = time.Parse(time.RFC822, v)
		}
	} else if tLen == len("2019-12-10T11:18:18.979878") ||
		tLen == len("2019-12-10T11:18:18.9798786") { //毫秒
		tempArr := strings.Split(v, ".")
		if len(tempArr) == 2 {
			timeTemp := tempArr[0]
			timeTemp = strings.Replace(timeTemp, "T", " ", 1)
			t, err = time.ParseInLocation(fullTimeForm, timeTemp, time.Local)
			if err != nil {
				t, err = time.Parse(time.RFC822, v)
			}
		}
	} else {
		if tLen > 19 {
			tempArr := strings.Split(v, ".")
			if len(tempArr) == 2 {
				timeTemp := tempArr[0]
				timeTemp = strings.Replace(timeTemp, "T", " ", 1)
				t, err = time.ParseInLocation(fullTimeForm, timeTemp, time.Local)
				if err == nil {
					return t, true
				}
			}
		}
		t, err = time.Parse(time.RFC1123, v)
	}

	if err != nil {
		{ //2023-04-14T10:09:00Z
			timePattern := "^(\\d{4})-(\\d{2})-(\\d{2})T(\\d{2}):(\\d{2}):(\\d{2})Z$"
			isFind, err := regexp.MatchString(timePattern, v)
			if err == nil {
				if isFind {
					regPattern, _ := regexp.Compile(timePattern)
					patternList := regPattern.FindAllStringSubmatch(v, -1)
					if len(patternList) == 1 {
						if len(patternList[0]) == 7 {
							v1 := fmt.Sprintf("%s-%s-%sT%s:%s:%s+00:00", patternList[0][1],
								patternList[0][2], patternList[0][3],
								patternList[0][4], patternList[0][5], patternList[0][6])
							return toTimeFromString(v1)
						}
					}
					return t, false
				}
			}
		}

		return t, false
	}
	return t, true
}

// UInt32ToIP 将uint32类型转化为ipv4地址
func UInt32ToIP(val uint32) string {
	ipData := net.IPv4(byte(val>>24), byte(val>>16&0xFF), byte(val>>8)&0xFF, byte(val&0xFF))
	return ipData.String()
}

// IPToUInt32 ip转数字
func IPToUInt32(ipAddr string) uint32 {
	bits := strings.Split(ipAddr, ".")
	if len(bits) == 4 {
		b0, _ := strconv.Atoi(bits[0])
		b1, _ := strconv.Atoi(bits[1])
		b2, _ := strconv.Atoi(bits[2])
		b3, _ := strconv.Atoi(bits[3])
		var sum uint32
		sum += uint32(b0) << 24
		sum += uint32(b1) << 16
		sum += uint32(b2) << 8
		sum += uint32(b3)
		return sum
	}
	return 0
}

//
//// int 转大端 []byte
//func IntToBytesBigEndian(n int64, bytesLength byte) ([]byte, error) {
//	switch bytesLength {
//	case 1:
//		tmp := int8(n)
//		bytesBuffer := bytes.NewBuffer([]byte{})
//		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
//		return bytesBuffer.Bytes(), nil
//	case 2:
//		tmp := int16(n)
//		bytesBuffer := bytes.NewBuffer([]byte{})
//		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
//		return bytesBuffer.Bytes(), nil
//	case 3:
//		tmp := int32(n)
//		bytesBuffer := bytes.NewBuffer([]byte{})
//		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
//		return bytesBuffer.Bytes()[1:], nil
//	case 4:
//		tmp := int32(n)
//		bytesBuffer := bytes.NewBuffer([]byte{})
//		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
//		return bytesBuffer.Bytes(), nil
//	case 5:
//		tmp := n
//		bytesBuffer := bytes.NewBuffer([]byte{})
//		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
//		return bytesBuffer.Bytes()[3:], nil
//	case 6:
//		tmp := n
//		bytesBuffer := bytes.NewBuffer([]byte{})
//		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
//		return bytesBuffer.Bytes()[2:], nil
//	case 7:
//		tmp := n
//		bytesBuffer := bytes.NewBuffer([]byte{})
//		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
//		return bytesBuffer.Bytes()[1:], nil
//	case 8:
//		tmp := n
//		bytesBuffer := bytes.NewBuffer([]byte{})
//		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
//		return bytesBuffer.Bytes(), nil
//	}
//	return nil, fmt.Errorf("IntToBytesBigEndian b param is inValid")
//}
//
////int 转小端 []byte
//func IntToBytesLittleEndian(n int64, bytesLength byte) ([]byte, error) {
//	switch bytesLength {
//	case 1:
//		tmp := int8(n)
//		bytesBuffer := bytes.NewBuffer([]byte{})
//		binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
//		return bytesBuffer.Bytes(), nil
//	case 2:
//		tmp := int16(n)
//		bytesBuffer := bytes.NewBuffer([]byte{})
//		binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
//		return bytesBuffer.Bytes(), nil
//	case 3:
//		tmp := int32(n)
//		bytesBuffer := bytes.NewBuffer([]byte{})
//		binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
//		return bytesBuffer.Bytes()[0:3], nil
//	case 4:
//		tmp := int32(n)
//		bytesBuffer := bytes.NewBuffer([]byte{})
//		binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
//		return bytesBuffer.Bytes(), nil
//	case 5:
//		tmp := n
//		bytesBuffer := bytes.NewBuffer([]byte{})
//		binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
//		return bytesBuffer.Bytes()[0:5], nil
//	case 6:
//		tmp := n
//		bytesBuffer := bytes.NewBuffer([]byte{})
//		binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
//		return bytesBuffer.Bytes()[0:6], nil
//	case 7:
//		tmp := n
//		bytesBuffer := bytes.NewBuffer([]byte{})
//		binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
//		return bytesBuffer.Bytes()[0:7], nil
//	case 8:
//		tmp := n
//		bytesBuffer := bytes.NewBuffer([]byte{})
//		binary.Write(bytesBuffer, binary.LittleEndian, &tmp)
//		return bytesBuffer.Bytes(), nil
//	}
//	return nil, fmt.Errorf("IntToBytesLittleEndian b param is inValid")
//}
//
////(大端) []byte 转 uint
//func BytesToUIntBigEndian(b []byte) (int, error) {
//	if len(b) == 3 {
//		b = append([]byte{0}, b...)
//	}
//	bytesBuffer := bytes.NewBuffer(b)
//	switch len(b) {
//	case 1:
//		var tmp uint8
//		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
//		return int(tmp), err
//	case 2:
//		var tmp uint16
//		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
//		return int(tmp), err
//	case 4:
//		var tmp uint32
//		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
//		return int(tmp), err
//	default:
//		return 0, fmt.Errorf("%s", "BytesToInt bytes length is inValid!")
//	}
//}
//
////(大端) []byte 转 int
//func BytesToIntBigEndian(b []byte) (int, error) {
//	if len(b) == 3 {
//		b = append([]byte{0}, b...)
//	}
//	bytesBuffer := bytes.NewBuffer(b)
//	switch len(b) {
//	case 1:
//		var tmp int8
//		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
//		return int(tmp), err
//	case 2:
//		var tmp int16
//		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
//		return int(tmp), err
//	case 4:
//		var tmp int32
//		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
//		return int(tmp), err
//	default:
//		return 0, fmt.Errorf("%s", "BytesToInt bytes length is inValid!")
//	}
//}
//
////(小端) []byte 转 uint
//func BytesToUIntLittleEndian(b []byte) (int, error) {
//	if len(b) == 3 {
//		b = append([]byte{0}, b...)
//	}
//	bytesBuffer := bytes.NewBuffer(b)
//	switch len(b) {
//	case 1:
//		var tmp uint8
//		err := binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
//		return int(tmp), err
//	case 2:
//		var tmp uint16
//		err := binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
//		return int(tmp), err
//	case 4:
//		var tmp uint32
//		err := binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
//		return int(tmp), err
//	default:
//		return 0, fmt.Errorf("%s", "BytesToInt bytes length is inValid!")
//	}
//}
//
//// 小端[]byte 转 int
//func BytesToIntLittleEndian(b []byte) (int, error) {
//	if len(b) == 3 {
//		b = append([]byte{0}, b...)
//	}
//	bytesBuffer := bytes.NewBuffer(b)
//	switch len(b) {
//	case 1:
//		var tmp int8
//		err := binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
//		return int(tmp), err
//	case 2:
//		var tmp int16
//		err := binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
//		return int(tmp), err
//	case 4:
//		var tmp int32
//		err := binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
//		return int(tmp), err
//	default:
//		return 0, fmt.Errorf("%s", "BytesToInt bytes length is inValid!")
//	}
//}

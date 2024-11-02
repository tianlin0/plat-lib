package conv

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"time"
	"unsafe"
)

const (
	fullTimeForm     = "2006-01-02 15:04:05"
	FullDateForm     = "2006-01-02"
	ShortTimeForm10  = "0102150405"
	ShortTimeForm12  = "060102150405"
	ShortTimeForm14  = "20060102150405"
	ShortDateForm08  = "20060102"
	ShortMonthForm06 = "200601"
)

var (
	sysTimeLocation = "Asia/Chongqing"
)

// SetTimeLocation 设置时区
func SetTimeLocation(location string) {
	sysTimeLocation = location
}

// GetTimeLocation 获得时区
func GetTimeLocation() *time.Location {
	if cst, err := time.LoadLocation(sysTimeLocation); err == nil {
		return cst
	}
	return nil
}

// timeCodec 时间格式转换
type timeCodec struct {
}

// Decode 转码
func (codec *timeCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	var ok bool
	s := iter.ReadString()
	*((*time.Time)(ptr)), ok = Time(s)
	if !ok {
		iter.ReportError("decode time.Time", fmt.Sprint(s, " is not valid time format"))
	}
}

// IsEmpty 是否为空时间
func (codec *timeCodec) IsEmpty(ptr unsafe.Pointer) bool {
	ts := *((*time.Time)(ptr))
	return ts.UnixNano() == 0
}

// Encode 转码
func (codec *timeCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	ts := *((*time.Time)(ptr))
	stream.WriteString(ts.Format("2006-01-02 15:04:05"))
}

func init() {
	jsoniter.RegisterTypeEncoder("time.Time", &timeCodec{})
	jsoniter.RegisterTypeDecoder("time.Time", &timeCodec{})

	//php兼容模式[]，{}
	extra.RegisterFuzzyDecoders()
}

// FormatFromUnixTime 将unix时间戳格式化为YYYYMMDD HH:MM:SS格式字符串
// FormatFromUnixTime(FullDate, 12321312)
func FormatFromUnixTime(formatStr ...interface{}) string {
	format := fullTimeForm
	var timeNum int64 = 0

	if len(formatStr) == 1 {
		if times, ok := formatStr[0].(int64); ok {
			timeNum = times
		} else if strTemp, ok := formatStr[0].(string); ok {
			if strTemp != "" {
				format = strTemp
			}
		}
	} else if len(formatStr) == 2 {
		if times, ok := formatStr[0].(int64); ok {
			timeNum = times
		} else if strTemp, ok := formatStr[0].(string); ok {
			if strTemp != "" {
				format = strTemp
			}
		}
		if times, ok := formatStr[1].(int64); ok {
			timeNum = times
		} else if strTemp, ok := formatStr[1].(string); ok {
			if strTemp != "" {
				format = strTemp
			}
		}
	}

	if timeNum > 0 {
		return time.Unix(timeNum, 0).Format(format)
	}
	return time.Now().Format(format)
}

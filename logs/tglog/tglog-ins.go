package tglog

import (
	"fmt"
	"github.com/tianlin0/plat-lib/conv"
	"github.com/tianlin0/plat-lib/utils"
	"time"
)

// ILogParam 将参数转换为数组
type ILogParam interface {
	ToArray() (tableName string, paramList []interface{}) //通过定义的对象，返回参数的数组
	GetLogString() string                                 //获取日志关键信息
	GetLogId() string                                     //获取日志唯一ID
}

// LogService tgLog的输出格式
type LogService struct {
	PlatName    string
	LogEnv      string
	CurrentTime time.Time
	LogId       string
	ServerIp    string
	ServerPort  string
	ClientIp    string
	ClientPort  string
	LogType     string
	LogLevel    string
	Prefix      string
	Method      string
	Tag         string
	LogMessage  interface{}
	ExpendTime  int64
	Sort        int
	LoginUser   string
}

// GetLogString logService的日志信息信息
func (l *LogService) GetLogString() string {
	if conv.String(l.LogMessage) == "" {
		return "" //为空则跳过
	}
	return fmt.Sprintf("%s: %s", l.LogId, conv.String(l.LogMessage))
}

// GetLogId logService的日志唯一ID
func (l *LogService) GetLogId() string {
	return l.LogId
}

// ToArray logService的格式
func (l *LogService) ToArray() (tableName string, paramList []interface{}) {
	arrayList := make([]interface{}, 0)
	arrayList = append(arrayList, l.PlatName)
	arrayList = append(arrayList, l.LogEnv)
	arrayList = append(arrayList, l.CurrentTime)
	arrayList = append(arrayList, l.LogId)
	arrayList = append(arrayList, l.ServerIp)
	arrayList = append(arrayList, l.ServerPort)
	arrayList = append(arrayList, l.ClientIp)
	arrayList = append(arrayList, l.ClientPort)
	arrayList = append(arrayList, l.LogType)
	arrayList = append(arrayList, l.LogLevel)
	arrayList = append(arrayList, l.Prefix)
	arrayList = append(arrayList, l.Method)
	arrayList = append(arrayList, l.Tag)
	arrayList = append(arrayList, l.LogMessage)
	arrayList = append(arrayList, l.ExpendTime)
	arrayList = append(arrayList, l.Sort)
	arrayList = append(arrayList, l.LoginUser)
	return "LogService", arrayList
}

// EventAlarm 告警格式
type EventAlarm struct {
	ProjectName string //paas_name
	GroupVer    string //AppName
	Env         string
	Uid         string //LoginUser
	EventType   string
	EventName   string
	EventMsg    interface{}
	Uri         string
	TimeOver    int64 //ExpendTime
	Ext1        string
	Ext2        string
	Ext3        string
}

// ToArray 告警的格式输出
func (l *EventAlarm) ToArray() (tableName string, paramList []interface{}) {
	arrayList := make([]interface{}, 0)
	arrayList = append(arrayList, l.ProjectName)
	arrayList = append(arrayList, l.GroupVer)
	arrayList = append(arrayList, l.Env)
	arrayList = append(arrayList, l.Uid)
	arrayList = append(arrayList, l.EventType)
	arrayList = append(arrayList, l.EventName)
	arrayList = append(arrayList, l.EventMsg)
	arrayList = append(arrayList, l.Uri)
	arrayList = append(arrayList, l.TimeOver)
	arrayList = append(arrayList, l.Ext1)
	arrayList = append(arrayList, l.Ext2)
	arrayList = append(arrayList, l.Ext3)
	return "EventAlarm", arrayList
}

// GetLogString 获取Event日志唯一信息
func (l *EventAlarm) GetLogString() string {
	return fmt.Sprintf("%s,%s,%s", l.ProjectName, l.Uid, l.EventMsg)
}

// GetLogId 获取日志唯一ID
func (l *EventAlarm) GetLogId() string {
	return utils.GetUUID(l.GetLogString())
}

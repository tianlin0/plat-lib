package logs

import (
	"github.com/tianlin0/plat-lib/cond"
	"time"
)

// LogCommData 不会变的数据
type LogCommData struct {
	createTime time.Time              `json:"createTime"`
	LogId      string                 `json:"id"`               //logId
	UserId     string                 `json:"userid,omitempty"` //userID
	Env        string                 `json:"env"`              //env
	Path       string                 `json:"path,omitempty"`   //当前请求的地址
	Method     string                 `json:"method,omitempty"` //当前请求的方法
	Extend     map[string]interface{} `json:"extend,omitempty"` //额外的参数
}

// LogData 每条日志的数据
type LogData struct {
	LogCommData
	Now      time.Time     `json:"now"` //初始化时间
	FileName string        //文件名
	Line     int           //行号
	LogLevel LogLevel      `json:"logLevel"`
	Message  []interface{} `json:"message"`
}

// NewLogData 初始化一个日志变量
func NewLogData(logCommData ...*LogCommData) *LogData {
	l := new(LogData)
	if logCommData != nil && len(logCommData) > 0 {
		if logCommData[0] != nil {
			l.LogCommData = *(logCommData[0])
		}
	}

	commData := &(l.LogCommData)
	commData.Init()
	return l
}

// Init 必须初始化的
func (l *LogCommData) Init() *LogCommData {
	if cond.IsTimeEmpty(l.createTime) {
		l.createTime = time.Now()
	}
	return l
}

func (l *LogData) String() string {
	return defaultConfig.CommLogString(l)
}

// AddLogMessage 将日志添加到列表中
func (l *LogData) AddLogMessage(level LogLevel, fileName string, line int, msg ...interface{}) {
	if len(msg) == 0 {
		return
	}
	l.Now = time.Now()
	l.FileName = fileName
	l.Line = line
	l.LogLevel = level
	l.Message = append([]interface{}{}, msg...)
}

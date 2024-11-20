package logs

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/tianlin0/plat-lib/utils"
	"log"
	"time"
)

// printLogger 自定义日志的使用方式
type printLogger struct {
	logExecute  LogExecute
	logCommData *LogCommData
	logLevel    LogLevel
	callerSkip  int
}

// NewPrintLogger 例子，一个完整的日志，需要实现如下方法
func NewPrintLogger(level LogLevel, logCommData ...*LogCommData) *printLogger {
	printLoggerInstance := &printLogger{
		logCommData: &LogCommData{},
		callerSkip:  2,
		logExecute: func(ctx context.Context, logInfo *LogData) {
			defaultPrintLogExecute(logInfo)
		},
	}
	if logCommData != nil && len(logCommData) > 0 {
		printLoggerInstance.logCommData = logCommData[0].Init()
	}

	if level >= DEBUG {
		printLoggerInstance.SetLevel(level)
	}

	return printLoggerInstance
}

func (x *printLogger) SetCallerSkip(skip int) {
	x.callerSkip = skip
}

func defaultPrintLogExecute(logInfo *LogData) {
	if logInfo == nil || logInfo.Message == nil || len(logInfo.Message) == 0 {
		return
	}

	msgStr := logInfo.String()
	if msgStr == "" {
		return
	}

	if logInfo.LogLevel < ERROR {
		log.Println(msgStr)
		return
	}

	//如果是错误，则直接红色打印
	nowTime := time.Now().Format("2006/01/02 15:04:05")
	slantedRed := color.New(color.FgRed, color.Bold)

	newList := make([]interface{}, 0)
	newList = append(newList, fmt.Sprintf("%s %s", nowTime, msgStr))

	_, err := slantedRed.Println(newList...)
	if err == nil {
		return
	}
}

func (x *printLogger) printlnComm(level LogLevel, msg ...interface{}) {
	if len(msg) == 0 {
		return
	}

	logNewInfo := NewLogData(x.logCommData)

	fileName := ""
	line := 0
	file, _ := utils.SpecifyContext(x.callerSkip)
	if file != nil {
		fileName = file.FileName
		line = file.Line
	}
	logNewInfo.AddLogMessage(level, fileName, line, msg...)

	x.logExecute(context.Background(), logNewInfo)
}

// Debug
func (x *printLogger) Debug(v ...interface{}) {
	if x.Level() > DEBUG {
		return
	}
	x.printlnComm(DEBUG, v...)
}

// Error
func (x *printLogger) Error(v ...interface{}) {
	if x.Level() > ERROR {
		return
	}
	x.printlnComm(ERROR, v...)
}

// Info
func (x *printLogger) Info(v ...interface{}) {
	if x.Level() > INFO {
		return
	}
	x.printlnComm(INFO, v...)
}

// Warn
func (x *printLogger) Warn(v ...interface{}) {
	if x.Level() > WARNING {
		return
	}
	x.printlnComm(WARNING, v...)
}

// Level
func (x *printLogger) Level() LogLevel { return x.logLevel }

// SetLevel
func (x *printLogger) SetLevel(l LogLevel) {
	x.logLevel = l
}

// LogId
func (x *printLogger) LogId(ctx context.Context) string {
	return x.logCommData.LogId
}

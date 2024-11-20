// Package filelog 文件日志
package filelog

import (
	"context"
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/tianlin0/plat-lib/logs"
	"path/filepath"
	"runtime"
	"time"
)

// FileLogger 本地文件日志
type FileLogger struct {
	LinkName      string //整个文件路径+文件名
	FileEndName   string //分隔后的文件名
	RotationTime  time.Duration
	RotationCount uint
}

// Debug 调试
func (x *FileLogger) Debug(v ...interface{}) {
	if x.Level() < logs.DEBUG {
		return
	}
	logrus.Debug(v...)
}

//Error 错误
func (x *FileLogger) Error(v ...interface{}) {
	if x.Level() < logs.ERROR {
		return
	}
	logrus.Error(v...)
}

//Info 普通
func (x *FileLogger) Info(v ...interface{}) {
	if x.Level() < logs.INFO {
		return
	}
	logrus.Info(v...)
}

//Warn 警告
func (x *FileLogger) Warn(v ...interface{}) {
	if x.Level() < logs.WARNING {
		return
	}
	logrus.Warn(v...)
}

//Level 级别
func (x *FileLogger) Level() logs.LogLevel {
	l := logrus.GetLevel()
	return getLogLevel(l)
}

func (x *FileLogger) LogId(ctx context.Context) string {
	return ""
}

//SetLevel 设置
func (x *FileLogger) SetLevel(l logs.LogLevel) {
	fLevel := getLogrusLevel(l)
	logrus.SetLevel(fLevel)
}

// NewFileLogger 新建文件日志
func NewFileLogger(fl *FileLogger, logLevel logs.LogLevel) logs.ILogger {
	if fl == nil {
		fl = new(FileLogger)
	}

	//lineHook, err := newLineHook(fl.FileLineFieldName)
	//if err == nil {
	//	logrus.AddHook(lineHook)
	//}

	lfsHook, err := newLfsHook(fl.LinkName, fl.FileEndName, fl.RotationTime, fl.RotationCount)
	if err == nil {
		logrus.AddHook(lfsHook)
	}

	logrus.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FieldsOrder:     []string{"component", "category"},
		HideKeys:        true,
		CallerFirst:     true,
		NoColors:        false,
		CustomCallerFormatter: func(frame *runtime.Frame) string {
			return fmt.Sprintf(" [%s](%s): %d", frame.Function, filepath.Base(frame.File), frame.Line)
		},
	})
	fl.SetLevel(logLevel)
	logrus.SetReportCaller(true)

	var logger logs.ILogger = fl
	return logger
}

var levelMap = map[logs.LogLevel]logrus.Level{
	logs.DEBUG:     logrus.DebugLevel,
	logs.INFO:      logrus.InfoLevel,
	logs.NOTICE:    logrus.TraceLevel,
	logs.WARNING:   logrus.WarnLevel,
	logs.ERROR:     logrus.ErrorLevel,
	logs.ALERT:     logrus.FatalLevel,
	logs.EMERGENCY: logrus.PanicLevel,
}

func getLogrusLevel(l logs.LogLevel) logrus.Level {
	for key, val := range levelMap {
		if key == l {
			return val
		}
	}
	return logrus.WarnLevel
}

func getLogLevel(l logrus.Level) logs.LogLevel {
	for key, val := range levelMap {
		if val == l {
			return key
		}
	}
	return logs.DEBUG
}

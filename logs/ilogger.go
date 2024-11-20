package logs

import (
	"context"
)

// ILogger is a logger interface
type ILogger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})

	Level() LogLevel
	SetLevel(l LogLevel)

	LogId(ctx context.Context) string
	//WithFields(fields Fields)
	//Field(key string) interface{}
}

type Fields map[string]interface{}

// LogExecute log处理方法
type LogExecute func(ctx context.Context, logInfo *LogData) //日志的处理函数
type LogString func(logInfo *LogData) string

type CallSkip interface {
	SetCallerSkip(skip int)
}

// Logger 直接根据等级打印所有日志
func Logger(logger ILogger, l LogLevel, msg ...interface{}) {
	if l <= DEBUG {
		logger.Debug(msg...)
	} else if l <= INFO {
		logger.Info(msg...)
	} else if l <= WARNING {
		logger.Warn(msg...)
	} else if l <= ERROR {
		logger.Error(msg...)
	} else {
		logger.Debug(msg...)
	}
}

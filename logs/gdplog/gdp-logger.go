package gdplog

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/tianlin0/plat-lib/httputil"
	"github.com/tianlin0/plat-lib/logs"
	"github.com/tianlin0/plat-lib/utils"
)

// gdpLogger 自定义日志的使用方式
type gdpLogger struct {
	logLevel   logs.LogLevel
	gdpLogger  logrus.FieldLogger
	callerSkip int
}

// NewGdpLogger 例子，一个完整的日志，需要实现如下方法
func NewGdpLogger(ctx context.Context, level logs.LogLevel, logger logrus.FieldLogger) (*gdpLogger, context.Context) {
	logs.SetConfig(&logs.Config{
		DefaultLogger: &gdpLogger{
			gdpLogger:  DefaultLogger(),
			logLevel:   level,
			callerSkip: 3,
		},
	})

	gdpLoggerInstance := new(gdpLogger)
	gdpLoggerInstance.callerSkip = 2
	newLogger := false
	if logger != nil {
		newLogger = true
		gdpLoggerInstance.gdpLogger = logger
	} else {
		gdpLoggerInstance.gdpLogger = Logger(ctx)
	}

	gdpLoggerInstance.SetLevel(level)

	if newLogger {
		newCtx := SetLoggerToContext(ctx, gdpLoggerInstance.gdpLogger)
		return gdpLoggerInstance, newCtx
	}

	return gdpLoggerInstance, ctx
}

// GetGdpLogger 获取公共日志
func GetGdpLogger(ctx context.Context) logs.ILogger {
	gdpLoggerInstance := &gdpLogger{
		gdpLogger: Logger(ctx),
	}
	return gdpLoggerInstance
}

func (x *gdpLogger) SetCallerSkip(skip int) {
	x.callerSkip = skip
}

// ctxPrintlnComm 新增日志到context的列表中
func (x *gdpLogger) ctxPrintlnComm(level logs.LogLevel, msg ...interface{}) {
	if msg == nil || len(msg) == 0 {
		return
	}

	oldStr := utils.Join(msg, " ")

	//logContext, _ := utils.SpecifyContext(x.callerSkip)
	fields := logrus.Fields{}
	//if logContext != nil {
	//	fields["file"] = fmt.Sprintf("%s:%d", logContext.FileName, logContext.Line)
	//}

	logger, ok := x.gdpLogger.(*logrus.Logger)
	if ok {
		logger.SetReportCaller(false)

		if level == logs.DEBUG {
			logger.WithFields(fields).Debug(oldStr)
		} else if level == logs.ERROR {
			logger.WithFields(fields).Error(oldStr)
		} else if level == logs.INFO {
			logger.WithFields(fields).Info(oldStr)
		} else if level == logs.WARNING {
			logger.WithFields(fields).Warn(oldStr)
		}
		return
	}

	logger2, ok := x.gdpLogger.(*logrus.Entry)
	if ok {
		logger2.Logger.SetReportCaller(false)

		if level == logs.DEBUG {
			logger2.WithFields(fields).Debug(oldStr)
		} else if level == logs.ERROR {
			logger2.WithFields(fields).Error(oldStr)
		} else if level == logs.INFO {
			logger2.WithFields(fields).Info(oldStr)
		} else if level == logs.WARNING {
			logger2.WithFields(fields).Warn(oldStr)
		}
		return
	}

	if level == logs.DEBUG {
		x.gdpLogger.Debug(oldStr)
	} else if level == logs.ERROR {
		x.gdpLogger.Error(oldStr)
	} else if level == logs.INFO {
		x.gdpLogger.Info(oldStr)
	} else if level == logs.WARNING {
		x.gdpLogger.Warn(oldStr)
	}
}

// Debug 调试
func (x *gdpLogger) Debug(v ...interface{}) {
	if x.logLevel > logs.DEBUG {
		return
	}
	x.ctxPrintlnComm(logs.DEBUG, v...)
}

// error
func (x *gdpLogger) Error(v ...interface{}) {
	if x.logLevel > logs.ERROR {
		return
	}
	x.ctxPrintlnComm(logs.ERROR, v...)
}

// Info 信息
func (x *gdpLogger) Info(v ...interface{}) {
	if x.logLevel > logs.INFO {
		return
	}
	x.ctxPrintlnComm(logs.INFO, v...)
}

// Warn 警告
func (x *gdpLogger) Warn(v ...interface{}) {
	if x.logLevel > logs.WARNING {
		return
	}
	x.ctxPrintlnComm(logs.WARNING, v...)
}

// Level 等级
func (x *gdpLogger) Level() logs.LogLevel { return x.logLevel }

// SetLevel SetLevel
func (x *gdpLogger) SetLevel(l logs.LogLevel) {
	x.logLevel = l

	var logTemp logrus.Level

	switch l {
	case logs.ERROR:
		logTemp = logrus.ErrorLevel
	case logs.WARNING:
		logTemp = logrus.WarnLevel
	case logs.INFO:
		logTemp = logrus.InfoLevel
	case logs.DEBUG:
		logTemp = logrus.DebugLevel
	}
	SetDefaultLoggerLevel(logTemp)
}

func (x *gdpLogger) LogId(ctx context.Context) string {
	logId := ctx.Value(ContextKeyRequestID)
	if logId != nil {
		if logIdStr, ok := logId.(string); ok {
			return logIdStr
		}
	}
	return httputil.GetLogId()
	//x.gdpLogger.
	//
	//requestID := c.GetHeader("X-Request-Id")
	//traceID := c.GetHeader("traceparent")
	//if requestID == "" {
	//	// 请求头X-Request-Id不存在，生成一个新的uuid
	//	requestUUID := uuid.New()
	//	requestID = fmt.Sprintf("g-%s", requestUUID.String()) // 用于区分本地生成的Request-ID
	//	c.Request.Header.Set("X-Request-Id", requestID)
	//}
	//c.Set(ContextKeyRequestID, requestID)
}

package gdplog

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

const (
	// RequestIDField log field of request id
	RequestIDField = "req_id"
	// QueueMessageIDField log field of queue message id
	QueueMessageIDField = "queue_msg_id"
	// TraceIDField log field of TraceID
	TraceIDField = "trace_id"
	// TargetField log field of target
	TargetField = "target"
	// UserIDField log field of user id
	UserIDField = "user_id"
)

var (
	defaultLogger    *logrus.Logger
	contextKeyLogger = &struct{}{}
)

func init() {
	defaultLogger = New()
}

// SetDefaultLoggerLevel 设置默认Logger的日志等级
func SetDefaultLoggerLevel(level logrus.Level) {
	defaultLogger.SetLevel(level)
}

// New 新建一个Logger
func New() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&nested.Formatter{
		FieldsOrder:     []string{RequestIDField, QueueMessageIDField, TargetField, UserIDField},
		TimestampFormat: "2006/01/02 15:04:05.000",
		HideKeys:        true,
		NoColors:        true,
		ShowFullLevel:   true,
		CallerFirst:     true,
		CustomCallerFormatter: func(frame *runtime.Frame) string {
			return fmt.Sprintf(" [%s:%d]", filepath.Base(frame.File), frame.Line)
		},
	})
	logger.SetReportCaller(true)
	return logger
}

// Logger 根据上下文返回Logger
func Logger(ctx context.Context) logrus.FieldLogger {
	if ctx == nil {
		return defaultLogger
	}

	loggerIntf := ctx.Value(contextKeyLogger)
	if loggerIntf == nil {
		return defaultLogger
	}

	logEntry, ok := loggerIntf.(*logrus.Entry)
	if ok {
		return logEntry
	}

	logger, ok := loggerIntf.(*logrus.Logger)
	if ok {
		return logger
	}

	return loggerIntf.(logrus.FieldLogger)
}

// SetLoggerToContext 将Logger保存到上下文中
func SetLoggerToContext(ctx context.Context, logger logrus.FieldLogger) context.Context {
	newCtx := context.WithValue(ctx, contextKeyLogger, logger)
	return newCtx
}

// DefaultLogger 返回默认Logger
func DefaultLogger() logrus.FieldLogger {
	return defaultLogger
}

package logs

import (
	"context"
	"fmt"
	"github.com/tianlin0/plat-lib/cond"
	"github.com/tianlin0/plat-lib/utils"
	"path/filepath"
)

type Config struct {
	DefaultLogger ILogger
	LogLevel      LogLevel
	WithCtxLogger func(ctx context.Context) (newLogger ILogger, newCtx context.Context) //设置CtxLogger的方法
	CtxLogger     func(ctx context.Context) ILogger
	CommLogString LogString
	LoggerCtxName string
}

var (
	defaultConfig = &Config{
		CommLogString: defaultLogString,
		LoggerCtxName: "context_logger_name",
	}
)

type defaultLogger struct {
	defaultOne ILogger
}

// Debug
func (x *defaultLogger) Debug(v ...interface{}) {
	x.defaultOne.Debug(v...)
}

// Error
func (x *defaultLogger) Error(v ...interface{}) {
	x.defaultOne.Error(v...)
}

// Info
func (x *defaultLogger) Info(v ...interface{}) {
	x.defaultOne.Info(v...)
}

// Warn
func (x *defaultLogger) Warn(v ...interface{}) {
	x.defaultOne.Warn(v...)
}

// Level
func (x *defaultLogger) Level() LogLevel { return x.defaultOne.Level() }

// SetLevel
func (x *defaultLogger) SetLevel(l LogLevel) {
	x.defaultOne.SetLevel(l)
}

// LogId 公共的只能返回空
func (x *defaultLogger) LogId(ctx context.Context) string {
	return ""
}

// SetConfig 设置默认日志，不能包含ctx，不然全局唯一会有问题
func SetConfig(cfg *Config) {
	if cfg == nil {
		return
	}
	if !cond.IsNil(cfg.DefaultLogger) {
		defaultConfig.DefaultLogger = &defaultLogger{
			defaultOne: cfg.DefaultLogger,
		}
	}
	if !cond.IsNil(cfg.CommLogString) {
		defaultConfig.CommLogString = cfg.CommLogString
	}
	if cfg.LoggerCtxName != "" {
		defaultConfig.LoggerCtxName = cfg.LoggerCtxName
	}
	if cfg.CtxLogger != nil {
		defaultConfig.CtxLogger = cfg.CtxLogger
	}
	if cfg.WithCtxLogger != nil {
		defaultConfig.WithCtxLogger = cfg.WithCtxLogger
	}
	if cfg.LogLevel > 0 {
		defaultConfig.LogLevel = cfg.LogLevel
	}

	// 初始化设置
	if defaultConfig.LogLevel > 0 {
		DefaultLogger().SetLevel(defaultConfig.LogLevel)
	}
}

// GetConfig 获取默认配置
func GetConfig() *Config {
	return defaultConfig
}

// DefaultLogger 获取系统默认的日志
func DefaultLogger() ILogger {
	if cond.IsNil(defaultConfig.DefaultLogger) {
		logger := NewPrintLogger(DEBUG)
		logger.SetCallerSkip(3)
		SetConfig(&Config{DefaultLogger: logger})
	}
	return defaultConfig.DefaultLogger
}

// CtxLogger 获取系统默认的日志
func CtxLogger(ctx context.Context) ILogger {
	if ctx == nil {
		ctx = context.Background()
	}
	if cond.IsNil(defaultConfig.CtxLogger) {
		defaultConfig.CtxLogger = defaultCtxLogger
	}
	return defaultConfig.CtxLogger(ctx)
}

// WithCtxLogger 初始化
func WithCtxLogger(ctx context.Context) (newLogger ILogger, newCtx context.Context) {
	if cond.IsNil(defaultConfig.WithCtxLogger) {
		defaultConfig.WithCtxLogger = defaultSetCtxLogger
	}
	return defaultConfig.WithCtxLogger(ctx)
}

func defaultCtxLogger(ctx context.Context) ILogger {
	return GetCtxLogger(ctx)
}
func defaultSetCtxLogger(ctx context.Context) (newLogger ILogger, newCtx context.Context) {
	newLogger, newCtx = NewCtxLogger(ctx, INFO, nil)
	return
}

func defaultLogString(l *LogData) string {
	if l.Message == nil || len(l.Message) == 0 {
		return ""
	}

	fileNameTemp := ""
	if l.FileName != "" {
		fileNameTemp = filepath.Base(l.FileName)
	}

	oldStr := utils.Join(l.Message, " ")
	if fileNameTemp != "" {
		oldStr = fmt.Sprintf("[%s:%d]%s", fileNameTemp, l.Line, oldStr)
	}
	if l.LogId != "" {
		oldStr = fmt.Sprintf("%s: %s", l.LogId, oldStr)
	}

	if l.LogLevel > 0 {
		oldStr = fmt.Sprintf("%s %s", l.LogLevel.GetName(), oldStr)
	}

	return oldStr
}

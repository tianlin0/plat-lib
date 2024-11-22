package database

import (
	"fmt"
	"github.com/tianlin0/plat-lib/logs"
	xormlog "xorm.io/xorm/log"
)

// setXormLogger 设置数据库的日志
func setXormLogger(l interface{}) xormlog.Logger {
	if l == nil {
		return nil
	}
	if xLog, ok := l.(xormlog.Logger); ok {
		return xLog
	}
	if xLog2, ok := l.(logs.ILogger); ok {
		return &xormLogger{
			commLog:    xLog2,
			showSql:    true, //默认打开
			callerSkip: 7,
		}
	}
	return nil
}

// xormLogger 内部的xormLogger
type xormLogger struct {
	commLog    logs.ILogger
	showSql    bool
	callerSkip int
}

func (x *xormLogger) getLogger() logs.ILogger {
	var logger logs.ILogger
	if x.commLog != nil {
		logger = x.commLog
	} else {
		logger = logs.DefaultLogger()
	}
	return logger
}

// Debug 调试
func (x *xormLogger) Debug(v ...interface{}) {
	x.getLogger().Debug(v...)
}

// Debugf 调试
func (x *xormLogger) Debugf(format string, v ...interface{}) {
	x.Debug(fmt.Sprintf(format, v...))
}

// Error 错误
func (x *xormLogger) Error(v ...interface{}) {
	x.getLogger().Error(v...)
}

// Errorf 错误
func (x *xormLogger) Errorf(format string, v ...interface{}) {
	x.Error(fmt.Sprintf(format, v...))
}

// Info 普通
func (x *xormLogger) Info(v ...interface{}) {
	x.getLogger().Info(v...)
}

// Infof 普通
func (x *xormLogger) Infof(format string, v ...interface{}) {
	x.Info(fmt.Sprintf(format, v...))
}

// Warn 警告
func (x *xormLogger) Warn(v ...interface{}) {
	x.getLogger().Warn(v...)
}

// Warnf 警告
func (x *xormLogger) Warnf(format string, v ...interface{}) {
	x.Warn(fmt.Sprintf(format, v...))
}

// Level 等级
func (x *xormLogger) Level() xormlog.LogLevel {
	level := x.commLog.Level()
	return getXormLevelFromLogLever(level)
}

func getXormLevelFromLogLever(level logs.LogLevel) xormlog.LogLevel {
	if level == logs.DEBUG {
		return xormlog.LOG_DEBUG
	}
	if level == logs.INFO {
		return xormlog.LOG_INFO
	}
	if level == logs.WARNING {
		return xormlog.LOG_WARNING
	}
	if level == logs.ERROR {
		return xormlog.LOG_ERR
	}
	if level <= 0 {
		return xormlog.LOG_OFF
	}
	return xormlog.LOG_UNKNOWN
}
func getLevelFromXormLogLever(l xormlog.LogLevel) logs.LogLevel {
	if l == xormlog.LOG_DEBUG {
		return logs.DEBUG
	}
	if l == xormlog.LOG_INFO {
		return logs.INFO
	}
	if l == xormlog.LOG_WARNING {
		return logs.WARNING
	}
	if l == xormlog.LOG_ERR {
		return logs.ERROR
	}
	return logs.LogLevel(10000000)
}

// SetLevel 设置级别
func (x *xormLogger) SetLevel(l xormlog.LogLevel) {
	x.commLog.SetLevel(getLevelFromXormLogLever(l))
}

// ShowSQL 显示sql
func (x *xormLogger) ShowSQL(show ...bool) {
	if len(show) > 0 {
		x.showSql = show[0]
	}
}

// IsShowSQL 是否显示sql
func (x *xormLogger) IsShowSQL() bool {
	return x.showSql
}

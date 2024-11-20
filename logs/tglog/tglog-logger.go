package tglog

import (
	"fmt"
	"github.com/tianlin0/plat-lib/cache"
	"github.com/tianlin0/plat-lib/cond"
	"github.com/tianlin0/plat-lib/conn"
	"github.com/tianlin0/plat-lib/conv"
	"github.com/tianlin0/plat-lib/logs"
	"github.com/tianlin0/plat-lib/utils"
	"log"
	"net"
	"context"
	"strings"
	"time"
)

// tgLogger
type tgLogger struct {
	logLevel  logs.LogLevel
	logServer *TGLoggerV1
}

/*
odpLog, err := tglog.NewTgLogger(valStr)
*/

// NewTgLogger 新建tg连接
func NewTgLogger(con *conn.Connect) (logs.ILogger, error) {
	logServ := &tgLogger{
		logLevel: logs.DEBUG,
	}
	if con == nil || con.Host == "" || con.Port == "" {
		return nil, fmt.Errorf("host or port is null")
	}
	logAddr := net.JoinHostPort(con.Host, con.Port)

	var newOdpLog *TGLoggerV1

	tgLogIns, ok := tgLogInstance.Get(logAddr)
	if ok {
		//正常可用的连接
		if oneLog, ok := tgLogIns.(*TGLoggerV1); ok {
			newOdpLog = oneLog
		}
	} else {
		//表示过期了的，则需要关闭连接
		if tgLogIns != nil {
			//需要关闭连接
			if oneLog, ok := tgLogIns.(*TGLoggerV1); ok {
				err := oneLog.UdpConn.Close()
				if err == nil {
					tgLogInstance.Del(logAddr)
				}
			}
			cacheData := new(cache.DataEntry)
			cacheData.Key = logAddr
			cacheData.Value = tgLogIns
			tgLogInstance.FlushCallback(cacheData, false, true)
		}
	}
	if newOdpLog == nil {
		odpLog, err := NewTGLogV1(logAddr)

		if odpLog == nil || err != nil {
			var logger logs.ILogger = logServ
			if cond.IsNil(logger) {
				return nil, err
			}
			return logger, err
		}
		tgLogInstance.Set(logAddr, odpLog, cacheTimeout)
		newOdpLog = odpLog
	}

	logServ.logServer = newOdpLog

	var logger logs.ILogger = logServ
	if cond.IsNil(logger) {
		return nil, nil
	}
	return logger, nil
}

// SetLevel 设置
func (x *tgLogger) SetLevel(l logs.LogLevel) {
	x.logLevel = l
}

func (x *tgLogger) tgLogInfo(method logs.LogLevel, v ...interface{}) {
	instance := getInstance(v...)
	if instance == nil {
		return
	}
	logAsync(x.logServer, instance)

	logInfo := instance.GetLogString()
	if logInfo != "" {
		logger := logs.NewPrintLogger(method, &logs.LogCommData{LogId: instance.GetLogId()})
		if method == logs.DEBUG {
			logger.Debug(logInfo)
		} else if method == logs.ERROR {
			logger.Error(logInfo)
		} else if method == logs.INFO {
			logger.Info(logInfo)
		} else if method == logs.WARNING {
			logger.Warn(logInfo)
		}
	}
}

// Debug 调试
func (x *tgLogger) Debug(v ...interface{}) {
	if x.logLevel < logs.DEBUG {
		return
	}
	x.tgLogInfo(logs.DEBUG, v...)
}

// error
func (x *tgLogger) Error(v ...interface{}) {
	if x.logLevel < logs.ERROR {
		return
	}
	x.tgLogInfo(logs.ERROR, v...)
}

// Info 普通
func (x *tgLogger) Info(v ...interface{}) {
	if x.logLevel < logs.INFO {
		return
	}
	x.tgLogInfo(logs.INFO, v...)
}

// Warn 警告
func (x *tgLogger) Warn(v ...interface{}) {
	if x.logLevel < logs.WARNING {
		return
	}
	x.tgLogInfo(logs.WARNING, v...)
}

// Level 级别
func (x *tgLogger) Level() logs.LogLevel {
	return x.logLevel
}

func (x *tgLogger) LogId(ctx context.Context) string {
	return ""
}

func getInstance(v ...interface{}) ILogParam {
	if len(v) == 0 {
		return nil
	}
	instance, ok := v[0].(ILogParam)
	if !ok {
		// 如果为空，则用默认值: LogService 来进行处理，这里需要重点注意
		log.Println("tgLog find ToArray()(tableName string, paramList []interface{}) function.")
		logMsg := conv.String(v[0])
		instance = &LogService{
			CurrentTime: time.Now(),
			LogId:       httputil.GetLogId(),
			LogMessage:  logMsg,
		}
	}
	return instance
}

// logAsync tgLog打印日志
func logAsync(logServer *TGLoggerV1, logParam ILogParam) bool {
	tableName, paramArray := logParam.ToArray()
	if tableName == "" || paramArray == nil {
		return false
	}

	newParamArray := make([]interface{}, 0)
	nowTime := time.Now().Format(datetimeLayout)
	newParamArray = append(newParamArray, nowTime) //头要插入当前时间，通用的

	for _, one := range paramArray {
		oneStr := conv.String(one)
		oneStr = strings.ReplaceAll(oneStr, "\r\n", "")
		oneStr = strings.ReplaceAll(oneStr, "\n", "")
		oneStr = strings.ReplaceAll(oneStr, "|", "$")
		newParamArray = append(newParamArray, oneStr)
	}

	if !cond.IsNil(logServer) {
		logServer.LogAsync(tableName, newParamArray...)
	}
	return true
}

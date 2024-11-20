package filelog

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"time"
)

func newLfsHook(logName, endName string, rotationTime time.Duration, maxRemainCnt uint) (logrus.Hook, error) {
	if endName == "" {
		endName = ".%Y%m%d"
	}
	if rotationTime == 0 {
		rotationTime = time.Hour * 24
	}
	if logName == "" {
		logName = "logfile.log"
	}
	if maxRemainCnt == 0 {
		maxRemainCnt = 7
	}

	writer, err := rotatelogs.New(
		logName+endName,
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		rotatelogs.WithLinkName(logName),

		// WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
		rotatelogs.WithRotationTime(rotationTime),

		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数。
		//rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithRotationCount(maxRemainCnt),
	)

	if err != nil {
		return nil, err
	}

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{DisableColors: true})

	return lfsHook, nil
}

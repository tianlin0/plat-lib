package logs

// LogLevel defines a log level
type LogLevel int

const (
	DEBUG     LogLevel = 100
	INFO      LogLevel = 200
	NOTICE    LogLevel = 250
	WARNING   LogLevel = 300
	ERROR     LogLevel = 400
	CRITICAL  LogLevel = 500
	ALERT     LogLevel = 550
	EMERGENCY LogLevel = 600
)

var logLevelMap = map[LogLevel]string{
	DEBUG:     "DEBUG",
	INFO:      "INFO",
	NOTICE:    "NOTICE",
	WARNING:   "WARNING",
	ERROR:     "ERROR",
	CRITICAL:  "CRITICAL",
	ALERT:     "ALERT",
	EMERGENCY: "EMERGENCY",
}

var logLevelNumberMap map[string]LogLevel

func getLevelMap() map[string]LogLevel {
	if len(logLevelNumberMap) > 0 {
		return logLevelNumberMap
	}
	logLevelNumberMap = make(map[string]LogLevel, len(logLevelMap))

	for k, v := range logLevelMap {
		logLevelNumberMap[v] = k
	}
	return logLevelNumberMap
}

// GetName 默认返回空
func (l LogLevel) GetName() string {
	if name, ok := logLevelMap[l]; ok {
		return name
	}
	return ""
}

// GetLogLevel 通过名称取得日志等级
func GetLogLevel(name string) LogLevel {
	levelMap := getLevelMap()

	if code, ok := levelMap[name]; ok {
		return code
	}
	return 0
}

package gdbc

const (
	LOG_TRACE = "TRACE"
	LOG_DEBUG = "DEBUG"
	LOG_INFO  = "INFO"
	LOG_WARN  = "WARN"
	LOG_ERROR = "ERROR"
	LOG_FATAL = "FATAL"
)

type DbLog interface {
	//Set the log parameters
	//@path log base path
	//@prefixFile log file prefix,filename :prefixFile.level.log.yyyymmdd
	//@level Log level, only the level contained in the parameter will be output,all:TRACE|DEBUG|INFO|WARN|ERROR|FATAL
	//@maxSize The maximum length of the file, Exceeding this size will be renamed to .1,.2,.3...
	SetDbLog(path string, prefixFile string, level string, maxSize int64)

	//Set the log level
	SetDbLogLevel(level string)

	//log record head
	GetDbLogHead(level string) string

	//Check the log level when included in the output log level range
	CheckDbLogLevel(level string) bool

	//write log
	WriteDbLog(level string, head string, v ...interface{}) bool

	//write log for format
	WriteDbLogf(level string, head string, format string, v ...interface{}) bool
}

var defaulDbLog DbLog = nil

func SetDefautDbLog(log DbLog) bool {
	defaulDbLog = log
	return true
}

func WriteDbLogf(level string, format string, v ...interface{}) bool {
	if defaulDbLog != nil && defaulDbLog.CheckDbLogLevel(level) {
		return defaulDbLog.WriteDbLogf(level, defaulDbLog.GetDbLogHead(level), format, v...)
	}
	return false
}

func WriteDbLog(level string, v ...interface{}) bool {
	if defaulDbLog != nil && defaulDbLog.CheckDbLogLevel(level) {
		return defaulDbLog.WriteDbLog(level, defaulDbLog.GetDbLogHead(level), v...)
	}
	return false
}

func FatalLog(v ...interface{}) bool {
	if defaulDbLog != nil && defaulDbLog.CheckDbLogLevel(LOG_FATAL) {
		return defaulDbLog.WriteDbLog(LOG_FATAL, defaulDbLog.GetDbLogHead(LOG_FATAL), v...)
	}
	return false
}

func ErrorLog(v ...interface{}) bool {
	if defaulDbLog != nil && defaulDbLog.CheckDbLogLevel(LOG_ERROR) {
		return defaulDbLog.WriteDbLog(LOG_ERROR, defaulDbLog.GetDbLogHead(LOG_ERROR), v...)
	}
	return false
}

func WarnLog(v ...interface{}) bool {
	if defaulDbLog != nil && defaulDbLog.CheckDbLogLevel(LOG_WARN) {
		return defaulDbLog.WriteDbLog(LOG_WARN, defaulDbLog.GetDbLogHead(LOG_WARN), v...)
	}
	return false
}

func InfoLog(v ...interface{}) bool {
	if defaulDbLog != nil && defaulDbLog.CheckDbLogLevel(LOG_INFO) {
		return defaulDbLog.WriteDbLog(LOG_INFO, defaulDbLog.GetDbLogHead(LOG_INFO), v...)
	}
	return false
}

func DebugLog(v ...interface{}) bool {
	if defaulDbLog != nil && defaulDbLog.CheckDbLogLevel(LOG_DEBUG) {
		return defaulDbLog.WriteDbLog(LOG_DEBUG, defaulDbLog.GetDbLogHead(LOG_DEBUG), v...)
	}
	return false
}

func FatalLogf(format string, v ...interface{}) bool {
	if defaulDbLog != nil && defaulDbLog.CheckDbLogLevel(LOG_FATAL) {
		return defaulDbLog.WriteDbLogf(LOG_FATAL, defaulDbLog.GetDbLogHead(LOG_FATAL), format, v...)
	}
	return false
}

func ErrorLogf(format string, v ...interface{}) bool {
	if defaulDbLog != nil && defaulDbLog.CheckDbLogLevel(LOG_ERROR) {
		return defaulDbLog.WriteDbLogf(LOG_ERROR, defaulDbLog.GetDbLogHead(LOG_ERROR), format, v...)
	}
	return false
}

func WarnLogf(format string, v ...interface{}) bool {
	if defaulDbLog != nil && defaulDbLog.CheckDbLogLevel(LOG_WARN) {
		return defaulDbLog.WriteDbLogf(LOG_WARN, defaulDbLog.GetDbLogHead(LOG_WARN), format, v...)
	}
	return false
}

func InfoLogf(format string, v ...interface{}) bool {
	if defaulDbLog != nil && defaulDbLog.CheckDbLogLevel(LOG_INFO) {
		return defaulDbLog.WriteDbLogf(LOG_INFO, defaulDbLog.GetDbLogHead(LOG_INFO), format, v...)
	}
	return false
}

func DebugLogf(format string, v ...interface{}) bool {
	if defaulDbLog != nil && defaulDbLog.CheckDbLogLevel(LOG_DEBUG) {
		return defaulDbLog.WriteDbLogf(LOG_DEBUG, defaulDbLog.GetDbLogHead(LOG_DEBUG), format, v...)
	}
	return false
}

func SetDbLog(path string, file string, level string, maxSize int64) {
	if defaulDbLog != nil {
		defaulDbLog.SetDbLog(path, file, level, maxSize)
	}
}

func SetDbLogLevel(level string) {
	if defaulDbLog != nil {
		defaulDbLog.SetDbLogLevel(level)
	}
}

package gdbc

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)


type DefaultLog struct {
	path       string
	file       string
	level      string
	filemutex  sync.Mutex
	tracemutex sync.Mutex
	maxSize    int64
}

func NewDbLog(path string, file string, level string, maxSize int64) *DefaultLog {

	return &DefaultLog{path: path, file: file, level: level, maxSize: maxSize}
}

func (log *DefaultLog) SetDbLog(path string, file string, level string, maxSize int64) {
	log.path = path
	log.file = file
	log.level = level
	log.maxSize = maxSize
}

func (log *DefaultLog) SetDbLogLevel(level string) {
	log.level = level
}

func (log *DefaultLog) CheckDbLogLevel(level string) bool {
	return strings.Contains(log.level, level)
}

func (log *DefaultLog) getDbLogFileName(level string) (fileName string, ok bool) {
	if len(log.file) < 1 {
		fmt.Printf("GetLogFileName error! filename(%s) is null", fileName)
		return "", false
	}

	if len(log.path) < 1 {
		log.path = "."
	}

	today := time.Now()
	fileName = fmt.Sprintf("%s/%s.%s.log.%d%02d%02d", log.path, log.file, level,
		today.Year(), int(today.Month()), today.Day())

	return fileName, true
}

func (log *DefaultLog) checkAndBackupFile(fileName string) bool {
	fileInfo, err := os.Stat(fileName)
	if err != nil && os.IsExist(err) {
		fmt.Printf("check file %s error!%s\n", fileName, err.Error())
		return false
	}

	if os.IsNotExist(err) {
		return true
	}

	//fmt.Println("filesize=", fileInfo.Size())
	if fileInfo.Size() >= log.maxSize {
		for i := 1; ; i++ {
			backupFileName := fmt.Sprintf("%s.%d", fileName, i)
			_, err := os.Stat(backupFileName)
			if err != nil && os.IsNotExist(err) {
				//fmt.Printf("check file not found!%s\n", backupFileName)
				err = os.Rename(fileName, backupFileName)
				if err != nil {
					fmt.Printf("rename file %s to %s error!%s\n", fileName, backupFileName, err)
					return false
				}
				break
			}

		}
	}

	return true
}

func getSourceLine(skip int) (file string, line int) {
	var ok bool
	_, file, line, ok = runtime.Caller(skip)
	if !ok {
		file = "???"
		line = 0
	}

	return file, line
}

func (log *DefaultLog) GetDbLogHead(level string) string {
	now := time.Now()
	source, line := getSourceLine(3)
	return fmt.Sprintf("%02d%02d%02d.%04d!%d#%s,%d:", now.Hour(), now.Minute(), now.Second(), now.Nanosecond(),
		os.Getpid(), source, line)
}

func (log *DefaultLog) openDbLogFile(level string) (*os.File, bool) {
	filename, ok := log.getDbLogFileName(level)
	if !ok {
		return nil, false
	}
	ok = log.checkAndBackupFile(filename)
	if !ok {
		return nil, false
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println("open log file %s error:%s", filename, err.Error())
		return nil, false
	}
	return file, true
}

func (log *DefaultLog) closeDbLogFile(file *os.File) {
	if file != nil {
		file.Close()
	}
}

func (log *DefaultLog) WriteDbLog(level string, head string, v ...interface{}) bool {
	if log.CheckDbLogLevel(level) {
		log.filemutex.Lock()
		defer log.filemutex.Unlock()

		file, ok := log.openDbLogFile(level)
		if !ok {
			return false
		}
		defer log.closeDbLogFile(file)

		fmt.Fprint(file, head)
		fmt.Fprintln(file, v...)

	}
	if log.CheckDbLogLevel(LOG_TRACE) {
		log.tracemutex.Lock()
		defer log.tracemutex.Unlock()
		fmt.Print(head)
		fmt.Println(v...)
	}
	return true
}

func (log *DefaultLog) WriteDbLogf(level string, head string, format string, v ...interface{}) bool {
	if log.CheckDbLogLevel(level) {
		log.filemutex.Lock()
		defer log.filemutex.Unlock()

		file, ok := log.openDbLogFile(level)
		if !ok {
			return false
		}
		defer log.closeDbLogFile(file)

		fmt.Fprint(file, head)
		fmt.Fprintf(file, format, v...)
		fmt.Fprintln(file)

	}
	if log.CheckDbLogLevel(LOG_TRACE) {
		log.tracemutex.Lock()
		defer log.tracemutex.Unlock()
		fmt.Print(head)
		fmt.Printf(format, v...)
		fmt.Println()
	}

	return true
}



var _ bool = SetDefautDbLog(NewDbLog(".", "db", "TRACE|DEBUG|INFO|WARN|ERROR|FATAL", 5*1024*1024))


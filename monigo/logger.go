package monigo

import (
	"errors"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"sync/atomic"
)

const (
	LogLevelTrace = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelAdmin
	LogLevelFatal

	LogLevelLast
)

var (
	CurrentLogLevel int32

	adminLog   *log.Logger
	adminQueue chan [2]interface{}

	commonLog   *log.Logger
	commonQueue chan [2]interface{}

	consoleLog   *log.Logger
	consoleQueue chan [2]interface{}

	LogLevelName [LogLevelLast]string
)

func init() {
	LogLevelName[LogLevelTrace] = "trace"
	LogLevelName[LogLevelDebug] = "debug"
	LogLevelName[LogLevelInfo] = "info"
	LogLevelName[LogLevelWarn] = "warn"
	LogLevelName[LogLevelError] = "error"
	LogLevelName[LogLevelAdmin] = "admin"
	LogLevelName[LogLevelFatal] = "fatal"

	adminLog = OpenLog("admin.log")
	adminQueue = make(chan [2]interface{}, 10)
	commonLog = OpenLog("app.log")
	commonQueue = make(chan [2]interface{}, 10)
	consoleLog = OpenLog("")
	consoleQueue = make(chan [2]interface{}, 10)
	CurrentLogLevel = LogLevelTrace

	go writer()
}

func OpenLog(fileName string) *log.Logger {
	if len(fileName) == 0 {
		return log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	}

	l := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     10,
	}

	return log.New(l, "", log.Ldate|log.Ltime|log.Lmicroseconds)
}

func SetLogLevel(l int32) error {
	if l < LogLevelTrace || l >= LogLevelLast {
		return errors.New(fmt.Sprintf("Invalid log level value %d", l))
	}
	atomic.StoreInt32(&CurrentLogLevel, l)
	return nil
}

func GetLevelNameToId(level string) int32 {
	for i := 0; i < LogLevelLast; i++ {
		if LogLevelName[i] == level {
			return int32(i)
		}
	}

	return -1
}

func SetLogOptions(config LogConfig) {
	level := GetLevelNameToId(config.Level)
	if level >= 0 {
		err := SetLogLevel(level)
		if err != nil {
			LogCritical("Err set log level: %s", err)
		}
	}

	if len(config.File) > 0 {
		commonLog.SetOutput(&lumberjack.Logger{
			Filename:   config.File,
			MaxSize:    config.MaxSizeMB,
			MaxBackups: config.MaxSizeBackups,
			MaxAge:     config.MaxAgeDays,
		})
	}

	if len(config.FileCritical) > 0 {
		adminLog.SetOutput(&lumberjack.Logger{
			Filename:   config.FileCritical,
			MaxSize:    config.MaxSizeMB,
			MaxBackups: config.MaxSizeBackups,
			MaxAge:     config.MaxAgeDays,
		})
	}
}

func LogWrite(format string, args ...interface{}) {
	commonLog.Printf(format, args...)
}

func ConsoleWrite(format string, args ...interface{}) {
	consoleLog.Printf(format, args...)
}

func writer() {
	for {
		select {
		case row := <-commonQueue:
			LogWrite(ClearColors(row[0].(string)), row[1].([]interface{})...)

		case row := <-adminQueue:
			AdminLogWrite(ClearColors(row[0].(string)), row[1].([]interface{})...)

		case row := <-consoleQueue:
			ConsoleWrite(row[0].(string), row[1].([]interface{})...)
		}
	}
}

func AdminLogWrite(format string, args ...interface{}) {
	adminLog.Printf(ClearColors(format), args...)
}

func Log(level int32, format string, args ...interface{}) {
	color := ""
	switch {

	case level == LogLevelInfo:
		color = Text.Sky

	case level == LogLevelDebug:
		color = Text.Yellow

	case level == LogLevelError || level == LogLevelFatal:
		color = Text.Red

	}

	format = color + LogLevelName[level] + Text.Reset + " " + Text.Gray + getStrGID() + Text.Reset + " " + format

	if level >= LogLevelAdmin {
		consoleQueue <- [2]interface{}{format, args}
		adminQueue <- [2]interface{}{format, args}
	}

	if level >= atomic.LoadInt32(&CurrentLogLevel) {
		consoleQueue <- [2]interface{}{format, args}
		commonQueue <- [2]interface{}{format, args}
	}
}

func LogDebug(format string, args ...interface{}) {
	Log(LogLevelDebug, format, args...)
}

func LogInfo(format string, args ...interface{}) {
	Log(LogLevelInfo, format, args...)
}

func LogError(format string, args ...interface{}) {
	Log(LogLevelError, format, args...)
}

func LogCritical(format string, args ...interface{}) {
	Log(LogLevelAdmin, format, args...)
}

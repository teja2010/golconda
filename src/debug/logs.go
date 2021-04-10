package debug

import (
	"fmt"
	"time"
)

const (
	_DebugLogs   = 0
	_InfoLogs    = 1
	_WarningLogs = 2

	_StackMaxlen = 20
)

// TODO move to config file
var _LogLevel = _InfoLogs

var _logger loggerInterface

type loggerInterface interface {
	Println(v ...interface{})
}

func header(log string) []interface{} {
	caller := getCallerFunc(3)
	callStr := caller.logPrefix()

	header := []interface{}{"[" + log + callStr + "]"}
	return header
}

// DebugLog is only for debugging
func DebugLog(v ...interface{}) {
	if _LogLevel > _DebugLogs {
		return
	}

	v = append(header("DEBUG: "), v...)

	_logger.Println(v...)
}

//Log is to inform about general event
func Log(v ...interface{}) {
	if _LogLevel > _InfoLogs {
		return
	}

	v = append(header("INFO: "), v...)

	_logger.Println(v...)
}

// Error is for errors
func Error(v ...interface{}) {
	v = append(header("ERROR: "), v...)
	_logger.Println(v...)
}

// Warning is about unusal events
func Warning(v ...interface{}) {
	v = append(header("WARNING: "), v...)
	_logger.Println(v...)
}

// Bug is to log and crash
func Bug(v ...interface{}) {
	s := fmt.Sprintln(v...)

	caller := getCallerFunc(2)
	header := []interface{}{"BUG: ", caller}
	v = append(header, v...)
	_logger.Println(v...)

	_logger.Println("Stack:")
	callers := getCallerFunctions()
	for i, caller := range callers {
		_logger.Println("#", i, "|", caller)
	}

	panic(s)
}

// InitLogging inits logging
func InitLogging(updateLevel func() int) {
	_logger = defLogInit()

	go func() {
		for {
			time.Sleep(time.Second)
			_LogLevel = updateLevel()
		}
	}()

}

func SetLogLevel(lvl int) {
	_LogLevel = lvl
}

func ToString(v interface{}) string {
	return fmt.Sprintf("%+v", v)
}

package debug

import (
	"os"
)

const (
	DEBUG_LOGS = 0
	INFO_LOGS = 1
	WARNING_LOGS = 2

	LOG_LEVEL = DEBUG_LOGS

	STACK_MAXLEN = 20
)

var _logger loggerInterface

type loggerInterface interface {
	Println(v ...interface{})
}

func DebugLog (v ...interface{}) {
	if LOG_LEVEL > DEBUG_LOGS {
		return
	}

	header := []interface{}{"DEBUG:"}
	v = append(header, v...)

	_logger.Println(v ...)
}

func Log(v ...interface{}) {
	if LOG_LEVEL > INFO_LOGS {
		return
	}

	_logger.Println(v ...)
}

func Error(v ...interface{}) {
	header := []interface{}{"ERROR:"}
	v = append(header, v...)
	_logger.Println(v ...)
}

func Warning(v ...interface{}) {
	call_site := getCallSite(3)
	header := []interface{}{"WARNING: ", call_site}

	v = append(header, v...)
	_logger.Println(v ...)
}

func Bug(v ...interface{}) {
	PrintStackToStdErr()

	caller := getCallerFunc(1)
	header := []interface{}{"BUG: ", caller}
	v = append(header, v...)
	_logger.Println(v ...)

	_logger.Println("Stack:")
	callers := getCallerFunctions()
	for i, caller := range callers {
		_logger.Println("#", i, "|", caller)
	}

	os.Exit(-1)
}

func InitLogging() {
	_logger = DefLogInit()
}

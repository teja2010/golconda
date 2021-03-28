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

func header(log string) []interface{} {
	caller := getCallerFunc(3)
	call_str := caller.logPrefix()

	header := []interface{}{"[" + log + call_str + "]"}
	return header
}

func DebugLog (v ...interface{}) {
	if LOG_LEVEL > DEBUG_LOGS {
		return
	}

	v = append(header("DEBUG: "), v...)

	_logger.Println(v ...)
}

func Log(v ...interface{}) {
	if LOG_LEVEL > INFO_LOGS {
		return
	}

	v = append(header("INFO: "), v...)

	_logger.Println(v ...)
}

func Error(v ...interface{}) {
	v = append(header("ERROR: "), v...)
	_logger.Println(v ...)
}

func Warning(v ...interface{}) {
	v = append(header("WARNING: "), v...)
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

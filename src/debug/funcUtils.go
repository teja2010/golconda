package debug

import (
	"fmt"
	"runtime"
	runtime_debug "runtime/debug"
	"strconv"
	"strings"
)

type logFunc struct {
	pc       uintptr
	name     string
	fileName string
	line     int
}

func (f logFunc) String() string {
	str := fmt.Sprintf("%x %s %s:%d", f.pc, f.name,
		f.fileName, f.line)
	return str
}
func removePathPrefix(name string, seperator string) string {
	lastIdx := strings.LastIndex(name, seperator)
	if lastIdx > 0 {
		name = name[lastIdx+1:]
	}

	return name
}

func (f logFunc) logPrefix() string {
	str := fmt.Sprintf("%s %s:%d",
		removePathPrefix(f.name, `.`),
		removePathPrefix(f.fileName, `/`),
		f.line)
	return str
}

func invalidLogFunc(pc uintptr) logFunc {
	return logFunc{pc, "INVALID_LOGFUNC", "", -1}
}

func getCallerFunc(skip int) logFunc {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return invalidLogFunc(0)
	}

	return getLogFunc(pc)
}

func getCallSite(skip int) string {
	caller := getCallerFunc(skip)
	callSite := caller.name + ":" + strconv.Itoa(caller.line)
	return callSite
}

func getLogFunc(pc uintptr) logFunc {
	funcDetails := runtime.FuncForPC(pc)
	if funcDetails == nil {
		return invalidLogFunc(0)
	}

	fileName, lineNum := funcDetails.FileLine(pc)

	return logFunc{pc, funcDetails.Name() + "()", fileName, lineNum}
}

func getCallerFunctions() []logFunc {

	pcSlice := make([]uintptr, _StackMaxlen)
	numPc := runtime.Callers(2, pcSlice)

	pcSlice = pcSlice[:numPc]

	funcs := make([]logFunc, numPc)

	for i, pc := range pcSlice {
		funcs[i] = getLogFunc(pc)
	}

	return funcs
}

// PrintStackToStdErr will print the current stack to stderr
func PrintStackToStdErr() {
	runtime_debug.PrintStack()
}

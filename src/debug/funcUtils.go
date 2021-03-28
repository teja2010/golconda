package debug

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	runtime_debug "runtime/debug"
)

type logFunc struct {
	pc uintptr
	name string
	file_name string
	line int
}

func (f logFunc) String() string {
	str := fmt.Sprintf("%x %s %s:%d", f.pc, f.name,
			f.file_name, f.line);
	return str
}
func removePathPrefix (name string, seperator string) string {
	last_idx := strings.LastIndex(name, seperator)
	if last_idx > 0 {
		name = name[last_idx+1:]
	}

	return name
}


func (f logFunc) logPrefix() string {
	str := fmt.Sprintf("%s %s:%d",
			   removePathPrefix(f.name, `/`),
			   removePathPrefix(f.file_name, `/`),
			   f.line);
	return str
}

func invalidLogFunc(pc uintptr) logFunc{
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
	call_site := caller.name + ":" + strconv.Itoa(caller.line)
	return call_site
}

func getLogFunc(pc uintptr) logFunc {
	func_details := runtime.FuncForPC(pc)
	if func_details == nil {
		return invalidLogFunc(0)
	}

	file_name, line_num := func_details.FileLine(pc)

	return logFunc{pc, func_details.Name() + "()", file_name, line_num}
}

func getCallerFunctions() []logFunc {

	pc_slice := make([]uintptr, STACK_MAXLEN)
	num_pc := runtime.Callers(2, pc_slice)

	pc_slice = pc_slice[:num_pc]

	funcs := make([]logFunc, num_pc)

	for i, pc := range pc_slice {
		funcs[i] = getLogFunc(pc)
	}

	return funcs
}

func PrintStackToStdErr() {
	runtime_debug.PrintStack()
}

package golconda

import (
	"time"
	. "github.com/teja2010/golconda/src/any"
	d "github.com/teja2010/golconda/src/debug"
	ui "github.com/teja2010/golconda/src/ui"
)

type RegisteredFunction func(chan<- ui.PrintData)

// TODO remove it
func testFunc(c chan<- ui.PrintData) {
	i:=0
	for {
		time.Sleep(1*time.Second)
		c <- ui.PrintData{ui.Tuple{i,0}, ui.Tuple{0,0}, []string{}}
		i++
	}
}

func registeredFunctions() []RegisteredFunction {
	return []RegisteredFunction{
			//testFunc,
			CPU_Usage,
			Meminfo,
		}
}

func Start(argsMap map[string]AnyValue) {
	d.DebugLog("Parsed", argsMap)

	printChan := make(chan ui.PrintData, 10)

	ui.Init()

	for _, f := range registeredFunctions() {
		go f(printChan)
	}

	for {
		pdata := <-printChan
		d.DebugLog("Print Data", pdata)
		ui.Update(pdata)
	}
}

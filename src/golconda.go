package golconda

import (
	d "github.com/teja2010/golconda/src/debug"
	ui "github.com/teja2010/golconda/src/ui"
)

// Start golconda
func Start() {
	d.InitLogging(func() int { return GetConfig().Global.DebugLevel })
	ConfigInit()

	printChan := make(chan ui.PrintData, 10)

	ui.Init()

	for _, f := range RegisteredFunctions() {
		go f(printChan)
	}

	for {
		pdata := <-printChan
		d.DebugLog("Print Data", pdata)
		ui.Update(pdata)
	}
}

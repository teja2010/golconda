package main

import (
	"github.com/teja2010/golconda/src/args"
	"github.com/teja2010/golconda/src"
	conf "github.com/teja2010/golconda/src/config"
	d "github.com/teja2010/golconda/src/debug"
)

func main() {
	d.InitLogging()
	conf.ConfigInit()

	argsMap := args.ParseArgs()

	golconda.Start(argsMap)
}

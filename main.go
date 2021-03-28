package main

import (
	golconda "github.com/teja2010/golconda/src"
	d "github.com/teja2010/golconda/src/debug"
)

func main() {
	d.InitLogging()
	golconda.Start()
}

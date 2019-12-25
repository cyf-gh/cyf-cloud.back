package main

import (
	"runtime"
	stlog "stgogo/comn/log"
	"strconv"
)

func initGlobal() {
	stlog.Info.Println("Server Started...")
	stlog.Info.Println( "goroutine Run with Core Count: " + strconv.Itoa(runtime.GOMAXPROCS(runtime.NumCPU())))
}

func main() {
	initGlobal()
}

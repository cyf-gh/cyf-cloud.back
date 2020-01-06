package main

import (
	vtConfig "./config"
	"github.com/kpango/glg"
	"runtime"
	"strconv"
)

func initGlobal() {
	glg.Info("Server Started...")
	glg.Info( "goroutine Run with Core Count: " + strconv.Itoa(runtime.GOMAXPROCS(runtime.NumCPU())))
}

func main() {
	vtConfig.ConfigAll()
	initGlobal()

	// TODO

	vtConfig.Dispose()
}

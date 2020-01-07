package main

import (
	vtConfig "./config"
	vtServer "./sync"
	"github.com/kpango/glg"
	"runtime"
	"stgogo/comn/config"
	"strconv"
)

func main() {
	glg.Info("Server Started...")
	glg.Info( "goroutine Run with Core Count: " + strconv.Itoa(runtime.GOMAXPROCS(runtime.NumCPU())))
	vtConfig.ConfigAll()

	st_config_log.Start()

	vtServer.Run()

	st_config_log.End()
}

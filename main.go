package main

import (
	vtConfig "./config"
	vtServer "./sync"
	"github.com/kpango/glg"
	"runtime"
	"stgogo/comn/config"
	"strconv"
	"sync"
	"time"
)

func main() {
	glg.Info("Server Started...")
	glg.Info( "goroutine Run with Core Count: " + strconv.Itoa(runtime.GOMAXPROCS(runtime.NumCPU())))
	vtConfig.ConfigAll()

	st_config_log.Start( vtConfig.VtLogAddr, time.Duration( vtConfig.VtFreshLogInterval ) )

	lock := &sync.Mutex{}
	go vtServer.RunTcpServer( vtConfig.VtTcpAddr, lock )
	vtServer.RunUDPSyncServer( vtConfig.VtUdpAddr, time.Duration( vtConfig.VtFreshUdpInterval ) )

	st_config_log.End()
}

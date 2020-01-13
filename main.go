package main

import (
	vtConfig "./config"
	vtServer "./sync"
	"github.com/kpango/glg"
	"runtime"
	"stgogo/comn/config"
	"strconv"
	"time"
)

func main() {
	glg.Info("Server Started...")
	glg.Info( "goroutine Run with Core Count: " + strconv.Itoa(runtime.GOMAXPROCS(runtime.NumCPU())))
	vtConfig.ConfigAll()

	st_config_log.Start( vtConfig.VtLogAddr, time.Duration( vtConfig.VtFreshLogInterval ) )

	go vtServer.RunUDPSyncServer( vtConfig.VtUdpAddr, vtServer.Lobbies, time.Duration( vtConfig.VtFreshUdpInterval ) )
	vtServer.RunTcpServer( vtConfig.VtTcpAddr )

	st_config_log.End()
}

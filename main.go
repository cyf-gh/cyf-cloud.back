package main

import (
	"fmt"
	"github.com/kpango/glg"
	"runtime"
	"stgogo/comn/config"
	"strconv"

	Config  "./config"
	ccHttp "./http"
)

func main() {
	glg.Info("Server Starting...")
	glg.Info( "{goroutine} Run with Core Count: " + strconv.Itoa(runtime.GOMAXPROCS(runtime.NumCPU())))
	// 进行配置
	Config.ConfigAll()
	// 启动http服务器
	go ccHttp.RunHttpServer( Config.TcpAddr )

	/// ======================= proc input ===========================
	var input string
	for {
		fmt.Scanln(&input)
		glg.Log(input)
	}
	st_config_log.End()
}

package main

import (
	"github.com/kpango/glg"
	"runtime"
	"stgogo/comn/config"
	"strconv"

	"./cli"
	Config "./config"
	ccHttp "./server"
)

func main() {
	_ = cli.PrintBanner(nil)

	glg.Info("server starting...")
	glg.Info( "{goroutine} run with Core Count: " + strconv.Itoa(runtime.GOMAXPROCS(runtime.NumCPU())))
	// 进行配置
	Config.All()
	// 启动http服务器
	go ccHttp.RunHttpServer( Config.TcpAddr, Config.RedisCfg, Config.SqlitePath )

	// 启动命令
	cli.Run()

	st_config_log.End()
}

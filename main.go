package main

import (
	"./cli"
	Config "./config"
	"./dm_1"
	ccHttp "./server"
	"github.com/kpango/glg"
	"runtime"
	"stgogo/comn/config"
	"strconv"
)

func main() {
	_ = cli.PrintBanner(nil)

	glg.Info("server starting...")
	glg.Info( "{goroutine} run with Core Count: " + strconv.Itoa(runtime.GOMAXPROCS(runtime.NumCPU())))
	// 进行配置
	Config.All()

	// run dm file monitor
	dnf := dm_1.DMNewNotifyFile()
	dnf.AddWatchDirRecruit( "L:/mount/" )
	go dnf.WatchEvent()

	// 启动http服务器
	go ccHttp.RunHttpServer( Config.TcpAddr, Config.RedisCfg, Config.SqlitePath )

	// 启动命令
	cli.Run()

	st_config_log.End()
}

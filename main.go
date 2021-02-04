package main

import (
	"./cli"
	Config "./config"
	ccHttp "./server"
	"./v1x1/orm"
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

	// 运行数据库协同watcher
	dnf := orm.DMNewNotifyFileDb()
	e := dnf.AddWatchDirRecruit( "L:/mount/" ); if e != nil { glg.Error( e ); glg.Warn( "watch event may not work properly") }
	go dnf.WatchEvent()

	// 启动http服务器
	go ccHttp.RunHttpServer( Config.TcpAddr, Config.RedisCfg, Config.SqlitePath )

	// 启动命令
	cli.Run()

	st_config_log.End()
}

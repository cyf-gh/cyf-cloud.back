// this file will complete configurations below
// server configuration
// log configuration
package config

import (
	"fmt"
	"github.com/kpango/glg"
	"gopkg.in/ini.v1"
	"os"
	stConfigLog "stgogo/comn/config"
	"time"
)

var (
	UdpAddr string
	TcpAddr string
	LogAddr string
	FreshLogInterval int64
 	FreshUdpInterval int64
 	SqlitePath string
 	RedisCfg RedisConfig
	RunMode string	// run_mode or dep
	DMGodId int64
	DMRootPath string
)
type RedisConfig struct {
	Addr string
	MaxIdle, MaxActive int
}

func IsRunModeDev() bool {
	return RunMode == "dev"
}

func IsRunModeDep() bool {
	return RunMode == "dep"
}

func configServerInfo() {
	var (
		cfg *ini.File
		err error
	)
	if cfg, err = ini.Load("./server.cfg"); err != nil {
		panic( err )
	}

	LogAddr = cfg.Section("server_address").Key("log").String()
	TcpAddr = cfg.Section("server_address").Key("tcp").String()
	UdpAddr = cfg.Section("server_address").Key("udp").String()

	FreshUdpInterval, _ = cfg.Section("fresh_interval").Key("udp").Int64()
	FreshLogInterval, _ = cfg.Section("fresh_interval").Key("log").Int64()

	RedisCfg.Addr = cfg.Section("redis").Key("address").String()
	RedisCfg.MaxIdle, _ = cfg.Section("redis").Key("max_idle").Int()
	RedisCfg.MaxActive, _ = cfg.Section("redis").Key("max_active").Int()

	RunMode = cfg.Section("common").Key("mode").String()
	println("server start with mode:\"" + RunMode + "\"")

	SqlitePath =  cfg.Section("sqlite3").Key("path").String()

	DMGodId, _ = cfg.Section("dm_whitelist").Key("god_id").Int64()
	DMRootPath = cfg.Section("dm_whitelist").Key("root_path").String()
	println( " *************** DM configuration loaded... ***************" )
	println( "\tgod ID:\t", DMGodId )
	println( "\troot path:\t" + DMRootPath )
	println( " **********************************************************" )

	defer func() {
		if err := recover(); err != nil {
			// 配置文件读取错误，直接退出程序
			print("Loading server.cfg with err:")
			print(err)
			os.Exit(1)
		}
	}()
}

func All() {
	// stgogo log
	// 必须启动，否则服务器不允许启动
	// TODO: 尚未检查log是否成功启动
	stConfigLog.Start( LogAddr, time.Duration( FreshLogInterval ) )

	configServerInfo()
	glg.Info(fmt.Sprintf("\nTCP: %s\nUDP: %s\nLOG: %s\n",  TcpAddr, UdpAddr, LogAddr ))
	glg.Info(fmt.Sprintf("\nLog Interval: %d\nUdp Interval: %d\n", FreshLogInterval, FreshUdpInterval))
	glg.Info("Redis configs:")
	glg.Info( RedisCfg )
	// 2020.9.4
	glg.Warn("UDP is never used and TCP is server port")
}


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
	Err "../v1x1/err"
)

var UdpAddr string
var TcpAddr string
var LogAddr string
var FreshLogInterval int64
var FreshUdpInterval int64

func configServerInfo() {
	cfg, err := ini.Load("./server.cfg")
	Err.CheckErr(err)

	LogAddr = cfg.Section("server").Key("log").String()
	TcpAddr = cfg.Section("server").Key("tcp").String()
	UdpAddr = cfg.Section("server").Key("udp").String()

	FreshUdpInterval, _ = cfg.Section("fresh_interval").Key("udp").Int64()
	FreshLogInterval, _ = cfg.Section("fresh_interval").Key("log").Int64()

	defer func() {
		if err := recover(); err != nil {
			// 配置文件读取错误，直接退出程序
			print("Loading server.cfg with err:")
			print(err)
			os.Exit(-1)
		}
	}()
}

func ConfigAll() {
	// stgogo log
	// 必须启动，否则服务器不允许启动
	// TODO: 尚未检查log是否成功启动
	stConfigLog.Start( LogAddr, time.Duration( FreshLogInterval ) )

	configServerInfo()
	glg.Info(fmt.Sprintf("\nTCP: %s\nUDP: %s\nLOG: %s\n",  TcpAddr, UdpAddr, LogAddr ))
	glg.Info(fmt.Sprintf("\nLog Interval: %d\nUdp Interval: %d\n", FreshLogInterval, FreshUdpInterval))
	// 2020.9.4
	glg.Warn("UDP is never used and TCP is http port")
}


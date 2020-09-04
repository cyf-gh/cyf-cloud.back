// this file will complete configurations below
// server configuration
// log configuration
package config

import (
	"fmt"
	"github.com/kpango/glg"
	"gopkg.in/ini.v1"
	st_config_log "stgogo/comn/config"
	"time"
)

var UdpAddr string
var TcpAddr string
var LogAddr string
var FreshLogInterval int64
var FreshUdpInterval int64

func configServerInfo() {
	cfg, _ := ini.Load("./server.cfg")
	// glg.Error(err)
	LogAddr = cfg.Section("server").Key("log").String()
	TcpAddr = cfg.Section("server").Key("tcp").String()
	UdpAddr = cfg.Section("server").Key("udp").String()

	FreshUdpInterval, _ = cfg.Section("fresh_interval").Key("udp").Int64()
	FreshLogInterval, _ = cfg.Section("fresh_interval").Key("log").Int64()
}

func ConfigAll() {
	// stgogo log
	st_config_log.Start( LogAddr, time.Duration( FreshLogInterval ) )

	configServerInfo()
	glg.Info(fmt.Sprintf("\nTCP: %s\nUDP: %s\nLOG: %s\n",  TcpAddr, UdpAddr, LogAddr ))
	glg.Info(fmt.Sprintf("\nLog Interval: %s\nUdp Interval: %s\n", FreshLogInterval, FreshUdpInterval))
	// 2020.9.4
	glg.Info("UDP is never used and TCP is http port")
}


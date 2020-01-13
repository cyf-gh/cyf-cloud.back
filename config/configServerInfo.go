// this file will complete configurations below
// server configuration
// log configuration
package vt_config

import (
	"fmt"
	"github.com/kpango/glg"
	"gopkg.in/ini.v1"
)

var VtUdpAddr string
var VtTcpAddr string
var VtLogAddr string
var VtFreshLogInterval int64
var VtFreshUdpInterval int64

func configServerInfo() {
	cfg, _ := ini.Load("./vtserver.cfg")
	// glg.Error(err)
	VtLogAddr = cfg.Section("server").Key("log").String()
	VtTcpAddr = cfg.Section("server").Key("tcp").String()
	VtUdpAddr = cfg.Section("server").Key("udp").String()

	VtFreshUdpInterval, _ = cfg.Section("fresh_interval").Key("udp").Int64()
	VtFreshLogInterval, _ = cfg.Section("fresh_interval").Key("log").Int64()
}

func ConfigAll() {
	configServerInfo()
	glg.Info(fmt.Sprintf("\nTCP: %s\nUDP: %s\nLOG: %s\n",  VtTcpAddr, VtUdpAddr, VtLogAddr ))
}


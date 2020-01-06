package vt_lobby

import (
	vtConfig "../config"
	vtsync "../sync"
	"fmt"
	"github.com/kpango/glg"
	"net"
	"strconv"
)

// start sync
func StartLobby( lobby *VTLobby ) {
	glg.Info(
		lobby.Name +
			" Lobby started with viewers count:" +
			strconv.Itoa( len( lobby.Viewers ) ) +
			" Max Location Offset: " +
			strconv.Itoa( lobby.MaxOffset ) )

	conn, _ := vtsync.Init()
	defer conn.Close()
	vtsync.Proc( conn, 1, lobby )
}

func WaitForBorrower() {
	listener, err := net.Listen("tcp", vtConfig.VtTcpAddr)
	defer listener.Close()
	if err != nil {
		glg.Error(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			glg.Error(err)
			continue
		}
		glg.Info("A Client connected :" + conn.RemoteAddr().String())
		go tcpPipe(conn)
	}
}

func AskForHostViewer( viewers []VTViewer) *VTViewer {
	for i, viewer := range viewers  {
		if viewer.IsHost {
			return &viewers[i]
		}
	}
	return nil
}
package vt_lobby

import (
	vtsync "../sync"
	stlog "stgogo/comn/log"
	"strconv"
)

type VTViewer struct {
	Name, Location string
	IsHost, IsPause bool
}

func GetHostViewer( viewers []VTViewer) *VTViewer {
	for i, viewer := range viewers  {
		if viewer.IsHost {
			return &viewers[i]
		}
	}
	return nil
}

type VTLobby struct {
	Name, Password string
	Viewers []VTViewer
	MaxOffset int
}

// start sync
func StartLobby( lobby *VTLobby ) {
	stlog.Info.Println(
		lobby.Name +
			" Lobby started with viewers count:" +
			strconv.Itoa( len( lobby.Viewers ) ) +
			" Max Location Offset: " +
			strconv.Itoa( lobby.MaxOffset ) )

	conn, _ := vtsync.Init()
	defer conn.Close()
	vtsync.Proc( conn, 1, lobby )
}
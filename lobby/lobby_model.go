package vt_lobby

import (
	"time"
)

type VTViewer struct {
	Name, Location string
	IsHost, IsPause bool
}

type VTLobby struct {
	Name, Password, VideoLs string
	Viewers []VTViewer
	MaxOffset, VideoIndex int
	LastUpdateTime time.Time
}

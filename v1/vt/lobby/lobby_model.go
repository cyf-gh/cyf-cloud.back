package vt_lobby

import (
	"time"
)

type VTViewer struct {
	Name, Location string
	IsHost, IsPause bool
}

type VTLobby struct {
	Name, Password, VideoLs, Md5, IsPause, Location string
	Viewers []VTViewer
	MaxOffset, VideoIndex int
	LastUpdateTime time.Time
}

type VTVideoDesc struct {
	Ls string
	Index int
	Md5 string
}

package vt_lobby

import (
	"time"
)

type VTViewer struct {
	Name, Location string
	IsHost, IsPause bool
}

type VTLobby struct {
	Name, Password, VideoUrl, Cookie string
	IsShareCookie bool
	Viewers []VTViewer
	MaxOffset int
	LastUpdateTime time.Time
}

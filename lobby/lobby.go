package vt_lobby

type VTViewer struct {
	Name, Location string
	IsHost, IsPause bool
}

type VTLobby struct {
	Name, Password, VideoUrl string
	IsShareCookie bool
	Viewers []VTViewer
	MaxOffset int
}

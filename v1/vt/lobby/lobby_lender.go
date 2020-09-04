package vt_lobby

import (
	"time"
)
var Lobbies = []*VTLobby{}

/// 已废弃
func FindLobbyByViewer( aa string, viewerName string, lobbies []*VTLobby ) ( *VTLobby, int ) {
	for i, lb := range lobbies {
		for _, v := range lb.Viewers {
			if v.Name == viewerName && aa == lb.Name {
				return lb, i
			}
		}
	}
	return nil, -1
}

func FindLobbyByHost( userName string, lobbies []*VTLobby ) ( *VTLobby, int, int ) {
	for i, lb := range lobbies {
		for hi, v := range lb.Viewers {
			if v.Name == userName {
				if v.IsHost {
					return lb, i, hi
				}
			}
		}
	}
	return nil, -1, -1
}

func FindLobbyByUser( userName string, lobbies []*VTLobby ) ( *VTLobby, int, bool ) {
	for i, lb := range lobbies {
		for _, v := range lb.Viewers {
			if v.Name == userName {
				return lb, i, v.IsHost
			}
		}
	}
	return nil, -1, false
}

func CheckUserStatus( userName string, lobbies []*VTLobby ) string {
	for _, lb := range lobbies {
		for _, v := range lb.Viewers {
			if v.Name == userName {
				if v.IsHost {
					return "HOST"
				} else {
					return "GUEST"
				}
			}
		}
	}
	return "IDLE"
}

func ExitLobby( userName string, lobbies []*VTLobby ) (string,[]*VTLobby)  {
	for i, lb := range lobbies {
		for _, v := range lb.Viewers {
			if v.Name == userName {
				if v.IsHost {
					lobbies = DeleteLobbyAt( lobbies, i )
					return "LOBBY_DELETED", lobbies
				} else {
					if DeleteViewerIn( userName, lb.Name, Lobbies ) {
						return "LOBBY_EXIT", lobbies
					}
				}
			}
		}
	}
	return "NO_SUCH_LOBBY", lobbies
}

func IsSameNameLobbyExist( name string, lobbies []*VTLobby ) bool {
	for _, lb := range lobbies  {
		if (*lb).Name == name {
			return true
		}
	}
	return false
}

func AskForWhoIsHostViewerIn( viewers []VTViewer ) *VTViewer {
	for i, viewer := range viewers  {
		if viewer.IsHost {
			return &viewers[i]
		}
	}
	return nil
}

func DeleteLobbyNamed( name string, lobbies []*VTLobby  ) (bool, []*VTLobby) {
	for i, lb := range lobbies  {
		if (*lb).Name == name {
			return true, DeleteLobbyAt( lobbies, i )
		}
	}
	return false, lobbies
}

func DeleteLobbyAt( lobbies []*VTLobby, i int ) []*VTLobby {
	return append(lobbies[:i], lobbies[i+1:]...)
}

func DeleteViewerIn( viewerName string, lobbyName string, lobbies []*VTLobby ) bool {
	for _, lb := range lobbies  {
		if (*lb).Name == lobbyName {
			for j, v := range (*lb).Viewers  {
				if v.Name == viewerName {
					(*lb).Viewers = append( (*lb).Viewers[:j], (*lb).Viewers[j+1:]...)
					return true
				}
			}
		}
	}
	return false
}

func IsLobbyExist( name string, lobbies []*VTLobby  ) ( bool, *VTLobby ) {
	for _, lb := range lobbies  {
		if (*lb).Name == name {
			return true, lb
		}
	}
	return false, nil
}

func IsViewerExist( name string, lobby VTLobby ) bool {
	for _, v := range lobby.Viewers  {
		if v.Name == name {
			return true
		}
	}
	return false
}

func ViewerToString( viewer VTViewer ) string {
	host := "GUEST"
	status := "Playing"
	if viewer.IsHost {
		host = "HOST"
	}
	if viewer.IsPause {
		status = "Pause"
	}
	return host + "\t" + viewer.Name + "\t" + viewer.Location + "\t" + status
}

func UpdateLobbyLastUsedTime( lobby *VTLobby ) {
	lobby.LastUpdateTime = time.Now()
}

func ClearDiscardLobby( lobbies []*VTLobby ) []*VTLobby {
	m, _ := time.ParseDuration("1m")
	for i, lb := range lobbies {
		lastUpdateTime := lb.LastUpdateTime
		if time.Now().After( lastUpdateTime.Add( m ) ) {
			lobbies = DeleteLobbyAt( lobbies, i )
		}
	}
	return lobbies
}
package vt_lobby

import (
	ypm_parse "GoYPM/parse"
	st_comn_err "stgogo/comn/err"
	"strconv"
	"time"
)

// contrast prototype
//   0      1         2            3           4 			5		  ...
// “$name,$password,$max_offset,$host_name,$video_url,$is_share_cookie”
//
func CreateNewLobbyByContrast( lobbyContrast string ) *VTLobby {
	lobbyData := ypm_parse.SplitParaments( lobbyContrast )
	name := lobbyData[0]
	password := lobbyData[1]
	offset, err := strconv.Atoi(lobbyData[2])
	hostName := lobbyData[3]
	videoUrl := lobbyData[4]
	cookieData := lobbyData[5]


	st_comn_err.Exsit( err )
	newLobby := &VTLobby{
		Name:      name,
		Password:  password,
		Viewers:   nil,
		MaxOffset: offset,
		VideoUrl: videoUrl,
		IsShareCookie: !( cookieData == "no" ),
		Cookie: cookieData,
		LastUpdateTime: time.Now(),
	}
	// add host viewer
	host := &VTViewer{
		Name:     hostName,
		Location: "00:00",
		IsHost:   true,
		IsPause:  false,
	}
	newLobby.Viewers = append(newLobby.Viewers, *host)
	return newLobby
}

func FindLobbyByViewer( viewerName string, lobbies []*VTLobby ) ( *VTLobby, int ) {
	for i, lb := range lobbies {
		for _, v := range lb.Viewers {
			if v.Name == viewerName {
				return lb, i
			}
		}
	}
	return nil, -1
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
	m, _ := time.ParseDuration("10m")
	for i, lb := range lobbies {
		lastUpdateTime := lb.LastUpdateTime
		if time.Now().After( lastUpdateTime.Add( m ) ) {
			lobbies = DeleteLobbyAt( lobbies, i )
		}
	}
	return lobbies
}
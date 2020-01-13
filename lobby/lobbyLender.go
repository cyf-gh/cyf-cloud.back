package vt_lobby

import (
	st_comn_err "stgogo/comn/err"
	"strconv"
	"strings"
)

// contrast prototype
//   0      1         2            3           4 			5		  ...
// “$name,$password,$max_offset,$host_name,$video_url,$is_share_cookie”
//
func CreateNewLobbyByContrast( lobbyContrast string ) *VTLobby {
	lobbyData := strings.Split( lobbyContrast, ",")
	name := lobbyData[0]
	password := lobbyData[1]
	offset, err := strconv.Atoi(lobbyData[2])
	hostName := lobbyData[3]
	videoUrl := lobbyData[4]
	isShareCookie := lobbyData[5]


	st_comn_err.Exsit( err )
	newLobby := &VTLobby{
		Name:      name,
		Password:  password,
		Viewers:   nil,
		MaxOffset: offset,
		VideoUrl: videoUrl,
		IsShareCookie: isShareCookie == "share",
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

func FindLobbyByViewer( viewerName string, lobbies []*VTLobby ) *VTLobby {
	for _, lb := range lobbies {
		for _, v := range lb.Viewers {
			if v.Name == viewerName {
				return lb
			}
		}
	}
	return nil
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

func DeleteLobbyNamed( name string, lobbies []*VTLobby  ) bool {
	for i, lb := range lobbies  {
		if (*lb).Name == name {
			lobbies = append(lobbies[:i], lobbies[i+1:]...)
			return true
		}
	}
	return false
}

func IsLobbyExist( name string, lobbies []*VTLobby  ) ( bool, *VTLobby) {
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
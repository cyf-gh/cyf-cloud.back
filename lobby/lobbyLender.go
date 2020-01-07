package vt_lobby

import (
	"github.com/kpango/glg"
	"strconv"
	"strings"
)

// contrast prototype
//   0      1         2            3           4 ...
// “$name,$password,$max_offset,$host_name,$viewer1_name,..”
//
func CreateNewLobbyByContrast( lobbyContrast string ) *VTLobby {
	lobbyData := strings.Split( lobbyContrast, ",")
	offset, err := strconv.Atoi(lobbyData[2])
	if err != nil {
		glg.Log( err )
	}
	newLobby := &VTLobby{
		Name:      lobbyData[0],
		Password:  lobbyData[1],
		Viewers:   nil,
		MaxOffset: offset,
	}
	// add host viewer
	host := &VTViewer{
		Name:     lobbyData[3],
		Location: "00:00",
		IsHost:   true,
		IsPause:  false,
	}
	newLobby.Viewers = append(newLobby.Viewers, *host)
	// add guest viewers
	for i := 4; i < len(lobbyData); i++  {
		viewer := &VTViewer{
			Name:     lobbyData[i],
			Location: "00:00",
			IsHost:   false,
			IsPause:  false,
		}
		newLobby.Viewers = append(newLobby.Viewers, *viewer)
	}
	return newLobby
}

// start sync
func StartLobby( lobby *VTLobby ) {

}

func AskForWhoIsHostViewerIn( viewers []VTViewer ) *VTViewer {
	for i, viewer := range viewers  {
		if viewer.IsHost {
			return &viewers[i]
		}
	}
	return nil
}
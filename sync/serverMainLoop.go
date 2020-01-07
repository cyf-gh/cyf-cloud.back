package vt_sync

import (
	vtConfig "../config"
	vtLobby "../lobby"
	stErr "stgogo/comn/err"
	"github.com/kpango/glg"
	"net"
	"strings"
)

var Lobbies = []*vtLobby.VTLobby{}

//
// main loop which listen to the tcp request permanently
//
func Run() {
	listener, err := net.Listen("tcp", vtConfig.VtTcpAddr)
	defer listener.Close()
	stErr.Exsit(err)
	for {
		conn, err := listener.Accept()
		if stErr.Exsit(err) { continue }

		glg.Info("client connected :" + conn.RemoteAddr().String())

		// process coroutine udp procession
		go func() {
			defer func() {
				glg.Info(" client disconnected : " + conn.RemoteAddr().String())
				conn.Close()
			}()

			bMsg := []byte{}
			conn.Read( bMsg )
			msg := string( bMsg )
			glg.Info( (conn).RemoteAddr().String() + "\tsays\t"+ msg )
			ProcRequest( msg )
		}()
	}
}

func ProcRequest( rawString string ) {
	reqTypeAndParams := strings.Split( rawString, "@")
	head := reqTypeAndParams[0]
	body := reqTypeAndParams[1]
	switch head {
	case "create_lobby":
		// add lobby to lobbies array
		newLobby := vtLobby.CreateNewLobbyByContrast( body )
		lobIndex := len(Lobbies)
		Lobbies = append(Lobbies, newLobby)
		// response ready signal
		// _, err := (*conn).Write([]byte("READY"))
		// stErr.Exsit(err)
		// start sync
		StartUsingLobby( 10, newLobby )
		// sync finished and delete lobby
		Lobbies = append(Lobbies[:lobIndex], Lobbies[lobIndex+1:]...)
		break
	case "join_lobby":
		lobNameAndViewerName := strings.Split( body, ",")
		lobName := lobNameAndViewerName[0]
		viewerName := lobNameAndViewerName[1]
		for _, lob:= range Lobbies {
			if lob.Name == lobName {
				lob.Viewers = append(lob.Viewers, vtLobby.VTViewer{
					Name:     viewerName,
					Location: "00:00",
					IsHost:   false,
					IsPause:  false,
				})
			}
		}

		break
	default:
		glg.Error("Unknown tcp request\t" + rawString)
		break
	}
}
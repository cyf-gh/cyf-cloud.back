package vt_sync

import (
	vtLobby "../lobby"
	ypm_parse "GoYPM/parse"
	"github.com/kpango/glg"
	"net"
	stErr "stgogo/comn/err"
	stMem "stgogo/mem"
)

var Lobbies = []*vtLobby.VTLobby{}

//
// main loop which listen to the tcp request permanently
//
func RunTcpServer( addr string ) {
	listener, err := net.Listen("tcp", addr)
	defer listener.Close()
	stErr.Exsit(err)
	for {
		conn, err := listener.Accept()
		glg.Info("client connected :" + conn.RemoteAddr().String())
		if stErr.Exsit(err) {
			continue
		}

		bMsg := make([]byte, 1024)
		msgLen, err := conn.Read(bMsg)
		stErr.Exsit(err)
		bbMsg := make([]byte, msgLen)
		stMem.CutAdapt( bMsg, bbMsg )
		msg := string(bbMsg)
		glg.Info( (conn).RemoteAddr().String() + "\tsays\t"+ msg )
		ProcRequest( msg, DoResponse, &conn )
		glg.Info(" client disconnected : " + conn.RemoteAddr().String())
		conn.Close()
	}
}


type CBResponse func( msg string, conn *net.Conn )

func DoResponse( msg string, conn *net.Conn ) {
	if conn == nil {
		return
	} else {
		glg.Log( "[TCP] " + msg )
		_, err := (*conn).Write([]byte(msg))
		stErr.Exsit(err)
	}
}

// tcp message request process function
func ProcRequest( rawString string, doResp CBResponse, conn *net.Conn ) {
	head, body ,err := ypm_parse.SplitHeadBody( rawString )
	stErr.Exsit(err)
	switch head {
	case "create_lobby":
		// add lobby to lobbies array
		newLobby := vtLobby.CreateNewLobbyByContrast( body )
		if vtLobby.IsSameNameLobbyExist( newLobby.Name, Lobbies ) {
			doResp("LOBBY_ALREADY_EXISTS", conn )
			glg.Info(newLobby.Name + " is already exist but someone still want to borrow one.")
			break
		}
		Lobbies = append(Lobbies, newLobby)
		// response ready signal
		doResp("OK", conn )
		glg.Info("lobby " + newLobby.Name + " created!")
		// sync finished and delete lobby
		break
	case "join_lobby":
		lobNameAndViewerName := ypm_parse.SplitParaments(body)
		if len(lobNameAndViewerName) != 3 {
			doResp("INVALID_PARAM", conn )
			return
		}
		lobName := lobNameAndViewerName[0]
		viewerName := lobNameAndViewerName[1]
		lobPswd := lobNameAndViewerName[2]
		for _, lob:= range Lobbies {
			if lob.Name == lobName {
				if lobPswd != lob.Password {
					doResp("PSWD_INCOR", conn )
					return
				}
				if vtLobby.IsViewerExist( viewerName, *lob ) {
					doResp("ALREADY_IN_LOBBY", conn )
					return
				}
				lob.Viewers = append(lob.Viewers, vtLobby.VTViewer{
					Name:     viewerName,
					Location: "00:00",
					IsHost:   false,
					IsPause:  false,
				})
				doResp(lob.VideoUrl, conn )
				return
			}
		}
		doResp("NO_SUCH_LOBBY", conn )
		return
	case "query_lobbies":
		lobbyNames := ""
		for _, lob := range Lobbies {
			 lobbyNames += lob.Name + ","
		}
		doResp( lobbyNames, conn )
		return
	case "delete_lobby":
		if !vtLobby.DeleteLobbyNamed( body, Lobbies ) {
			doResp("NO_SUCH_LOBBY", conn )
		} else {
			doResp("OK", conn )
		}
	case "get_lobby_viewers":
		lobbyName := body
		viewerStr := ""
		exists, lob := vtLobby.IsLobbyExist(lobbyName, Lobbies)
		if !exists {
			doResp( "NO_SUCH_LOBBY", conn )
		} else {
			for _, v := range lob.Viewers {
				viewerStr += v.Name + ","
			}
		}
		doResp( viewerStr, conn )
		return
	default:
		glg.Error("Unknown tcp request\t" + rawString)
		break
	}
}
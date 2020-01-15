package vt_sync

import (
	vtlobby "../lobby"
	ypm_parse "GoYPM/parse"
	"github.com/kpango/glg"
	"net"
	stError "stgogo/comn/err"
	stmath "stgogo/comn/math"
	"strconv"
	"strings"
	"time"
)

// initialize udp connect object
func RunUDPSyncServer( addr string, lobs []*vtlobby.VTLobby, freshInterval time.Duration ) {
	udpAddr, err := net.ResolveUDPAddr("udp",addr )
	stError.Exsit(err)
	glg.Info( "\tUDP Server started at\t" + addr)

	for {
		conn, err := net.ListenUDP("udp", udpAddr)
		stError.Exsit(err)
		StartUdpSync( conn, freshInterval, lobs )
	}
}

// main process of video sync
func StartUdpSync( conn *net.UDPConn, freshInterval time.Duration, lobs []*vtlobby.VTLobby ) {
	/*
	glg.Info(
		lby.Name +
			"\nLobby started with viewers count:\t" +
			strconv.Itoa( len( lby.Viewers ) ) +
			"\nMax Location Offset:\t" +
			strconv.Itoa( lby.MaxOffset ) )
	*/
	for {
		time.Sleep(freshInterval)

		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buf)
		if stError.Exsit(err) { return }

		recvUdpMsg := string(buf[:n])
		glg.Log( recvUdpMsg )

		head, body, err := ypm_parse.SplitHeadBody( recvUdpMsg )

		if stError.Exsit(err) {
			_, err = conn.WriteToUDP( []byte( CheckLocationAndReturn(recvUdpMsg, lobs ) ), addr)
			if stError.Exsit(err) { continue }
		} else {
			switch head {
			// TODO: process operations
			case "get_lobby_viewers":
				lobbyName := body
				viewerStr := ""
				exists, lob := vtlobby.IsLobbyExist(lobbyName, Lobbies)
				if !exists {
					_, err = conn.WriteToUDP( []byte( "NO_SUCH_LOBBY"), addr )
				} else {
					for _, v := range lob.Viewers {
						viewerStr += v.Name + ","
					}
				}
				_, err = conn.WriteToUDP( []byte( "NO_SUCH_LOBBY"), addr )
				return
			default:
				glg.Error("Unknown tcp request\t" + recvUdpMsg)
				break
			}
		}
	}
}

func CheckLocationAndReturn( udpMsg string, lobs []*vtlobby.VTLobby ) string {
	curName, curLocation, isPause := splitNameAndLocationAndIsPauseFlag( udpMsg )
	lby := vtlobby.FindLobbyByViewer( curName, lobs )
	if lby == nil {
		glg.Error("No viewer called " + curName)
		 return "NO SUCH GUEST"
	}
	pHostViewer := vtlobby.AskForWhoIsHostViewerIn(lby.Viewers)

	for i, viewer := range lby.Viewers {
		if viewer.Name == curName {
			// update location and pause info first
			lby.Viewers[i].Location = curLocation
			lby.Viewers[i].IsPause = isPause == "p"
			// host viewer is always OK
			if viewer.IsHost { return "OK" }

			if pHostViewer.IsPause != lby.Viewers[i].IsPause {
				if pHostViewer.IsPause { return "p" } else { return "s" }
			}
			hm, hs := splitMinusAndSecond( pHostViewer.Location )
			cm, cs := splitMinusAndSecond( curLocation )
			offset := ( hm * 60 + hs ) - ( cm * 60 + cs )
			offset = stmath.Abs( offset )

			if offset >= lby.MaxOffset  {
				return pHostViewer.Location
			}
		}
	}
	return "OK"
}



func splitMinusAndSecond( currentTime string ) (int, int) {
	mAnds := strings.Split( currentTime, ":" )
	m, _ := strconv.Atoi( mAnds[0] )
	s, _ := strconv.Atoi( mAnds[1] )
	return m, s
}

func splitNameAndLocationAndIsPauseFlag( udpMsg string ) ( string, string, string ) {
	nAl := strings.Split(udpMsg, "," )
	if len(nAl) == 1 {
		glg.Error("Invalid UDP Message")
		return "", "", ""
	}
	return nAl[0], nAl[1], nAl[2]
}



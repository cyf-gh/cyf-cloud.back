package vt_sync

import (
	vtlobby "../lobby"
	stError "stgogo/comn/err"
	"github.com/kpango/glg"
	"net"
	stmath "stgogo/comn/math"
	"strconv"
	"strings"
	"time"
)

// initialize udp connect object
func Init() (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8800")
	if err != nil {
		glg.Error("ResolveUDPAddr err:", err)
		return nil, err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		glg.Error("ListenUDP err:", err)
		return nil, err
	}
	return conn, err
}

// main process of video sync
func StartUsingLobby(freshInterval time.Duration, lby *vtlobby.VTLobby) {
	glg.Info(
		lby.Name +
			"\nLobby started with viewers count:" +
			strconv.Itoa( len( lby.Viewers ) ) +
			"\nMax Location Offset: " +
			strconv.Itoa( lby.MaxOffset ) )

	conn, _ := Init()
	defer conn.Close()

	for {
		time.Sleep(freshInterval)

		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buf)
		if stError.Exsit(err) { return }

		recvUdpMsg := string(buf[:n])
		glg.Info("[Client]", recvUdpMsg)

		_, err = conn.WriteToUDP( []byte( CheckLocationAndReturn(recvUdpMsg, lby) ), addr,)
		if stError.Exsit(err) { return }
	}
}

func CheckLocationAndReturn( udpMsg string, lby *vtlobby.VTLobby ) string {
	curName, curLocation, isPause := splitNameAndLocationAndIsPauseFlag( udpMsg )
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



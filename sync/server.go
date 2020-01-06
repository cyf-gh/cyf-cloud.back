package vt_sync

import (
	vtlobby "../lobby"
	"net"
	stlog "stgogo/comn/log"
	stmath "stgogo/comn/math"
	"strconv"
	"strings"
	"time"
)

// initialize udp connect object
func Init() (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8800")
	if err != nil {
		stlog.Error.Println("ResolveUDPAddr err:", err)
		return nil, err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		stlog.Error.Println("ListenUDP err:", err)
		return nil, err
	}
	return conn, err
}

// main process of video sync
func Proc(conn *net.UDPConn, duration time.Duration, lby *vtlobby.VTLobby) {
	for {
		time.Sleep(duration)

		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			stlog.Error.Println( err )
			return
		}
		recvUdpMsg := string(buf[:n])
		stlog.Info.Println("[Client]", recvUdpMsg)

		_, err = conn.WriteToUDP( []byte( checkLocationAndReturn(recvUdpMsg, lby) ), addr,)
		if err != nil {
			stlog.Error.Println( err )
			return
		}
	}
}

func checkLocationAndReturn( udpMsg string, lby *vtlobby.VTLobby ) string {
	curName, curLocation, isPause := splitNameAndLocationAndIsPauseFlag( udpMsg )
	pHostViewer := vtlobby.AskForHostViewer(lby.Viewers)

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
	if len(nameAndLocation) == 1 {
		stlog.Error.Println("Invalid UDP Message")
		return "", "", ""
	}
	return nAl[0], nAl[1], nAl[2]
}
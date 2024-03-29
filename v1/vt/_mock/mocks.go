package vt_test_mock

import (
	vtlobby "../lobby"
	"time"
)

func GetMockLobby() *vtlobby.VTLobby {
	viewer_cyf := &vtlobby.VTViewer{
		Name:     "cyf",
		Location: "3:00",
		IsHost:   false,
		IsPause:  false,
	}

	viewer_yj := &vtlobby.VTViewer{
		Name:     "yj",
		Location: "4:00",
		IsHost:   true,
		IsPause:  false,
	}

	viewer_b1 := &vtlobby.VTViewer{
		Name:     "b1",
		Location: "4:00",
		IsHost:   false,
		IsPause:  false,
	}

	viewer_b2 := &vtlobby.VTViewer{
		Name:     "b2",
		Location: "4:00",
		IsHost:   false,
		IsPause:  false,
	}
	viewer_b3 := &vtlobby.VTViewer{
		Name:     "b3",
		Location: "5:00",
		IsHost:   false,
		IsPause:  false,
	}
	ff := "2006-01-02 15:04:05"
	t, _ := time.Parse(ff, "2019-03-10 11:00:00")
	mockLobby := &vtlobby.VTLobby{
		Name:          "g",
		Password:      "123",
		VideoUrl:      "https://bilibili.com/av123456",
		IsShareCookie: false,
		Viewers:       []vtlobby.VTViewer{},
		MaxOffset:     3,
		LastUpdateTime: t,
	}
	mockLobby.Viewers = append(mockLobby.Viewers, *viewer_yj, *viewer_cyf, *viewer_b1, *viewer_b2, *viewer_b3)

	return mockLobby
}

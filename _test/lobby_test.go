package vt_test

import (
	vtmocks "../_mock"
	vtLobby "../lobby"
	vtsync "../sync"
	"testing"
)

func TestCreateLobbyFromContrast(t *testing.T) {
	lob := vtLobby.CreateNewLobbyByContrast("new_lobby,123123,2,yj,https://baidu.com,share")
	if lob.Name != "new_lobby" || lob.Password != "123123" || lob.Viewers[0].Name != "yj" {
		t.Fatal()
	}
}

func TestFindLobbyByViewer( t *testing.T ) {
	vtsync.Lobbies = append(vtsync.Lobbies, vtmocks.GetMockLobby())

	if !( vtLobby.FindLobbyByViewer( "cyf", vtsync.Lobbies ).Name == "g" ) {
		t.Fatal()
	}
}
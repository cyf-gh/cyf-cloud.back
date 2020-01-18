package vt_test

import (
	vtmocks "../_mock"
	vtLobby "../lobby"
	vtsync "../sync"
	"testing"
)

func Test_CreateLobbyFromContrast(t *testing.T) {
	lob := vtLobby.CreateNewLobbyByContrast("new_lobby,123123,2,yj,https://baidu.com,share")
	if lob.Name != "new_lobby" || lob.Password != "123123" || lob.Viewers[0].Name != "yj" {
		t.Fatal()
	}
}

func Test_FindLobbyByViewer( t *testing.T ) {
	vtsync.Lobbies = append(vtsync.Lobbies, vtmocks.GetMockLobby())

	if !( vtLobby.FindLobbyByViewer( "cyf", vtsync.Lobbies ).Name == "g" ) {
		t.Fatal()
	}
}

func Test_DeleteViewerInLobby( t *testing.T ) {
	vtsync.Lobbies = append(vtsync.Lobbies, vtmocks.GetMockLobby())

	if !(vtLobby.DeleteViewerIn("cyf", "g", vtsync.Lobbies)) {
		t.Fatal()
	}
	if vtLobby.IsViewerExist("cyf", *vtsync.Lobbies[0] ) {
		t.Fatal()
	}
}
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

func Test_DeleteOneSlice(t *testing.T) {
	ints := []int { 1 }
	ints = append(ints[:0], ints[0+1:]...)
	if len(ints) != 0 {
		t.Fatal()
	}
}

func Test_ClearDiscardLobby( t*testing.T ) {
	lobs := []*vtLobby.VTLobby{ vtmocks.GetMockLobby() }
	lobs = vtLobby.ClearDiscardLobby( lobs )
	if !( len(lobs) == 0 ) {
		t.Log(lobs[0])
		t.Fatal()
	}

	// m, _ := time.ParseDuration("10m")
	// t.Log( time.Now().After( lobs[0].LastUpdateTime.Add( m ) ) )
}

func Test_DeleteLobbyByName( t*testing.T ) {
	lobs := []*vtLobby.VTLobby{ vtmocks.GetMockLobby() }
	_, lobs = vtLobby.DeleteLobbyNamed( "g", lobs )
	if !( len(lobs)==0 ) {
		t.Log(lobs[0])
		t.Fatal()
	}
}
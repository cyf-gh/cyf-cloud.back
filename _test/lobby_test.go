package vt_test

import (
	vtLobby "../lobby"
	"testing"
)

func TestCreateLobbyFromContrast(t *testing.T) {
	lob := vtLobby.CreateNewLobbyByContrast("new_lobby,123123,2,yj,cyf,c1,c2")
	if lob.Name != "new_lobby" || lob.Password != "123123" || len(lob.Viewers) != 4 || lob.Viewers[0].Name != "yj" {
		t.Fail()
	}
	t.Log("ok\tcreate lobby from contrast")
}
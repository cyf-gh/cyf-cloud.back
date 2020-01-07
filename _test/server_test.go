package vt_test

import (
	vtsync "../sync"
	"testing"
	vtmocks "../_mock"
)

func TestJoinLobby(t *testing.T) {
	vtsync.Lobbies = append(vtsync.Lobbies, vtmocks.GetMockLobby())
	vtsync.ProcRequest("join_lobby@g,cyf1")

	for _, lob:= range vtsync.Lobbies {
		for _, v := range lob.Viewers {
			if v.Name == "cyf1" {
				t.Log("ok\tjoin_lobby")
				t.SkipNow()
			}
		}
	}
	t.Fail()
}


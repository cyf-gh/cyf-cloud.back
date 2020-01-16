package vt_test

import (
	vtsync "../sync"
	"testing"
)

func TestJoinLobby(t *testing.T) {
	vtsync.ProcRequest("join_lobby@g,cyf1,123", vtsync.DoResponse, nil )
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


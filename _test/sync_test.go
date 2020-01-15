package vt_test

import (
	vtTestMock "../_mock"
	vtsync "../sync"
	vtlobby "../lobby"
	"testing"
)

func BenchmarkSsync(b *testing.B) {
	lobs := []*vtlobby.VTLobby{ vtTestMock.GetMockLobby() }
	b.Log("normal play benchmark")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i % 2 == 1 {
			vtsync.CheckLocationAndReturn( "cyf,00:03,s", lobs )
			vtsync.CheckLocationAndReturn( "yj,00:03,s", lobs )
		} else {
			vtsync.CheckLocationAndReturn( "cyf,00:02,s", lobs )
			vtsync.CheckLocationAndReturn( "yj,00:02,s", lobs )
		}
	}
}

func TestCheckLocationAndReturn(t *testing.T) {
	lobs := []*vtlobby.VTLobby{ vtTestMock.GetMockLobby() }
	vtsync.CheckLocationAndReturn( "yj,00:03,s", lobs )

	if vtsync.CheckLocationAndReturn( "cyf,00:34,s", lobs ) != "00:03" {
		t.Error()
	}
	t.Log("ok\ttime adjust")

	if vtsync.CheckLocationAndReturn( "cyf,00:03,s", lobs ) != "OK" {
		t.Error()
	}
	t.Log("ok\tsame time and return OK")

	if vtsync.CheckLocationAndReturn( "cyf,00:04,s", lobs ) != "OK" {
		t.Error()
	}
	t.Log("ok\tin offset time and return OK")

	vtsync.CheckLocationAndReturn( "yj,00:03,p", lobs )

	if vtsync.CheckLocationAndReturn( "cyf,00:03,s", lobs ) != "p" {
		t.Error()
	}
	t.Log("ok\tpause")

	vtsync.CheckLocationAndReturn( "yj,00:03,s", lobs )

	if vtsync.CheckLocationAndReturn( "cyf,00:03,p", lobs ) != "s" {
		t.Error()
	}
	t.Log("ok\tstart")

	/*
	for _, viewer := range mlob.Viewers {
		if viewer.Name == "cyf" {
			if viewer.Location != "00:03" {
				t.Error()
			}
		}
	}
	*/
}
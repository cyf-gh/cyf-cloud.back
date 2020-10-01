package vt_test

import (
	vtTestMock "../_mock"
	vtlobby "../lobby"
	vtsync "../sync"
	"testing"
)

func Benchmark_Ssync(b *testing.B) {
	lobs := []*vtlobby.VTLobby{ vtTestMock.GetMockLobby() }
	b.Log("normal play benchmark")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i % 2 == 1 {
			vtsync.CheckLocationAndReturn( "cyf,00:03,s" )
			vtsync.CheckLocationAndReturn( "yj,00:03,s" )
		} else {
			vtsync.CheckLocationAndReturn( "cyf,00:02,s" )
			vtsync.CheckLocationAndReturn( "yj,00:02,s" )
		}
	}
}

func Test_CheckLocationAndReturn(t *testing.T) {
	vtsync.CheckLocationAndReturn( "yj,00:03,s" )

	if vtsync.CheckLocationAndReturn( "cyf,00:34,s" ) != "00:03" {
		t.Error()
	}
	t.Log("ok\ttime adjust")

	if vtsync.CheckLocationAndReturn( "cyf,00:03,s" ) != "OK" {
		t.Error()
	}
	t.Log("ok\tsame time and return OK")

	if vtsync.CheckLocationAndReturn( "cyf,00:04,s" ) != "OK" {
		t.Error()
	}
	t.Log("ok\tin offset time and return OK")

	vtsync.CheckLocationAndReturn( "yj,00:03,p" )

	if vtsync.CheckLocationAndReturn( "cyf,00:03,s" ) != "p" {
		t.Error()
	}
	t.Log("ok\tpause")

	vtsync.CheckLocationAndReturn( "yj,00:03,s" )

	if vtsync.CheckLocationAndReturn( "cyf,00:03,p" ) != "s" {
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
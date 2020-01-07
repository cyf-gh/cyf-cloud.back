package vt_test

import (
	vtTestMock "../_mock"
	vtsync "../sync"
	"testing"
)

func BenchmarkSsync(b *testing.B) {
	mlob := vtTestMock.GetMockLobby()
	b.Log("normal play benchmark")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i % 2 == 1 {
			vtsync.CheckLocationAndReturn( "cyf,00:03,s", mlob )
			vtsync.CheckLocationAndReturn( "yj,00:03,s", mlob )
		} else {
			vtsync.CheckLocationAndReturn( "cyf,00:02,s", mlob )
			vtsync.CheckLocationAndReturn( "yj,00:02,s", mlob )
		}
	}
}

func TestCheckLocationAndReturn(t *testing.T) {
	mlob := vtTestMock.GetMockLobby()
	vtsync.CheckLocationAndReturn( "yj,00:03,s", mlob )

	if vtsync.CheckLocationAndReturn( "cyf,00:34,s", mlob ) != "00:03" {
		t.Error()
	}
	t.Log("ok\ttime adjust")

	if vtsync.CheckLocationAndReturn( "cyf,00:03,s", mlob ) != "OK" {
		t.Error()
	}
	t.Log("ok\tsame time and return OK")

	if vtsync.CheckLocationAndReturn( "cyf,00:04,s", mlob ) != "OK" {
		t.Error()
	}
	t.Log("ok\tin offset time and return OK")

	vtsync.CheckLocationAndReturn( "yj,00:03,p", mlob )

	if vtsync.CheckLocationAndReturn( "cyf,00:03,s", mlob ) != "p" {
		t.Error()
	}
	t.Log("ok\tpause")

	vtsync.CheckLocationAndReturn( "yj,00:03,s", mlob )

	if vtsync.CheckLocationAndReturn( "cyf,00:03,p", mlob ) != "s" {
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
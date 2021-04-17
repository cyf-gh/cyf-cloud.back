package cc

import (
	"testing"
)

func TestGet(t *testing.T) {
	_, e := Get( "https://baidu.com" )
	t.Log( e )
}

func TestGetByProxy(t *testing.T) {
	r, e := GetByProxy( "https://google.com", "http://192.168.0.103:1080" )
	t.Log( r, e )
}


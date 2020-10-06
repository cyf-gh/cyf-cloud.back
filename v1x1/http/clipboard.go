package http

import (
	"net/http"
	err "../err"
	cache "../cache"
)

func MakeClipboardKey( r *http.Request ) (string, error) {
	atk, e := GetAtk( r )
	return "clipboard" + atk, e;
}

func ClipboardPush(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	key, e := MakeClipboardKey( r ); err.Check( e )
	bs, e := Body2String( r ); err.Check( e )
	_, e = cache.Set( key, bs ); err.Check( e )
	err.HttpReturnOk( &w )
}

func ClipboardFetch(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	key, e := MakeClipboardKey( r ); err.Check( e )
	v, e := cache.Get( key ); err.Check( e )
	err.HttpReturnOkWithData( &w, v )
}
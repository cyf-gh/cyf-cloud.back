package http

import (
	"../../cc"
	err "../../cc/err"
	cache "../cache"
	"net/http"
)

func MakeClipboardKey( r *http.Request ) (string, error) {
	a, e := GetAccountByAtk( r )
	return "$clipboard$" + a.Name, e
}
func init() {
	cc.AddActionGroup( "/v1x1/clipboard", func( a cc.ActionGroup ) error {
	    a.POST( "/push", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			key, e := MakeClipboardKey( ap.R ); err.Check( e )
			bs, e := Body2String( ap.R ); err.Check( e )
			_, e = cache.Set( key, bs ); err.Check( e )
			return cc.HerOk()
	    } )

	    a.GET( "/fetch", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			key, e := MakeClipboardKey( ap.R ); err.Check( e )
			v, e := cache.Get( key ); err.Check( e )
			return cc.HerOkWithString( v )
	    } )

	    return nil
	} )
}
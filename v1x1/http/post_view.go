package http

import (
	"../../cc"
	"../../cc/err"
	"../cache"
	"stgogo/comn/convert"
)

func init() {
	cc.AddActionGroup( "/v1x1/post/view", func( a cc.ActionGroup ) error {
		a.GET( "/count",  func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var pid string
			if pid = ap.R.FormValue("id"); pid == "" {
				return cc.HerArgInvalid( "id ")
			}
			return cc.HerOkWithString( convert.I64toa( getPostView( pid ) ) )
		} )
		a.GET( "",  func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var pid string

			if pid = ap.R.FormValue("id"); pid == "" {
				return cc.HerArgInvalid( "id ")
			}
			e := doPostView( pid ); err.Check( e )
			return cc.HerOk()
		} )
		return nil

	} )
}

func getPostView( pid string ) int64 {
	k :=  _postViewPrefix + pid
	if ps, e := cache.Get( k ); e != nil {
		// 还没有键，创建键
		cache.Set( k, "0")
		return 0
	} else {
		p, _ := convert.Atoi64( ps )
		return p
	}
}

func doPostView( pid string ) error {
	k :=  _postViewPrefix + pid

	_, e := cache.Set( k, convert.I64toa( getPostView( pid ) + 1) ) // no error
	return e
}

package http

import (
	"../../cc"
	err "../../cc/err"
	"../orm"
)

func init() {
	cc.AddActionGroup( "/v1x1/search", func( a cc.ActionGroup ) error {

		a.GET( "/user", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var ( e error )
			text := ap.R.FormValue("a")
			res, e := orm.VagueSearchAccountName(text); err.Check( e )
			return cc.HerOkWithData( res )
		} )

		a.GET( "/post", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var ( e error )
			text := ap.R.FormValue("a")
			ps, tags, e := orm.VagueSearchPostAndTagName(text); err.Check( e )
			eps, e := extendPostInfo( ps ); err.Check( e )
			return cc.HerOkWithData( cc.H {
				"PostInfos": eps,
				"Tags": tags,
			} )
		} )

	    return nil
	} )
}
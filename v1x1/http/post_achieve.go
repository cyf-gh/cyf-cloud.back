package http

import (
	"../../cc"
	"../../cc/err"
	"../orm"
)

func init() {
	cc.AddActionGroup( "/v1x1/posts/achieve", func( a cc.ActionGroup ) error {

		a.GET( "/date", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			au := ap.GetFormValue("au")
			a, e := orm.GetAccountByName( au ); err.Check( e )
			cd, e := orm.GetOnesAllPostInfoDate( a.Id ); err.Check( e )
			return cc.HerOkWithData( cd )
		} )

		a.GET( "/tag", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			au := ap.GetFormValue("au")
			a, e := orm.GetAccountByName( au ); err.Check( e )
			tags, e := orm.GetOnesAllPostInfoTags( a.Id )
			return cc.HerOkWithData( tags )
		} )

		a.GET( "/recent", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			au := ap.GetFormValue("au")
			a, e := orm.GetAccountByName( au ); err.Check( e )
			titles, e := orm.GetOnesRecentPostTitle( a.Id ); err.Check( e )
			return cc.HerOkWithData( titles )
		} )

		return nil
	} )
}
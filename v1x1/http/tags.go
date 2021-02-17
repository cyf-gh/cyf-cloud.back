package http

import (
	"../../cc"
	"../../cc/err"
	"../orm"
)

func init() {
	cc.AddActionGroup( "/v1x1/tags", func( a cc.ActionGroup ) error {

		a.GET( "", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			tags, e := orm.GetAllTags(); err.Assert( e )
			return cc.HerOkWithData( tags )
		} )

	    return nil
	} )
}
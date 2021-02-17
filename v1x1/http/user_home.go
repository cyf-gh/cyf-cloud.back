// 个人主页API
package http

import (
	"../../cc"
	"../../cc/err"
	"../orm"
	"stgogo/comn/convert"
)

func init() {
	cc.AddActionGroup( "/v1x1/account/info", func( a cc.ActionGroup ) error {

		// 返回的infocomponentmodel一定不为空，如果不存在则为""
	    a.GET( "/home", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			strid := ap.R.FormValue("uid")
			id, e := convert.Atoi64( strid ); err.Assert( e )
			hi, e := orm.GetUserHomeInfo( id ); err.Assert( e )
	        return cc.HerOkWithData( hi )
	    } )
	    return nil
	} )
}

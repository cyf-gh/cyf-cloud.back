package http

import (
	"../../cc"
	"../../cc/err"
	"../orm"
	"stgogo/comn/convert"
)

func init() {
	cc.AddActionGroupDeprecated( "/v1x1/dm/1/indicate", func( a cc.ActionGroup ) error {
		// \brief 用于指定某资源为二进制资源，如果某
		// \arg[id] 资源的id
		a.GET( "/binary", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			id, e := convert.Atoi64( ap.GetFormValue( "id" ) )
			e = orm.DMIndicateResourceBinary( id ); err.Assert( e )
			return cc.HerOk()
		} )
		// \brief 用于指定某资源为二进制资源，如果某
		// \arg[rid] 资源的id
		// \arg[bid] 备份资源的id，为TargetResource的Id
		a.GET( "/backup", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			rid, e := convert.Atoi64( ap.GetFormValue( "rid" ) )
			bid, e := convert.Atoi64( ap.GetFormValue("bid" ) )
			e = orm.DMIndicateResourceBackupOf( rid, bid ); err.Assert( e )
			return cc.HerOk()
		} )
		return nil
	} )
}

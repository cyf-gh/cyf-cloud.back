package http

import (
	"../../cc"
	"../../cc/err"
	"../orm"
)

func init() {
	cc.AddActionGroupDeprecated( "/v1x1/dm/1/info", func( a cc.ActionGroup ) error {
		// \brief 修改单个资源的信息
		// \body orm.DMTargetResource
		// \return ok
		a.POST( "/modifies", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			tr := orm.DMTargetResource{}
			e = ap.GetBodyUnmarshal( &tr ); err.Assert( e )
			tr.Update()
			return cc.HerOk()
		})
		// \brief 修改多个资源的信息
		// \body orm.DMTargetResource
		// \return ok
		a.POST( "/modify", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			var tr []orm.DMTargetResource
			e = ap.GetBodyUnmarshal( &tr ); err.Assert( e )
			for _, t := range tr {
				t.Update()
			}
			return cc.HerOk()
		})
		// \brief 修改单个资源的ex信息
		// \body orm.DMTargetResource
		// \return ok
		a.POST( "/ex/modifies", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			var tr orm.DMTargetResourceEx
			e = ap.GetBodyUnmarshal( &tr ); err.Assert( e )
			tr.Update()
			return cc.HerOk()
		})
		// \brief 修改多个资源的ex信息
		// \body orm.DMTargetResourceEx
		// \return ok
		a.POST( "/ex/modify", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			var tr []orm.DMTargetResourceEx
			e = ap.GetBodyUnmarshal( &tr ); err.Assert( e )
			for _, t := range tr {
				t.Update()
			}
			return cc.HerOk()
		})
		return nil
	} )
}
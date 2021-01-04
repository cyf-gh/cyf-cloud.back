package http

import (
	"../../cc"
	"../../cc/err"
	"../orm"
)

func init() {
	cc.AddActionGroup( "/v1x1/dm/1/info", func( a cc.ActionGroup ) error {
		// \brief 修改单个资源的信息
		// \body orm.DMTargetResource
		a.POST( "/modifies", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			tr := orm.DMTargetResource{}
			e = ap.GetBodyUnmarshal( &tr ); err.Check( e )
			tr.Update()
			return cc.HerOk()
		})
		// \brief 修改多个资源的信息
		// \body orm.DMTargetResource
		a.POST( "/modify", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			var tr []orm.DMTargetResource
			e = ap.GetBodyUnmarshal( &tr ); err.Check( e )
			for _, t := range tr {
				t.Update()
			}
			return cc.HerOk()
		})
		// \brief 修改单个资源的ex信息
		// \body orm.DMTargetResource
		a.POST( "/ex/modifies", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			var tr orm.DMTargetResourceEx
			e = ap.GetBodyUnmarshal( &tr ); err.Check( e )
			tr.Update()
			return cc.HerOk()
		})
		// \brief 修改多个资源的ex信息
		// \body orm.DMTargetResourceEx
		a.POST( "/ex/modify", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			var tr []orm.DMTargetResourceEx
			e = ap.GetBodyUnmarshal( &tr ); err.Check( e )
			for _, t := range tr {
				t.Update()
			}
			return cc.HerOk()
		})
		return nil
	} )
}
package http

import (
	"../../cc"
	"../../cc/err"
	"../orm"
	"time"
)

func init() {
	cc.AddActionGroupDeprecated( "/v1x1/dm/1/query", func( a cc.ActionGroup ) error {
		// \brief 获取某个资源的所有克隆资源
		// \body orm.DMTargetResource
		// \return []DMTargetResource
		a.POST( "/clones", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			dmRes := &orm.DMTargetResource{}
			e = ap.GetBodyUnmarshal( dmRes ); err.Assert( e )
			start := time.Now()
			clones, e := dmRes.GetClones(); err.Assert( e )
			return cc.HerOkWithDataAndUsedTime( clones, time.Since( start ) )
		} )
		// \brief 获取某个目录下的所有目标资源
		// \arg[d] 目录
		// \return []DMTargetResource
		a.GET( "/dir", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Assert( e )
			dir := ap.GetFormValue( "d" )
			dmDir, e := checkDir( dir ); err.Assert( e )
			lsRes, e := dmDir.Ls(); err.Assert( e )

			var trs []orm.DMTargetResource
			for _, r := range lsRes {
				tr, e := orm.DMGetTargetResourceByPath( r.Path ); err.Assert( e )
				if tr != nil {
					trs = append(trs, *tr)
				}
			}
			return cc.HerOkWithData( trs )
		} )
		return nil
	} )
}
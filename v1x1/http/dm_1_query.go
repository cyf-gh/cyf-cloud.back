package http

import (
	"../../cc"
	"../../cc/err"
	"../orm"
	"time"
)

func init() {
	cc.AddActionGroup( "/v1x1/dm/1/query", func( a cc.ActionGroup ) error {
		// \brief 获取某个资源的所有克隆资源
		// \body orm.DMTargetResource
		// \return []DMTargetResource
		a.POST( "/clones", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			dmRes := &orm.DMTargetResource{}
			e = ap.GetBodyUnmarshal( dmRes ); err.Check( e )
			start := time.Now()
			clones, e := dmRes.GetClones(); err.Check( e )
			return cc.HerOkWithDataAndUsedTime( clones, time.Since( start ) )
		} )
		// \brief 获取某个目录下的所有目标资源
		// \arg[d] 目录
		// \return []DMTargetResource
		a.GET( "/dir", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			dir := ap.GetFormValue( "d" )
			dmDir, e := checkDir( dir ); err.Check( e )
			lsRes, e := dmDir.Ls(); err.Check( e )

			var trs []orm.DMTargetResource
			for _, r := range lsRes {
				tr, e := orm.DMGetTargetResourceByPath( r.Path ); err.Check( e )
				if tr != nil {
					trs = append(trs, *tr)
				}
			}
			return cc.HerOkWithData( trs )
		} )
		return nil
	} )
}
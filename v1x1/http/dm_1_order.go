// /order所有操作都会添加数据至数据库
package http

import (
	"../../cc"
	"../../cc/err"
	"../../dm_1"
	"../orm"
	"time"
)

func init() {
	cc.AddActionGroup( "/v1x1/dm/1/order", func( a cc.ActionGroup ) error {
		// \brief 开始递归所有目录进行资源索引
		// \return 递归所用的时间
		a.GET( "/recruit", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			dmRootDir := &dm_1.DMResource{
				Path: dm_1.DMRootPath(),
			}
			start := time.Now()
			lsRootRes := dmRootDir.LsRecruit()
			usedTime := time.Since( start )
			e = orm.DMAddResource( lsRootRes ); err.Check( e )
			return cc.HerOkWithData( usedTime )
		} )
		// \brief 添加某个目录下的所有资源
		// \return ok
		a.GET( "/ls", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			dir := ap.GetFormValue( "d" )
			dmDir, e := checkDir( dir ); err.Check( e )
			lsRes, e := dmDir.Ls(); err.Check( e )
			e = orm.DMAddResource( lsRes ); err.Check( e )
			return cc.HerOk()
		} )
		// \brief 添加一个或多个资源
		// \return ok
		a.POST( "", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			var dmRes []dm_1.DMResource
			e = ap.GetBodyUnmarshal( &dmRes ); err.Check( e )
			e = orm.DMAddResource( dmRes ); err.Check( e )
			return cc.HerOk()
		} )
		return nil
	} )
}
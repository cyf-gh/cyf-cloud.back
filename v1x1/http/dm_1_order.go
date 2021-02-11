// /order所有操作都会添加数据至数据库
package http

import (
	"../../cc"
	"../../cc/err"
	"../../dm_1"
	"../orm"
)

func init() {
	cc.AddActionGroup( "/v1x1/dm/1/order", func( a cc.ActionGroup ) error {
		// \brief 开始递归所有目录进行资源索引
		// \arg[path] 要递归索引的目录
		// \note 会导致并发
		// \return ok
		a.GET( "/recruit", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )

			rootDir := ap.GetFormValue("d")
			if rootDir == "" {
				rootDir = dm_1.DMRootPath()
			}
			dmRootDir := &dm_1.DMResource{
				Path: rootDir,
			}
			go func() {
				if e := dm_1.TaskSharedList.AddTask( "order_recruit", true, 600, dmRootDir.LsRecruitCount() ); e != nil {
					return
				}
				orTask := &dm_1.TaskSharedList.Lists["order_recruit"][0]

				lsRootRes := dmRootDir.LsRecruit( orTask )
				e = orm.DMAddResources( lsRootRes, orTask ); orTask.Error( e )

				orTask.Finished()
			} ()
			return cc.HerOk()
		} )
		// \brief 添加某个目录下的所有资源
		// \return ok
		a.GET( "/ls", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			dir := ap.GetFormValue( "d" )
			dmDir, e := checkDir( dir ); err.Check( e )
			lsRes, e := dmDir.Ls(); err.Check( e )
			e = orm.DMAddResources( lsRes, nil ); err.Check( e )
			return cc.HerOk()
		} )
		// \brief 添加一个或多个资源
		// \return ok
		a.POST( "", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			var dmRes []dm_1.DMResource
			e = ap.GetBodyUnmarshal( &dmRes ); err.Check( e )
			e = orm.DMAddResources( dmRes, nil ); err.Check( e )
			return cc.HerOk()
		} )
		return nil
	} )
}
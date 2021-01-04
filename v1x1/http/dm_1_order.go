// /order所有操作都会添加数据至数据库
package http

import (
	"../../cc"
	"../../cc/err"
	"../../dm_1"
	"../orm"
	"errors"
	"strings"
)

func init() {
	cc.AddActionGroup( "/v1x1/dm/1/ls", func( a cc.ActionGroup ) error {
		// \brief 返回dm目录的资源，用于索引数据库
		// \arg[d] 路径，附加于root_path之后的路径
		a.GET( "", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			dir := ap.GetFormValue( "d" )
			if dir == "" {
				panic( errors.New("do nothing with empty dir") )
			}
			if strings.Contains( dir, "..") {
				panic( errors.New(".. is not allowed in param d ") )
			}
			dmDir := &dm_1.DMResource{ Path: dm_1.DMRootPath() + dir }
			if !dmDir.Exists() {
				panic( errors.New("specific path is invalid") )
			}
			if !dmDir.IsDire() {
				panic( errors.New("specific path is a file") )
			}
			lsRes, e := dmDir.Ls(); err.Check( e )
			e = orm.DMAddResource( lsRes ); err.Check( e )
			return cc.HerOkWithData( lsRes )
		} )
		// \brief 返回dm根目录的资源，用于索引数据库
		a.GET( "/root", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			dmRootDir := &dm_1.DMResource{
				Path: dm_1.DMRootPath(),
			}
			lsRootRes, e := dmRootDir.Ls(); err.Check( e )
			e = orm.DMAddResource( lsRootRes ); err.Check( e )
			return cc.HerOkWithData( lsRootRes )
		} )
		a.GET( "/recruit", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			dmRootDir := &dm_1.DMResource{
				Path: dm_1.DMRootPath(),
			}
			lsRootRes, e := dmRootDir.Ls(); err.Check( e )
			e = orm.DMAddResource( lsRootRes ); err.Check( e )
			return cc.HerOkWithData( lsRootRes )
		} )
		return nil
	} )
}
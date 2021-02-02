// 用于查看目录与文件
package http

import (
	"../../cc"
	"../../cc/err"
	"../../dm_1"
	"errors"
	"strings"
)

// 检查目录是否正常，并返回一个从根目录出发的目录
func checkDir( dir string ) ( dmDir *dm_1.DMResource, e error ) {
	if dir == "" {
		return nil, errors.New("do nothing with empty dir")
	}
	if strings.Contains( dir, "..") {
		return nil, errors.New(".. is not allowed in param d")
	}
	dmDir = &dm_1.DMResource{ Path: dir }
	if !dmDir.Exists() {
		return nil, errors.New("specific path is invalid")
	}
	if !dmDir.IsDire() {
		return nil, errors.New("specific path is a file")
	}
	return dmDir, nil
}

func init() {
	cc.AddActionGroup( "/v1x1/dm/1/raw", func( a cc.ActionGroup ) error {
		// \brief 返回dm根目录
		a.GET( "/root", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			return cc.HerOkWithData( dm_1.DMRootPath() )
		} )
		// \brief 返回dm目录的资源，用于索引数据库
		// \arg[d] 路径，附加于root_path之后的路径
		a.GET( "/dir", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			dir := ap.GetFormValue( "d" )
			dmDir, e := checkDir( dir ); err.Check( e )
			lsRes, e := dmDir.Ls(); err.Check( e )
			fivm := []dm_1.DMFileInfoViewModel{}
			for _, res := range lsRes {
				fivm = append(fivm, *res.ToReadable())
			}
			return cc.HerOkWithData( fivm )
		} )
		// \brief 返回dm根目录的资源，用于索引数据库
		a.GET( "/dir/root", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			dmRootDir := &dm_1.DMResource{
				Path: dm_1.DMRootPath(),
			}
			lsRootRes, e := dmRootDir.Ls(); err.Check( e )
			fivm := []dm_1.DMFileInfoViewModel{}
			for _, res := range lsRootRes {
				fivm = append(fivm, *res.ToReadable())
			}
			return cc.HerOkWithData( fivm )
		} )
		// \brief 返回dm目录的递归资源，用于索引数据库
		// \arg[d] 路径，附加于root_path之后的路径
		a.GET( "/recruit/dir", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			dir := ap.GetFormValue( "d" )
			dmDir, e := checkDir( dir ); err.Check( e )
			lsRes := dmDir.LsRecruit( nil ); err.Check( e )
			return cc.HerOkWithData( lsRes )
		} )
		// \brief 返回该目录的大小
		// \arg[d] 路径，附加于root_path之后的路径
		a.GET( "/recruit/dir/size", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			e := DM1CheckPermission( ap.R ); err.Check( e )
			dir := ap.GetFormValue( "d" )
			dmDir, e := checkDir( dir ); err.Check( e )
			size, e := dmDir.GetSize(); err.Check( e )
			return cc.HerOkWithData( size )
		} )
		return nil
	} )
}
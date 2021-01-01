package http

import (
	"../orm"
	"../../cc"
	"../../cc/err"
	"../../dm_1"
	"errors"
	"strings"
)

func init() {
	// \brief 返回dm目录的资源
	// \type GET
	// \arg[d] 路径，附加于root_path之后的路径
	cc.AddActionGroup( "/v1x1/dm/1/ls", func( a cc.ActionGroup ) error {
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
			lsRes, e := dmDir.Ls()
			return cc.HerOkWithData( lsRes )
		} )
		return nil
	} )
}
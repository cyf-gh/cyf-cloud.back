package http

import (
	"../../cc"
	"../../ccDoc"
	cfg "../../config"
)

func init() {
	cc.AddActionGroup( "/v1x1/meta", func( a cc.ActionGroup ) error {
		// \brief 返回所有v1x1的restful API与描述
		a.GET( "/api/ref", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			return cc.HerOkWithString( ccDoc.GenerateDocJson( cfg.V1X1SrcPath, "" ) )
		} )
		return nil
	} )
}
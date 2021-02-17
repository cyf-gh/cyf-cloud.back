package http

import (
	"../../cc"
	"../../cc/err"
	"../orm"
)

func init() {

cc.AddActionGroup( "/v1x1/account/fav", func( a cc.ActionGroup ) error {
	// \brief 获取收藏夹的文章
	a.GET( "/post/info", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
		id, e := GetIdByAtk( ap.R ); err.Assert( e )
		pis, e := orm.GetAllFavPostInfos( id ); err.Assert( e )
		epis, e :=  extendPostInfo( pis ); err.Assert( e )

		return cc.HerOkWithData( epis )
	} )
	return nil
})

}

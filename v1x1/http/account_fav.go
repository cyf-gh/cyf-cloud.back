package http

import (
	"../../cc"
	"../../cc/err"
	"../orm"
)

func init() {

cc.AddActionGroup( "/v1x1/account/fav", func( a cc.ActionGroup ) error {
	// 获取收藏夹的文章
	a.GET( "/post/info", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, int ) {
		id, e := GetIdByAtk( ap.R ); err.Check( e )
		pis, e := orm.GetAllFavPostInfos( id ); err.Check( e )
		epis, e :=  extendPostInfo( pis ); err.Check( e )

		return cc.HerOkWithData( epis )
	} )
	return nil
})

}

package http

import (
	"../../cc"
	"../../cc/err"
	"../orm"
	"encoding/json"
	"io/ioutil"
	"stgogo/comn/convert"
)

func init() {
	cc.AddActionGroup( "/v1x1/post/fav", func( a cc.ActionGroup ) error {
		a.GET( "/add", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var (
				e error
				pid string
			)
			if pid = ap.R.FormValue("id"); pid == "" {
				return cc.HerArgInvalid( "id" )
			}
			id, e := GetIdByAtk( ap.R ); err.Assert( e )
			npid, e := convert.Atoi64( pid ); err.Assert( e )
			_, e = orm.AddFav( id, npid ); err.Assert( e )

			return cc.HerOk()
		} )

		a.POST( "/update", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var (
				e error
				favList []int64
			)
			b, e := ioutil.ReadAll( ap.R.Body); err.Assert( e )
			e = json.Unmarshal( b, &favList ); err.Assert( e )

			id, e := GetIdByAtk( ap.R ); err.Assert( e )
			_, e = orm.UpdateFav( id, favList ); err.Assert( e )
			return cc.HerOk()
		} )

		a.GET( "/remove", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var (
				e error
				pid string
			)
			if pid = ap.R.FormValue("id"); pid == "" {
				return cc.HerArgInvalid( "id" )
			}
			id, e := GetIdByAtk( ap.R ); err.Assert( e )
			npid, e := convert.Atoi64( pid ); err.Assert( e )

			_, e = orm.RemoveFav( id, npid ); err.Assert( e )
			return cc.HerOk()
		} )

		a.GET( "/check", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var (
				e error
				pid string
				isFav bool
			)
			if pid = ap.R.FormValue("id"); pid == "" {
				return cc.HerArgInvalid( "id" )
			}
			id, e := GetIdByAtk( ap.R ); err.Assert( e )
			npid, e := convert.Atoi64( pid ); err.Assert( e )

			isFav, e = orm.IsPostFav( id, npid ); err.Assert( e )
			return cc.HerOkWithString( convert.Bool2a( isFav ) )
		} )

		return nil
	} )
}
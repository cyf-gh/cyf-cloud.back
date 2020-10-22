package http

import (
	"../err"
	"../orm"
	"encoding/json"
	"net/http"
)

func GetAccountFav( w http.ResponseWriter, r *http.Request ) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	var (
		e error
		postsB []byte
	)

	id, e := GetIdByAtk( r ); err.Check( e )
    pis, e := orm.GetAllFavPostInfos( id ); err.Check( e )
	epis, e :=  extendPostInfo( pis ); err.Check( e )

	{
		postsB, e = json.Marshal( epis ); err.Check( e )
	}
	err.HttpReturnOkWithData( &w, string(postsB) )
}

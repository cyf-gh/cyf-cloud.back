package http

import (
	"../../cc/err"
	"../orm"
	"errors"
	"net/http"
)

func DM1CheckPermission( r *http.Request ) ( error ) {
	id, e := GetIdByAtk( r )
	has, e := orm.DMCheckPermission( id ); err.Check( e )
	if !has {
		return errors.New("no permission to dm" )
	}
	return nil
}
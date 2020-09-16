package http

import (
	"net/http"
	orm "../orm"
)

func GetCid( r *http.Request ) ( string, error ) {
	cl, e  := r.Cookie("cid")
	if e == nil {
		return "", e
	}
	cid := cl.Value
	return cid, e
}

func GetAccountByCid( r *http.Request ) ( *orm.Account, error ) {
	cid, e :=  GetCid(r)
	if e == nil {
		return nil, e
	}
	return orm.GetAccount( AccessTokens[cid] )
}
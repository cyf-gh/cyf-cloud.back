package http

import (
	"github.com/kpango/glg"
	"net/http"
	orm "../orm"
)

func GetCid( r *http.Request ) ( string, error ) {
	cl, e  := r.Cookie("cid")
	if e != nil {
		glg.Error("cid not found. it may be a post proxy problem")
		return "", e
	}
	cid := cl.Value
	glg.Success("cid is (" + cid + ")")
	return cid, e
}

func GetAccountByCid( r *http.Request ) ( *orm.Account, error ) {
	cid, e :=  GetCid(r)
	if e != nil {
		return nil, e
	}
	return orm.GetAccount( AccessTokens[cid] )
}
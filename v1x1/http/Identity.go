package http

import (
	orm "../orm"
	"github.com/kpango/glg"
	"net/http"
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

func GetAtk( r *http.Request ) ( string, error ) {
	cl, e  := r.Cookie("atk")
	if e != nil {
		glg.Error("cid not found. it may be a post proxy problem")
		return "", e
	}
	atk := cl.Value
	glg.Success("atk is (" + atk + ")")
	return atk, e
}

func GetAccountByAtk( r *http.Request ) ( *orm.Account, error ) {
	cid, e :=  GetAtk(r)
	if e != nil {
		return nil, e
	}
	return orm.GetAccount( AccessTokens[cid] )
}
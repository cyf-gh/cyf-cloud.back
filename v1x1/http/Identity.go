package http

import (
	orm "../orm"
	"github.com/kpango/glg"
	"net/http"
	sec "../security"
	"../cache"
	"strconv"
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
		glg.Error("atk not found. it may be a post proxy problem")
		return "", e
	}
	atk := cl.Value
	glg.Success("atk is (" + atk + ")")
	return atk, e
}

func GetIdByAtk( r *http.Request ) ( int64, error ) {
	atk, e :=  GetAtk(r)
	if e != nil {
		return -1, e
	}
	ids, e := cache.Get( atk )
	if e != nil {
		return -1, e
	}
	return strconv.ParseInt(ids, 10, 64)
}

func GetAccountByAtk( r *http.Request ) ( *orm.Account, error ) {
	id, e := GetIdByAtk( r )
	if e != nil {
		return nil, e
	}
	return orm.GetAccount( id )
}

func GetAccountExByAtk( r *http.Request ) ( *orm.AccountEx, error ) {
	id, e := GetIdByAtk( r )
	if e != nil {
		return nil, e
	}
	return orm.GetAccountEx( id )
}

func CreateAtk( id int64, exp int ) (string, error) {
	var token string

	if exp == orm.TIME_EXPIRE_ONE_DAY {
		token = sec.GenerateAtkSession()
	} else {
		token = sec.GenerateAtk()
	}
	// 添加token到token列表中
	_, e := cache.SetExp(token, id, exp )
	return token, e
}
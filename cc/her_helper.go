package cc

import (
	"./err"
	"./err_code"
	"encoding/json"
	"net/http"
)

// data将会自动转化为json
func HerOkWithData( data interface{} ) ( HttpErrReturn, int ) {
	jn, e := json.Marshal( data ); err.Check( e )
	return HttpErrReturn{
		ErrCod: err_code.ERR_OK,
		Desc:   "ok",
		Data:   string(jn),
	}, http.StatusOK
}

// data将会自动转化为json
func HerOk() ( HttpErrReturn, int ) {
	return HttpErrReturn{
		ErrCod: err_code.ERR_OK,
		Desc:   "ok",
		Data:   "",
	}, http.StatusOK
}
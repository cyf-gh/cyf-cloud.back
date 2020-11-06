package cc

import (
	"./err"
	"./err_code"
	"encoding/json"
	"net/http"
)

// data将会自动转化为json
func HerOkWithData( data interface{} ) ( HttpErrReturn, StatusCode ) {
	jn, e := json.Marshal( data ); err.Check( e )
	return HttpErrReturn{
		ErrCod: err_code.ERR_OK,
		Desc:   "ok",
		Data:   string(jn),
	}, http.StatusOK
}

// data将会自动转化为json
func HerOkWithString( str string ) ( HttpErrReturn, StatusCode ) {
	return HttpErrReturn{
		ErrCod: err_code.ERR_OK,
		Desc:   "ok",
		Data:   str,
	}, http.StatusOK
}

// data将会自动转化为json
func HerOk() ( HttpErrReturn, StatusCode ) {
	return HttpErrReturn{
		ErrCod: err_code.ERR_OK,
		Desc:   "ok",
		Data:   "",
	}, http.StatusOK
}

func HerArgInvalid( argName string ) ( HttpErrReturn, StatusCode ) {
	return HttpErrReturn{
		ErrCod: err_code.ERR_INVALID_ARGUMENT,
		Desc:   "invalid argument: \""+argName +"\"",
		Data:   "",
	}, http.StatusOK
}
package cc

import (
	"./err"
	"./err_code"
	"encoding/json"
	"net/http"
	"time"
)

// data将会自动转化为json
func HerOkWithData( data interface{} ) ( HttpErrReturn, StatusCode ) {
	jn, e := json.Marshal( data ); err.Assert( e )
	return HttpErrReturn{
		ErrCod: err_code.ERR_OK,
		Desc:   "ok",
		Data:   string(jn),
	}, http.StatusOK
}

// 携带花费时间
// Data块会分为 raw 与 usedTime 标记；raw 装载原始数据，usedTime 装载运行所用时间，单位为秒
func HerOkWithDataAndUsedTime( data interface{}, time time.Duration ) ( HttpErrReturn, StatusCode ) {
	d := H {
		"raw": data,
		"usedTime" : time.Seconds(),
	}
	jn, e := json.Marshal( d ); err.Assert( e )
	return HttpErrReturn{
		ErrCod: err_code.ERR_OK,
		Desc:   "ok",
		Data:   string(jn),
	}, http.StatusOK
}

// string
func HerOkWithString( str string ) ( HttpErrReturn, StatusCode ) {
	return HttpErrReturn{
		ErrCod: err_code.ERR_OK,
		Desc:   "ok",
		Data:   str,
	}, http.StatusOK
}

// http仅返回data，兼容resp
func HerData( str string ) ( HttpErrReturn, StatusCode ) {
	return HttpErrReturn{
		ErrCod: err_code.ERR_OK,
		Desc:   "nil",
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

// 用于在弃用的API中直接返回
// 请使用( a ActionGroup ) Deprecated
func HerDeprecated() ( HttpErrReturn, StatusCode ) {
	return HttpErrReturn{
		ErrCod: err_code.ERR_DEPRECATED,
		Desc:   "deprecated",
		Data:   "",
	}, http.StatusOK
}

func HerRaw () {

}
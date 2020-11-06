package cc

import (
	cfg "../config"
	"./err"
	"./err_code"
	"encoding/json"
	"fmt"
	"github.com/kpango/glg"
	"net/http"
	"runtime/debug"
)


// 用于返回http状态信息，格式为json
type (
	HttpErrReturn struct {
		ErrCod string // 内部错误代码，与http状态码不同见第四行
		Desc   string // 错误描述
		Data   string // 携带数据
	}
	MakeHERxxx func( desc, errcode string ) ( *HttpErrReturn, int )
	StatusCode int
)

// 创建一个错误返回
func MakeHER( desc, errcode string) *HttpErrReturn {
	her := new(HttpErrReturn)

	// 如果为部署模式，则隐藏错误信息
	if cfg.IsRunModeDep() && errcode == err_code.ERR_SYS {
		her.Desc = "server internal error"
		her.ErrCod = errcode
		return her
	}

	her.Desc = desc
	her.ErrCod = errcode
	return her
}

func HttpReturnHER( w* http.ResponseWriter, her *HttpErrReturn, statusCode StatusCode ) {
	defer func() {
		if e := recover(); e != nil {
			glg.Error( e )
			(*w).WriteHeader(http.StatusInternalServerError)
			// 这时data返回体为空，客户端应当作出null检查动作
		}
	}()

	(*w).WriteHeader( int(statusCode) )

	// 将her结构体转化为json
	bs, e := json.Marshal(*her); err.Check(e)
	_, e = (*w).Write(bs); err.Check(e)

	glg.Log( fmt.Sprintf( "[HttpReturn] - StatusCode:(%d) - HER (%s)", statusCode, her ))
}

// server Ok 请求返回成功
func MakeHER200( desc, errcode string ) ( *HttpErrReturn, int ) {
	return MakeHER( desc, errcode ), 200
}

// server Bad Request 请求非法，请检查错误信息后重试
func MakeHER400( desc, errcode string ) ( *HttpErrReturn, int ) {
	return MakeHER( desc, errcode ), 400
}

// server Unauthorized 未授权
func MakeHER401( desc, errcode string )  ( *HttpErrReturn, int ) {
	return MakeHER( desc, errcode ), 401
}

// server Not Found 没有这个资源
func MakeHER404( desc, errcode string )  ( *HttpErrReturn, int ) {
	return MakeHER( desc, errcode ), 404
}

// server Server Error 服务器内部错误
func MakeHER500( desc, errcode string )  ( *HttpErrReturn, int ) {
	return MakeHER( desc, errcode ) , 500
}

func HttpReturn( w* http.ResponseWriter, desc, errCode, data string, MakeHERxxx MakeHERxxx ) {
	defer func() {
		if e := recover(); e != nil {
			glg.Error( e )
			(*w).WriteHeader(http.StatusInternalServerError)
			// 这时data返回体为空，客户端应当作出null检查动作
		}
	}()

	her, statusCode := MakeHERxxx( desc, errCode )
	her.Data = data
	(*w).WriteHeader( statusCode )

	// 将her结构体转化为json
	bs, e := json.Marshal(her); err.Check(e)
	_, e = (*w).Write(bs); err.Check(e)

	glg.Log( fmt.Sprintf( "[HttpReturn] - StatusCode:(%d) - HER (%s)", statusCode, her ))
}

// 封装

func HttpRecoverBasic( w *http.ResponseWriter, re interface{} ) {
	debug.PrintStack()
	_ = glg.Error( re )
	HttpReturn( w, fmt.Sprint( re ), err_code.ERR_SYS, "", MakeHER200 )
}

func HttpReturnArgInvalid( w *http.ResponseWriter, argName string ) {
	HttpReturn( w, "invalid argument: \""+argName +"\"", err_code.ERR_INVALID_ARGUMENT, "", MakeHER200 )
}

func HttpReturnOk( w *http.ResponseWriter ) {
	HttpReturn( w, "ok", err_code.ERR_OK, "", MakeHER200 )
}

func HttpReturnOkWithData( w *http.ResponseWriter, data string ) {
	HttpReturn( w, "ok", err_code.ERR_OK, data, MakeHER200 )
}

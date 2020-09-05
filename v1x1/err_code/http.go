package err_code

// 用于返回http状态信息，格式为json
type HttpErrReturn struct {
	ErrCod string // 内部错误代码，与http状态码不同
	Desc   string // 错误描述
	Data   string // 携带数据
}

// 创建一个错误返回
func MakeHER( desc, errcode string) *HttpErrReturn {
	her := new(HttpErrReturn)

	her.Desc = desc
	her.ErrCod = errcode
	return her
}

type MakeHERxxx func( desc, errcode string ) ( *HttpErrReturn, int )

// http Ok 请求返回成功
func MakeHER200( desc, errcode string ) ( *HttpErrReturn, int ) {
	return MakeHER( desc, errcode ), 200
}

// http Bad Request 请求非法，请检查错误信息后重试
func MakeHER400( desc, errcode string ) ( *HttpErrReturn, int ) {
	return MakeHER( desc, errcode), 400
}

// http Unauthorized 未授权
func MakeHER401( desc, errcode string )  ( *HttpErrReturn, int ) {
	return MakeHER( desc, errcode), 401
}

// http Server Error 服务器内部错误
func MakeHER500( desc, errcode string )  ( *HttpErrReturn, int ) {
	return MakeHER( desc, errcode) , 401
}
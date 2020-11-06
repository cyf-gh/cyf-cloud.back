package err_code

// errcode一览
const (
	ERR_SECURITY = "-2" // 安全问题
	ERR_SYS = "-1"		// 系统错误；同时为服务器内部错误的标识
	ERR_INCORRECT = "-3"// 输入错误
	ERR_OK = "0"		// ok
	ERR_INVALID_ARGUMENT = "-4" // 参数错误
	ERR_NO_AUTH = "-5"
)

const (
	ERR_INIT_ORM = "abort at creating orm" // 数据库初始化错误
)

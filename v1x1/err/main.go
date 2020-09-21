package err

import (
	errCode "../err_code"
	"encoding/json"
	"fmt"
	"github.com/kpango/glg"
	"net/http"
	"os"
	"runtime/debug"
)

func Exit( desc string ) {
	glg.Error("server abort with description: "+  desc )
	os.Exit( 1 )
}

func HttpReturn( w* http.ResponseWriter, desc, errCode, data string, MakeHERxxx errCode.MakeHERxxx ) {
	defer func() {
		if err := recover(); err != nil {
			glg.Error( err )
			(*w).WriteHeader(http.StatusInternalServerError)
			// 这时data返回体为空，客户端应当作出null检查动作
		}
	}()

	her, statusCode := MakeHERxxx( desc, errCode )
	her.Data = data
	(*w).WriteHeader( statusCode )
	bs, err := json.Marshal(her) // 将her结构体转化为json
	CheckErr(err)
	_, err = (*w).Write(bs)
	CheckErr(err)

	glg.Log( fmt.Sprintf( "[HttpReturn] - StatusCode:(%d) - HER (%s)", statusCode, her ))
}

// 检查是否存在错误，如果有则抛出异常
// catch 例子：
//    defer func() {
//        if err := recover(); err != nil {
//	  	  }
//    }()
func CheckErr( err error ) {
	if err != nil {
		panic( err )
	}
}

func HttpRecoverBasic( w *http.ResponseWriter, re interface{} ) {
	debug.PrintStack()
	_ = glg.Error( re )
	HttpReturn( w, fmt.Sprint( re ), errCode.ERR_SYS, "", errCode.MakeHER200 )
}
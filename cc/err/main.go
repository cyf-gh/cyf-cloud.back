package err

import (
	"github.com/kpango/glg"
	"os"
)

func Exit( desc string ) {
	glg.Error("server abort with description: "+  desc )
	os.Exit( 1 )
}

// 检查是否存在错误，如果有则抛出异常
// catch 例子：
//    defer func() {
//        if err := recover(); err != nil {
//	  	  }
//    }()
func Check( err error ) {
	if err != nil {
		panic( err )
	}
}


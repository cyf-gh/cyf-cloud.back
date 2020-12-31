// 进行log封装，用于对cc log进行区分
package dm_1

import "github.com/kpango/glg"

func Info( val ...interface{}) {
	glg.Log( "[dm_1]", val )
}

func Fail( val ...interface{}) {
	glg.Fail( "[dm_1]", val )
}

func Success( val ...interface{}) {
	glg.Success( "[dm_1]", val )
}

func Fatal( val ...interface{}) {
	glg.Fatal( "[dm_1]", val )
}
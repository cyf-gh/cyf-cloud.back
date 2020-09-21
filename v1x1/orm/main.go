// orm/main.go
// 负责数据库读取
package orm

import (
	err "../err"
	err_code "../err_code"
	"github.com/kpango/glg"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var engine *xorm.Engine
var engine_post *xorm.Engine
var engine_chat *xorm.Engine

// 有可能panic创建数据库引擎的错误
func connectDb( ppEnginePost **xorm.Engine, dbName, dbPath string ) {
	var e error
	*ppEnginePost, e = xorm.NewEngine("sqlite3", dbPath + dbName )
	err.CheckErr( e )
	glg.Success("orm to " + dbName + " (sqlite3)")
}

func InitEngine( dbPath string ) {
	glg.Log("initizing orm...")
	defer func() {
		if r := recover(); r != nil {
			_ = glg.Error(r)
			err.Exit( err_code.ERR_INIT_ORM )
		}
	}()
	connectDb( &engine, "account.db", dbPath )
	Sync2Account()

	connectDb( &engine_post, "post.db", dbPath )
	Sync2Post()

	connectDb( &engine_chat, "chat.db", dbPath )
	Sync2Chat()

	// e = NewAccount("cyf","cyf-ms@hotmail.com","18217203406","19990908cyfcyfcyfcyf")
	// err.CheckErr( e )
}
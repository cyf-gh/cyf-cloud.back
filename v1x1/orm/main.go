// orm/mainServer.go
// 负责数据库读取
package orm

import (
	err "../../cc/err"
	err_code "../../cc/err_code"
	"github.com/kpango/glg"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var (
	engine_account *xorm.Engine
	engine_post *xorm.Engine
	engine_chat *xorm.Engine
	engine_dm *xorm.Engine
)

// 初始化sqlite数据库orm
// 有可能panic创建数据库引擎的错误
func connectDb( ppEnginePost **xorm.Engine, dbName, dbPath string ) {
	var e error
	*ppEnginePost, e = xorm.NewEngine("sqlite3", dbPath + dbName )
	err.Check( e )
	glg.Success("orm to " + dbName + " (sqlite3)")
}

func InitEngine( dbPath string ) {
	glg.Log("orm loading...")
	defer func() {
		if r := recover(); r != nil {
			_ = glg.Error(r)
			err.Exit( err_code.ERR_INIT_ORM )
		}
	}()
	connectDb( &engine_account, "account.db", dbPath )
	Sync2Account()

	connectDb( &engine_post, "post.db", dbPath )
	Sync2Post()

	connectDb( &engine_chat, "chat.db", dbPath )
	Sync2Chat()

	connectDb( &engine_dm, "dm_1.db", dbPath )
	Sync2DM1()

	// e = NewAccount("cyf","cyf-ms@hotmail.com","18217203406","19990908cyfcyfcyfcyf")
	// err.Check( e )
	glg.Log("orm finished loading...")
}
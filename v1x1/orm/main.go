// orm/main.go
// 负责数据库读取
package orm

import (
	err "../err"
	"github.com/kpango/glg"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var engine *xorm.Engine
var engine_post *xorm.Engine

func InitEngine() {
	glg.Log("initizing orm...")
	var e error
	defer func() {
		if r := recover(); r != nil {
			_ = glg.Error(r)
		}
	}()
	engine, e = xorm.NewEngine("sqlite3", "./.db/account.db")
	err.CheckErr( e )
	glg.Success("orm to account.db [sqlite]")
	Sync2Account()

	engine_post, e = xorm.NewEngine("sqlite3", "./.db/post.db")
	err.CheckErr( e )
	glg.Success("orm to post.db [sqlite]")
	Sync2Post()

	// e = NewAccount("cyf","cyf-ms@hotmail.com","18217203406","19990908cyfcyfcyfcyf")
	// err.CheckErr( e )
}
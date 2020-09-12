package orm

import (
	err "../err"
	"github.com/kpango/glg"
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

var engine *xorm.Engine

func InitEngine() {
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

	// e = NewAccount("cyf","cyf-ms@hotmail.com","18217203406","19990908cyfcyfcyfcyf")
	// err.CheckErr( e )
}
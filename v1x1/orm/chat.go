package orm

import (
	err "../err"
	"time"
)

type Chat struct {
	Id int64
	AccountId int64 `xorm:"unique"`
	Text string
	Date time.Time
}

func Sync2Chat() {
	e := engine_post.Sync2(new(Chat))
	err.Check( e )
}

func NewChat( accountId int64, text string, date time.Time ) error {
	//TODO
	return nil
}
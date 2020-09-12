package orm

import (
	err "../err"
	"errors"
)

type Account struct {
	Id int64
	Name string `xorm:"unique"`
	Email  string `xorm:"unique"`
	Phone  string `xorm:"unique"`
	Passwd string
}

func Sync2Account() {
	e := engine.Sync2(new(Account))
	err.CheckErr( e )
}

func NewAccount( name, email, phone, passwd string ) error {
	_, e := engine.Insert( &Account{
		Name:   name,
		Email:  email,
		Phone:  phone,
		Passwd: passwd,
	})
	return e
}

func GetAccount( id int64 ) (*Account, error) {
	a := &Account{}
	has, e := engine.ID(id).Get(a)
	if e != nil {
		return nil, e
	} else if !has {
		return nil, errors.New("account not found")
	}
	return a, nil
}

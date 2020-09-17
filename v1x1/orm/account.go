package orm

import (
	err "../err"
	"errors"
)

// 账户的基本信息
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
	e = engine.Sync2(new(AccountEx))
	err.CheckErr( e )
	e = engine.Sync2(new(AccountPermission))
	err.CheckErr( e )
}

func NewAccount( name, email, phone, passwd string ) error {
	_, e := engine.Table("Account").Insert( &Account{
		Name:   name,
		Email:  email,
		Phone:  phone,
		Passwd: passwd,
	})
	return e
}

func GetAccount( id int64 ) (*Account, error) {
	a := &Account{}
	has, e := engine.Table("Account").ID(id).Get(a)
	if e != nil {
		return nil, e
	} else if !has {
		return nil, errors.New("account not found")
	}
	return a, nil
}

func GetAccountByLoginType( login ,cryPswd, loginType string) (*Account, error) {
	a := new(Account)
	exists, _ := engine.Table("Account").Where( loginType + " = ?", login).Get(a)

	if !exists {
		return nil, errors.New("no such account")
	}

	if a.Passwd != cryPswd {
		return nil, errors.New("wrong password")
	}
	return a, nil
}


// 账户的额外信息
type AccountEx struct {
	Id int64
	AccountId int64 `xorm:"unique"`
	Avatar string // base64 数据
	Info string // 个人简介，markdown数据
}

type AccountPermission struct {
	Id int64
	AccountId int64 `xorm:"unique"`
	Level string // \see ACCOUNT_LEVEL_xxx
}

const (
	ACCOUNT_LEVEL_ADMIN = "admin"	// 网站管理员
	ACCOUNT_LEVEL_VIP = "vip"		// 网站vip
	ACCOUNT_LEVEL_NORMAL = "n"		// normal - 一般会员
)
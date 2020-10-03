package orm

import (
	err "../err"
	"errors"
)

// 账户的基本信息
type (
	Account struct {
		Id int64
		Name string `xorm:"unique"`
		Email  string `xorm:"unique"`
		Phone  string `xorm:"unique"`
		Passwd string
	}
	// 账户的额外信息
	 AccountEx struct {
		Id int64
		AccountId int64 `xorm:"unique"`
		Avatar string // base64 数据
		Info string // 个人简介，markdown数据
		Level string // \see ACCOUNT_LEVEL_xxx

		Exp int64 // 等级？
		PrivateInfoMask string // 用于遮罩哪些信息不可被外部访问
	}
)

const (
	ACCOUNT_LEVEL_ADMIN = "admin"	// 网站管理员
	ACCOUNT_LEVEL_VIP = "vip"		// 网站vip
	ACCOUNT_LEVEL_NORMAL = "n"		// 一般会员
	ACCOUNT_LEVEL_UNREGISTERED = "unrgstr" // 还未验证账户

	ACCOUNT_LEVEL_TEST = "t" 		// 测试账户，鸟用没有
)

func Sync2Account() {
	e := engine.Sync2(new(Account))
	err.Check( e )
	e = engine.Sync2(new(AccountEx))
	err.Check( e )
}

func NewAccount( name, email, phone, passwd string ) error {
	id, e := engine.Table("Account").Insert( &Account{
		Name:   name,
		Email:  email,
		Phone:  phone,
		Passwd: passwd,
	})
	if e != nil { return e }

	_, e = engine.Table("account_ex").Insert( &AccountEx{
		AccountId: id,
		Avatar:    "",
		Info:      "",
		Level:     ACCOUNT_LEVEL_NORMAL,
		Exp:       0,
		PrivateInfoMask: "Phone",
	})
	return e
}

func SetAccountPhone( phone string, id int64 ) error {
	a := &Account{}
	a.Phone = phone
	_, e := engine.Table("account").ID( id ).Update( a )
	return e
}

func SetAccountExInfo( info string, id int64 ) error {
	ae := &AccountEx{}
	ae.Info = info
	_, e := engine.Table("account_ex").Where("account_id = ?", id ).Update( ae )
	return e
}

func SetAccountExAvatar( avatar string, id int64 ) error {
	ae := &AccountEx{}
	ae.Avatar = avatar
	_, e := engine.Table("account_ex").Where("account_id = ?", id ).Update( ae )
	return e
}

func GetAccountEx( id int64 ) ( *AccountEx, error ) {
	ae := &AccountEx{}
	has, e := engine.Table("account_ex").Where("account_id = ?", id ).Get(ae)
	if e != nil {
		return nil, e
	} else if !has {
		return nil, errors.New("account not found")
	}
	return ae, nil
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

func GetAccountByName( name string )  (*Account, error) {
	a := new(Account)
	exists, _ := engine.Table("Account").Where(  "name = ?", name).Get(a)
	if !exists {
		return nil, errors.New("no such account")
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


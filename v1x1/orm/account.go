package orm

import (
	err "../../cc/err"
	"errors"
	"fmt"
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

		BgUrl string // 装扮信息

		Exp int64 // 等级？
		PrivateInfoMask string // 用于遮罩哪些信息不可被外部访问

		FavPosts []int64 // 收藏文章

		SteamId string
		SteamApi string
	}
)

const (
	ACCOUNT_LEVEL_ADMIN = "admin"	// 网站管理员
	ACCOUNT_LEVEL_VIP = "vip"		// 网站vip
	ACCOUNT_LEVEL_NORMAL = "n"		// 一般会员
	ACCOUNT_LEVEL_UNREGISTERED = "unrgstr" // 还未验证账户

	ACCOUNT_LEVEL_TEST = "t" 		// 测试账户，鸟用没有

	TIME_EXPIRE_ONE_MONTH = 2626560
	TIME_EXPIRE_ONE_DAY = 87552

)

func Sync2Account() {
	e := engine_account.Sync2(new(Account))
	err.Assert( e )
	e = engine_account.Sync2(new(AccountEx))
	err.Assert( e )
}

func NewAccount( name, email, phone, passwd string ) error {
	_, e := engine_account.Table("Account").Insert( &Account{
		Name:   name,
		Email:  email,
		Phone:  phone,
		Passwd: passwd,
	})
	if e != nil { return e }

	a, e := GetAccountByName( name ); if e != nil { return e }

	_, e = engine_account.Table("account_ex").Insert( &AccountEx{
		AccountId: a.Id,
		Avatar:    "",
		Info:      "",
		Level:     ACCOUNT_LEVEL_NORMAL,
		Exp:       0,
		PrivateInfoMask: "Phone",
		BgUrl: "",
		FavPosts: []int64{},
	})
	return e
}

func SetAccountPhone( phone string, id int64 ) error {
	a := &Account{}
	a.Phone = phone
	_, e := engine_account.Table("account").ID( id ).Update( a )
	return e
}

func SetAccountExInfo( info string, id int64 ) error {
	ae := &AccountEx{}
	ae.Info = info
	_, e := engine_account.Table("account_ex").Where("account_id = ?", id ).Update( ae )
	return e
}

func SetAccountExAvatar( avatar string, id int64 ) error {
	ae := &AccountEx{}
	ae.Avatar = avatar
	_, e := engine_account.Table("account_ex").Where("account_id = ?", id ).Update( ae )
	return e
}

func GetAccountEx( id int64 ) ( *AccountEx, error ) {
	ae := &AccountEx{}
	has, e := engine_account.Table("account_ex").Where("account_id = ?", id ).Get(ae)
	if e != nil {
		return nil, e
	} else if !has {
		return nil, errors.New("account ex not found")
	}
	return ae, nil
}

func GetAccount( id int64 ) (*Account, error) {
	a := &Account{}
	has, e := engine_account.Table("Account").ID(id).Get(a)
	if e != nil {
		return nil, e
	} else if !has {
		return nil, errors.New("account not found")
	}
	return a, nil
}

func GetAccountByName( name string )  (*Account, error) {
	a := new(Account)
	exists, _ := engine_account.Table("Account").Where(  "name = ?", name).Get(a)
	if !exists {
		return nil, errors.New("no such account")
	}
	return a, nil
}

func GetAccountByLoginType( login ,cryPswd, loginType string) (*Account, error) {
	a := new(Account)
	exists, _ := engine_account.Table("Account").Where( loginType + " = ?", login).Get(a)

	if !exists {
		return nil, errors.New("no such account")
	}

	if a.Passwd != cryPswd {
		return nil, errors.New("wrong password")
	}
	return a, nil
}

// --------- 搜索模块 ---------

type (
	UserSearchResult struct {
		Id int64
		Avatar string // base64 数据
		Name string
	}
)

func VagueSearchAccountName( text string ) ( usrs []UserSearchResult, e error ) {
	var (
		as []Account
		ex *AccountEx
	)
	where := fmt.Sprintf("name like '%%%s%%'", text )
	if e = engine_account.Table("Account").Where(where).Find(&as); e != nil {return}

	for _, a := range as {
		if ex, e = GetAccountEx( a.Id ); e != nil {
			return
		}
		usrs = append(usrs, UserSearchResult{
			Id:     a.Id,
			Avatar: ex.Avatar,
			Name:   a.Name,
		})
	}
	return
}


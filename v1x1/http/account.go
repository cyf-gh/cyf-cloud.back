// 处理账户的登录，注册等情况
package http

import (
	"../cache"
	err "../err"
	err_code "../err_code"
	orm "../orm"
	sec "../security"
	"encoding/json"
	"github.com/kpango/glg"
	"io/ioutil"
	"net/http"
	"strings"
)

type (
	RegisterModel struct {
		Name string `json:"name"`
		Email string `json:"email"`
		Phone string  `json:"phone"`
		Pswd string   `json:"pswd"`
		Cap string  `json:"cap"`
	}
	LoginModel struct {
		Login string `json:"login"`
		Pswd string `json:"pswd"`
		LoginType  string  `json:"loginType"`
		// LoginType 应在前端进行完解析为
		// email name phone 三种之一
		KeepLogin bool `json:"keepLogin"`
	}
	InfoModel struct {
		Name string
		Email string
		Phone string
		Avatar string
		Info string
		Level string
		BgUrl string
	}
)
// 注册
func Register(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	var e error
	cid, e := GetCid( r )
	var registerModel RegisterModel

	b, e := ioutil.ReadAll(r.Body)
	err.Check(e)
	e = json.Unmarshal( b, &registerModel )
	glg.Log( registerModel )
	err.Check(e)

	if len( registerModel.Cap ) != 4 {
		err.HttpReturn(&w, "wrong captcha", err_code.ERR_INCORRECT, "", err_code.MakeHER200 )
		return
	}

	if false == sec.CaptchaVerify( &w, registerModel.Cap, cid ) {
		return
	}

	cryPswd := sec.CryptoPasswd( registerModel.Pswd )

	if registerModel.Phone == "" {
		registerModel.Phone = sec.GetRandom();
	}
	e = orm.NewAccount( registerModel.Name, registerModel.Email, registerModel.Phone, cryPswd )
	err.Check(e)
	err.HttpReturnOk( &w )
}


// 登录
func Login( w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	var (
		loginModel LoginModel
		maxAge int
		tokenCl http.Cookie
	)
	b, e := ioutil.ReadAll(r.Body); 		err.Check(e)
	e = json.Unmarshal( b, &loginModel ); 	err.Check(e)

	account, e := orm.GetAccountByLoginType( loginModel.Login,sec.CryptoPasswd( loginModel.Pswd ), loginModel.LoginType )
	err.Check( e )

	if loginModel.KeepLogin {
		maxAge = orm.TIME_EXPIRE_ONE_MONTH
	} else {
		maxAge = orm.TIME_EXPIRE_ONE_DAY
}
	token, e := CreateAtk( account.Id, maxAge ); err.Check( e )
	tokenCl = http.Cookie{ Name:"atk", Value:token, Path:"/", MaxAge: maxAge }

	http.SetCookie(w, &tokenCl)
	err.HttpReturnOk( &w )
}

// 直接返回所有的信息
func PrivateUserInfo( w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()

	info, e := getRawInfoByAtk( r ); err.Check( e )

	b, e := json.Marshal( info ); err.Check( e )
	err.HttpReturnOkWithData( &w, string( b ) )
}

// 返回公开的信息
func PublicUserInfo( w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	user := r.FormValue("user")
	info, mask, e := getRawInfoByName( r, user ); err.Check( e )

	b, e := json.Marshal( info ); err.Check( e )
	pinfo, e := createInfoMask( b, mask ); err.Check( e )
	err.HttpReturnOkWithData( &w, pinfo )
}

func UploadAvatar( w http.ResponseWriter, r *http.Request ) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()

	b, e := ioutil.ReadAll(r.Body); err.Check( e )
	id, e := GetIdByAtk( r )
	e = orm.SetAccountExAvatar( string( b ), id ); err.Check( e )
	err.HttpReturnOk( &w )
}

func UploadPhone( w http.ResponseWriter, r *http.Request ) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()

	phone := r.FormValue("phone")
	id, e := GetIdByAtk( r )
	e = orm.SetAccountPhone( phone, id ); err.Check( e )
	err.HttpReturnOk( &w )
}

// 移除cookie操作由前端完成，后端仅消除atk
func Logout( w http.ResponseWriter, r *http.Request ) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()

	atk, e := GetAtk( r ); err.Check( e )

	tokenCl := http.Cookie{Name:"atk", Path:"/", MaxAge: -1 }
	http.SetCookie( w, &tokenCl )
	e = cache.Del( atk ); err.Check( e )
	err.HttpReturnOk( &w )
}

func UploadInfo( w http.ResponseWriter, r *http.Request ) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()

	info := r.FormValue("info")
	atk, e := GetIdByAtk( r )
	e = orm.SetAccountExInfo( info, atk ); err.Check( e )
	err.HttpReturnOk( &w )
}

func copyInfoFromAAE( a *orm.Account, ae *orm.AccountEx) *InfoModel {
	return &InfoModel{
		Name:   a.Name,
		Email:  a.Email,
		Phone:  a.Phone,
		Avatar: ae.Avatar,
		Info:   ae.Info,
		Level:  ae.Level,
		BgUrl:  ae.BgUrl,
	}
}

func getRawInfoByName( r *http.Request, userName string ) (*InfoModel, string, error) {
	a, e := orm.GetAccountByName( userName )
	if e != nil {
		return nil, "", e
	}
	ae, e := orm.GetAccountEx( a.Id )
	if e != nil {
		return nil, "", e
	}

	info := copyInfoFromAAE( a ,ae )
	return info, ae.PrivateInfoMask, e
}

// 将账户信息转化为InfoModel
func getRawInfoByAtk( r *http.Request ) (*InfoModel, error) {
	a, e := GetAccountByAtk( r )
	if e != nil {
		return nil, e
	}
	ae, e := GetAccountExByAtk( r )
	if e != nil {
		return nil, e
	}

	info := copyInfoFromAAE( a ,ae )
	return info, e
}

// 创建遮罩，使被遮罩的数据不可被访问
func createInfoMask( b []byte, mask string ) ( string, error ) {
	ms := strings.Split( mask, ",")
	info := make(map[string]string)

	if e := json.Unmarshal( b, &info ); e != nil {
		return "", e
	}
	for _, m := range ms {
		if _, ok := info[m]; ok {
			info[m]="___cyfcloud_secret___"
		}
	}
	bi, e := json.Marshal(info)
	if e != nil {
		return "", e
	}
	return string( bi ), e
}

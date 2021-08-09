// 处理账户的登录，注册等情况
package http

import (
	"../../cc"
	err "../../cc/err"
	"../../cc/err_code"
	"../cache"
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
		Id int64
		Name string
		Email string
		Phone string
		Avatar string
		Info string
		Level string
		BgUrl string
		FavPost []orm.PostInfo
	}
)

func init() {
	cc.AddActionGroup("/v1x1/account", func(a cc.ActionGroup) error {
		a.POST("/register", func(ap cc.ActionPackage) (cc.HttpErrReturn, cc.StatusCode) {
			var e error
			cid, e := GetCid(ap.R)
			var registerModel RegisterModel

			b, e := ioutil.ReadAll(ap.R.Body)
			err.Assert(e)
			e = json.Unmarshal(b, &registerModel)
			glg.Log(registerModel)
			err.Assert(e)

			if len(registerModel.Cap) != 4 {
				return cc.HttpErrReturn{
					ErrCod: err_code.ERR_INCORRECT,
					Desc:   "wrong captcha",
					Data:   "",
				}, http.StatusOK
			}

			if false == sec.CaptchaVerify(registerModel.Cap, cid) {
				return cc.HttpErrReturn{
					ErrCod: err_code.ERR_INCORRECT,
					Desc:   "wrong captcha",
					Data:   "",
				}, http.StatusOK
			}

			cryPswd := sec.CryptoPasswd(registerModel.Pswd)

			if registerModel.Phone == "" {
				registerModel.Phone = sec.GetRandom();
			}
			e = orm.NewAccount(registerModel.Name, registerModel.Email, registerModel.Phone, cryPswd); err.Assert(e)
			return cc.HerOk()
		})
		// \brief 登陆账户
		// \arg[loginModel] LoginModel
		// \note 应在前端进行完解析为
		// email name phone 三种之一
		a.POST("/login", func(ap cc.ActionPackage) (cc.HttpErrReturn, cc.StatusCode) {
			var (
				loginModel LoginModel
				maxAge     int
				tokenCl    http.Cookie
			)
			b, e := ioutil.ReadAll(ap.R.Body); err.Assert(e)
			e = json.Unmarshal(b, &loginModel); err.Assert(e)

			account, e := orm.GetAccountByLoginType(loginModel.Login, sec.CryptoPasswd(loginModel.Pswd), loginModel.LoginType);err.Assert(e)

			if loginModel.KeepLogin {
				maxAge = orm.TIME_EXPIRE_ONE_MONTH
			} else {
				maxAge = orm.TIME_EXPIRE_ONE_DAY
			}
			token, e := CreateAtk(account.Id, maxAge); err.Assert(e)
			tokenCl = http.Cookie{Name: "atk", Value: token, Path: "/", MaxAge: maxAge}

			ap.SetCookie( &tokenCl )
			return cc.HerOkWithString(token)
		})

		a.POST( "/logout", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			atk, e := ap.GetCookie( "atk" ); err.Assert( e )

			tokenCl := http.Cookie{ Name:"atk", Path:"/", MaxAge: -1 }
			ap.SetCookie( &tokenCl )
			e = cache.Del( atk ); err.Assert( e )
			return cc.HerOk()
		} )

		a.GET("/private/info", func(ap cc.ActionPackage) (cc.HttpErrReturn, cc.StatusCode) {
			info, e := getRawInfoByAtk(ap.R); err.Assert(e)
			return cc.HerOkWithData(info)
		})

		a.Deprecated("/v1x1/account/info/home?uid=").GET("/public/info", func(ap cc.ActionPackage) (cc.HttpErrReturn, cc.StatusCode) {
			user := ap.R.FormValue("user")
			info, mask, e := getRawInfoByName(ap.R, user); err.Assert(e)

			b, e := json.Marshal(info); err.Assert(e)
			pInfo, e := createInfoMask(b, mask); err.Assert(e)
			return cc.HerOkWithString(pInfo)
		})

		a.POST("/upload/avatar", func(ap cc.ActionPackage) (cc.HttpErrReturn, cc.StatusCode) {
			b, e := ioutil.ReadAll( ap.R.Body ); err.Assert(e)
			id, e := GetIdByAtk( ap.R ); err.Assert( e )
			e = orm.SetAccountExAvatar( string( b ), id ); err.Assert( e )
			return cc.HerOk()
		})

		a.POST( "/update/phone", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			phone := ap.R.FormValue("phone")
			id, e := GetIdByAtk( ap.R ); err.Assert( e )
			e = orm.SetAccountPhone( phone, id ); err.Assert( e )
			return cc.HerOk()
		} )

		a.GET( "/update/description", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			info := ap.R.FormValue("info")
			atk, e := GetIdByAtk( ap.R )
			e = orm.SetAccountExInfo( info, atk ); err.Assert( e )
			return cc.HerOk()
		} )
		
		return nil
	})
}

func copyInfoFromAAE( a *orm.Account, ae *orm.AccountEx) (*InfoModel, error) {
	ps, e := orm.GetPostInfosByIds( ae.FavPosts )
	return &InfoModel{
		Id: a.Id,
		Name:   a.Name,
		Email:  a.Email,
		Phone:  a.Phone,
		Avatar: ae.Avatar,
		Info:   ae.Info,
		Level:  ae.Level,
		BgUrl:  ae.BgUrl,
		FavPost: ps,
	}, e
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

	info, e := copyInfoFromAAE( a ,ae ); err.Assert( e )
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

	info, e := copyInfoFromAAE( a ,ae ); err.Assert( e )
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


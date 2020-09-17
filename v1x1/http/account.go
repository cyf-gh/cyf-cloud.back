// 处理账户的登录，注册等情况
package http

import (
	err "../err"
	err_code "../err_code"
	orm "../orm"
	sec "../security"
	"encoding/json"
	"fmt"
	"github.com/kpango/glg"
	"io/ioutil"
	"net/http"
)

type RegisterModel struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Phone string  `json:"phone"`
	Pswd string   `json:"pswd"`
	Cap string  `json:"cap"`
}

// 注册
func Register(w http.ResponseWriter, r *http.Request) {
	enableCookies( &w )
	defer func() {
		if r := recover(); r  != nil {
			_ = glg.Error(r)
			err.HttpReturn(&w, fmt.Sprint( r ), err_code.ERR_SYS, "", err_code.MakeHER200 )
		}
	}()
	var e error
	cid, e := GetCid( r )
	var registerModel RegisterModel

	b, e := ioutil.ReadAll(r.Body)
	err.CheckErr(e)
	e = json.Unmarshal( b, &registerModel )
	glg.Log( registerModel )
	err.CheckErr(e)

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
	err.CheckErr(e)
	err.HttpReturn(&w, "ok", err_code.ERR_OK, "", err_code.MakeHER200 )
}

type LoginModel struct {
	Login string `json:"login"`
	Pswd string `json:"pswd"`
	LoginType  string  `json:"loginType"`
	// LoginType 应在前端进行完解析为
	// email name phone 三种之一
}
var AccessTokens map[string]int64

// 登录
func Login( w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r  != nil {
			_ = glg.Error(r)
			err.HttpReturn(&w, fmt.Sprint( r ), err_code.ERR_SYS, "", err_code.MakeHER200 )
		}
	}()
	var loginModel LoginModel

	b, e := ioutil.ReadAll(r.Body)
	err.CheckErr(e)
	e = json.Unmarshal( b, &loginModel )
	err.CheckErr(e)

	account, e := orm.GetAccountByLoginType( loginModel.Login,sec.CryptoPasswd( loginModel.Pswd ), loginModel.LoginType )
	err.CheckErr( e )

	token := sec.GenerateAccessToken()
	// 添加token到token列表中
	AccessTokens[token] = account.Id

	tokenCl := http.Cookie{Name:"atk", Value:token, Path:"/", MaxAge:2592000}
	http.SetCookie(w, &tokenCl)

	err.HttpReturn(&w, "ok", err_code.ERR_OK, "", err_code.MakeHER200 )
}
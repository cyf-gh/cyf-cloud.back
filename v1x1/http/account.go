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
	Name, Email, Phone, Pswd, Cap string
}

// 注册
func Register(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r  != nil {
			_ = glg.Error(r)
			err.HttpReturn(&w, fmt.Sprint( r ), err_code.ERR_SYS, "", err_code.MakeHER200 )
		}
	}()
	var e error
	cl, e  := r.Cookie("cid")
	cid := cl.Value
	err.CheckErr(e)
	var registerModel RegisterModel

	b, e := ioutil.ReadAll(r.Body)
	err.CheckErr(e)
	e = json.Unmarshal( b, &registerModel )
	err.CheckErr(e)

	sec.CaptchaVerify( &w, registerModel.Cap, cid )

	cryPswd := sec.CryptoPasswd( registerModel.Pswd )

	e = orm.NewAccount( registerModel.Name, registerModel.Phone, registerModel.Cap, cryPswd )
	err.CheckErr(e)
	err.HttpReturn(&w, "ok", err_code.ERR_OK, "", err_code.MakeHER200 )
}

type LoginModel struct {
	Login, Pswd, LoginType  string
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

	account, e := orm.GetAccountByLoginType( loginModel.Login, loginModel.Pswd, loginModel.LoginType )
	err.CheckErr( e )

	token := sec.GenerateAccessToken()
	// 添加token到token列表中
	AccessTokens[token] = account.Id

	tokenCl := http.Cookie{Name:"atk", Value:token, Path:"/", MaxAge:2592000}
	http.SetCookie(w, &tokenCl)

	err.HttpReturn(&w, "ok", err_code.ERR_OK, "", err_code.MakeHER200 )
}
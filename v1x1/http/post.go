package http

import (
	"encoding/json"
	"fmt"
	"github.com/kpango/glg"
	"io/ioutil"
	"net/http"
	err "../err"
	err_code "../err_code"
	orm "../orm"
)

// 发布新文章
type PostModel struct {
	Title string
	Text string
	TagIds[] string
}

func NewPost( w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r  != nil {
			_ = glg.Error(r)
			err.HttpReturn(&w, fmt.Sprint( r ), err_code.ERR_SYS, "", err_code.MakeHER200 )
		}
	}()

	var post PostModel

	b, e := ioutil.ReadAll(r.Body)
	err.CheckErr( e )
	e = json.Unmarshal( b, &post )
	err.CheckErr( e )

	account, e := GetAccountByCid( r )
	err.CheckErr( e )
	e = orm.NewPost( post.Title, post.Text, account.Id, post.TagIds )
	err.CheckErr( e )
	err.HttpReturn(&w, "ok", err_code.ERR_OK, "", err_code.MakeHER200 )
}

// 修改文章
type ModifiedPostModel struct {
	Id int64
	Title string
	Text string
	TagIds[] string
}

func ModifyPost( w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r  != nil {
			_ = glg.Error(r)
			err.HttpReturn(&w, fmt.Sprint( r ), err_code.ERR_SYS, "", err_code.MakeHER200 )
		}
	}()

	var post ModifiedPostModel

	b, e := ioutil.ReadAll(r.Body)
	err.CheckErr( e )
	e = json.Unmarshal( b, &post )
	err.CheckErr( e )

	account, e := GetAccountByCid( r )
	err.CheckErr( e )
	e = orm.ModifyPost( post.Id, post.Title, post.Text, account.Id, post.TagIds )
	err.CheckErr( e )
	err.HttpReturn(&w, "ok", err_code.ERR_OK, "", err_code.MakeHER200 )
}

// 更改文章，没有文本内容
// 应对流量节约的情况
type ModifiedPostNoTextModel struct {
	Id int64
	Title string
	TagIds[] string
}

func ModifiedPostNoText( w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r  != nil {
			_ = glg.Error(r)
			err.HttpReturn(&w, fmt.Sprint( r ), err_code.ERR_SYS, "", err_code.MakeHER200 )
		}
	}()

	var post ModifiedPostNoTextModel

	b, e := ioutil.ReadAll(r.Body)
	err.CheckErr( e )
	e = json.Unmarshal( b, &post )
	err.CheckErr( e )

	account, e := GetAccountByCid( r )
	err.CheckErr( e )
	e = orm.ModifyPostNoText( post.Id, post.Title, account.Id, post.TagIds )
	err.CheckErr( e )
	err.HttpReturn(&w, "ok", err_code.ERR_OK, "", err_code.MakeHER200 )
}
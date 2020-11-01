
package http

import (
	err "../err"
	"../err_code"
	orm "../orm"
	"encoding/json"
	"github.com/kpango/glg"
	"io/ioutil"
	"net/http"
	"stgogo/comn/convert"
)

// 发布新文章
type (
	PostModel struct {
		Title string
		Text string
		TagIds[] string
		IsPrivate bool
	}
	PostReaderModel struct {
		Title string
		Text string
		Tags[] string
		Author string
		Date string
		MyPost bool
		IsPrivate bool
		ViewedCount int64
	}
)

func NewPost( w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()

	var (
		post PostModel
		id int64
	)

	b, e := ioutil.ReadAll(r.Body); err.Check( e )
	e = json.Unmarshal( b, &post ); err.Check( e )

	account, e := GetAccountByAtk( r ); err.Check( e ); glg.Log( account ); glg.Log( post )
	id, e = orm.NewPost( post.Title, post.Text, account.Id, post.TagIds, post.IsPrivate ); err.Check( e )
	err.HttpReturnOkWithData( &w, convert.I64toa(id) )
}

// 修改文章
type ModifiedPostModel struct {
	Id int64
	Title string
	Text string
	TagIds[] string
	IsPrivate bool
}

func ModifyPost( w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()

	var post ModifiedPostModel

	b, e := ioutil.ReadAll(r.Body); err.Check( e )
	e = json.Unmarshal( b, &post ); err.Check( e )

	account, e := GetAccountByAtk( r ); err.Check( e ); glg.Log( account ); glg.Log( post )
	e = orm.ModifyPost( post.Id, post.Title, post.Text, account.Id, post.IsPrivate, post.TagIds ); err.Check( e )
	err.HttpReturnOk( &w )
}

// 更改文章，没有文本内容
// 应对流量节约的情况
type ModifyPostNoTextModel struct {
	Id int64
	Title string
	TagIds[] string
}

func ModifyPostNoText( w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()

	var post ModifyPostNoTextModel

	b, e := ioutil.ReadAll(r.Body); err.Check( e )
	e = json.Unmarshal( b, &post ); err.Check( e )

	account, e := GetAccountByAtk( r ); err.Check( e )
	e = orm.ModifyPostNoText( post.Id, post.Title, account.Id, post.TagIds ); err.Check( e )
	err.HttpReturnOk( &w )
}

func GetPost( w http.ResponseWriter, r *http.Request ) {
	defer func() {
		if r := recover(); r != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	var (
		id int64
		e error
		postsB []byte
		p orm.Post
		myPost bool
	)
	myPost = false
	strId := r.FormValue("id")
	id, e = convert.Atoi64( strId ); err.Check( e )
	// 获取文章
	p, e = orm.GetPostById( id ); err.Check( e )

	myId, _ := GetIdByAtk( r ) // 没有权限也可以访问，可以为-1
	// 只有不是本人的私有文章才不返回
	if p.IsPrivate && myId != p.OwnerId {
		err.HttpReturn( &w, "target post is private, cannot access", err_code.ERR_NO_AUTH, "", err_code.MakeHER200)
		return
	}

	if myId == p.OwnerId {
		myPost = true
	}

	// 找出作者名字与tag名字
	a, e := orm.GetAccount( p.OwnerId ); err.Check( e )
	tags, e := orm.GetTagNames( p.TagIds ); err.Check( e )

	tP := &PostReaderModel{
		Title:  p.Title,
		Text:   p.Text,
		Tags:    tags,
		Author: a.Name,
		Date: p.CreateDate,
		MyPost: myPost,
		IsPrivate: p.IsPrivate,
		ViewedCount: getPostView( convert.I64toa( p.Id ) ),
	}

	{
		postsB, e = json.Marshal( tP ); err.Check( e )
	}
	err.HttpReturnOkWithData( &w, string(postsB) )
}

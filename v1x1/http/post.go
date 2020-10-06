package http

import (
	err "../err"
	errCod "../err_code"
	orm "../orm"
	"encoding/json"
	"errors"
	"github.com/kpango/glg"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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
			err.HttpRecoverBasic( &w, r )
		}
	}()

	var post PostModel

	b, e := ioutil.ReadAll(r.Body)
	err.Check( e )
	e = json.Unmarshal( b, &post )
	err.Check( e )

	account, e := GetAccountByAtk( r )
	err.Check( e )
	glg.Log( account )
	glg.Log( post )
	e = orm.NewPost( post.Title, post.Text, account.Id, post.TagIds )
	err.Check( e )
	err.HttpReturnOk( &w )
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
			err.HttpRecoverBasic( &w, r )
		}
	}()

	var post ModifiedPostModel

	b, e := ioutil.ReadAll(r.Body)
	err.Check( e )
	e = json.Unmarshal( b, &post )
	err.Check( e )

	account, e := GetAccountByAtk( r )
	err.Check( e )
	glg.Log( account )
	glg.Log( post )
	e = orm.ModifyPost( post.Id, post.Title, post.Text, account.Id, post.TagIds )
	err.Check( e )
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

}

func GetPosts( w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	user := r.FormValue("user")
	rg := r.FormValue("range")
	var (

		e error
		posts []orm.Post
		postsB []byte
		a * orm.Account
	)

	// 如果user参数为空，则获取所有人的文章
	if user != "" {
		a, e = orm.GetAccountByName( user )			; err.Check( e )
		posts, e = orm.GetPostsByOwnerPublic( a.Id ); err.Check( e )
	} else {
		posts, e = orm.GetPostsAll(); err.Check( e )
	}

	if rg != "" {
		head, end, e := getRange( rg ); err.Check( e )
		posts = posts[head:end]
	}

	{
		postsB, e = json.Marshal( posts ); 	err.Check( e )
	}
	err.HttpReturnOkWithData( &w, string(postsB) )
}

func GetMyPosts( w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	var (
		posts []orm.Post
		postsB []byte
		e error
	)
	rg := r.FormValue("range")

	a, e := GetAccountByAtk( r );	err.Check( e )
	posts, e = orm.GetPostsByOwnerAll( a.Id ); err.Check( e )

	if rg != "" {
		head, end, e := getRange( rg ); err.Check( e )
		posts = posts[head:end]
	}

	{
		postsB, e = json.Marshal( posts ); 	err.Check( e )
	}
	err.HttpReturnOkWithData( &w, string(postsB) )
}

func getRange( rg string ) ( int, int, error ){
	var (
		rga []string
		head, end int
		e error
	)
	// 如果range该参数为空，则不限定
	if rg != "" {
		// 限定获取文章的篇数
		if rga = strings.Split( rg, ":"); len(rga) != 2 {
			return -1, -1, errors.New("invalid range argument")
		}
		if head, e = strconv.Atoi( rga[0] ); e != nil {
			return -1, -1, e
		}
		if end, e = strconv.Atoi( rga[1] ); e != nil {
			return -1, -1, e
		}
		return head, end, nil
	}
}
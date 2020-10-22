
package http

import (
	err "../err"
	"../err_code"
	orm "../orm"
	"encoding/json"
	"errors"
	"github.com/kpango/glg"
	"io/ioutil"
	"net/http"
	"stgogo/comn/convert"
	"stgogo/comn/refactor"
	"strconv"
	"strings"
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
	}
	PostInfoModel struct {
		orm.Post
		Author string
		ViewedCount int64
		Tags []string
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
	}

	{
		postsB, e = json.Marshal( tP ); err.Check( e )
	}
	err.HttpReturnOkWithData( &w, string(postsB) )
}

func GetPostInfos( w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	user := r.FormValue("user")
	rg := r.FormValue("range")
	var (
		e error
		posts []orm.PostInfo
		postsB []byte
		a * orm.Account
	)

	// 如果user参数为空，则获取所有人的文章
	if user != "" {
		a, e = orm.GetAccountByName( user )			; err.Check( e )
		posts, e = orm.GetPostInfosByOwnerPublic( a.Id ); err.Check( e )
	} else {
		// 如果设定了范围，则取范围
		if rg != "" {
			start, count, e := getRange( rg ); err.Check( e )
			posts, e = orm.GetAllPublicPostInfosLimited( start, count ); err.Check( e )
		} else {
			posts, e = orm.GetPostInfosAll(); err.Check( e )
		}
	}

	epi, e := extendPostInfo(posts); err.Check( e )

	{
		postsB, e = json.Marshal( epi ); err.Check( e )
	}
	err.HttpReturnOkWithData( &w, string(postsB) )
}

func extendPostInfo( posts []orm.PostInfo ) ([]PostInfoModel, error) {
	var (
		pims []PostInfoModel
		a * orm.Account
		e error
	)
	aMap := make(map[int64]*orm.Account)

	for _, p := range posts {
		pim := PostInfoModel{}
		if e = refactor.CopyFields( &pim, p ); e != nil {
			return nil, e
		}
		// 获取阅读量
		pim.ViewedCount = getPostView( convert.I64toa( p.Id ) )
		pims = append(pims, pim)
	}

	// 获取笔者名字
	for i, p := range pims {
		if a = aMap[p.OwnerId]; a == nil {
			if a, e = orm.GetAccount( p.OwnerId ); e != nil {
				return nil, e
			}
			aMap[a.Id] = a
		}
		pims[i].Author = a.Name

		var tgs []string
		if tgs, e = orm.GetTagNames( pims[i].TagIds); e != nil {
			return nil, e
		}
		pims[i].Tags = tgs
	}


	return pims, nil
}

func GetMyPostInfos( w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	var (
		posts []orm.PostInfo
		postsB []byte
		e error
	)
	rg := r.FormValue("range")

	a, e := GetAccountByAtk( r );	err.Check( e )

	if rg != "" {
		start, count, e := getRange( rg ); err.Check( e )
		posts, e = orm.GetAllPublicPostInfosLimited( start, count ); err.Check( e )
	} else {
		posts, e = orm.GetPostInfosByOwnerAll( a.Id ); err.Check( e )
	}

	epi, e := extendPostInfo(posts); err.Check( e )

	{
		postsB, e = json.Marshal( epi ); 	err.Check( e )
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

func GetAllTags( w http.ResponseWriter, r *http.Request ) {
	var (
		e error
	)
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	tags, e := orm.GetAllTags(); err.Check( e )
	tb, e := json.Marshal( tags ); 	err.Check( e )

	err.HttpReturnOkWithData( &w, string(tb) )
}

func GetPostInfosByTags( w http.ResponseWriter, r *http.Request ) {
	var (
		e error
	)
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	tags := r.FormValue("tags")
	tgs := strings.Split( tags, ",")
	pis, e := orm.GetPostInfosByTags( tgs ); err.Check( e )

	ps, e := extendPostInfo( pis )

	pb, e := json.Marshal(ps); err.Check( e )

    err.HttpReturnOkWithData( &w, string( pb ) )
}
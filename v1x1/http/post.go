
package http

import (
	"../../cc"
	err "../../cc/err"
	"../../cc/err_code"
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
	// 修改文章
	ModifiedPostModel struct {
		Id int64
		Title string
		Text string
		TagIds[] string
		IsPrivate bool
	}
	// 更改文章，没有文本内容
	// 应对流量节约的情况
	ModifyPostNoTextModel struct {
		Id int64
		Title string
		TagIds[] string
	}
)

func init() {
	cc.AddActionGroup( "/v1x1/post", func( a cc.ActionGroup ) error {
		a.POST( "/create", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var (
				post PostModel
				id int64
			)

			b, e := ioutil.ReadAll( ap.R.Body ); err.Check( e )
			e = json.Unmarshal( b, &post ); err.Check( e )

			account, e := GetAccountByAtk( ap.R ); err.Check( e ); glg.Log( account ); glg.Log( post )
			id, e = orm.NewPost( post.Title, post.Text, account.Id, post.TagIds, post.IsPrivate ); err.Check( e )
			return cc.HerOkWithString( convert.I64toa(id) )
		} )

		a.POST( "/modify", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var post ModifiedPostModel

			b, e := ioutil.ReadAll( ap.R.Body); err.Check( e )
			e = json.Unmarshal( b, &post ); err.Check( e )

			account, e := GetAccountByAtk( ap.R ); err.Check( e ); glg.Log( account ); glg.Log( post )
			e = orm.ModifyPost( post.Id, post.Title, post.Text, account.Id, post.IsPrivate, post.TagIds ); err.Check( e )
			return cc.HerOk()
		} )

		a.POST( "/modifyNT", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var post ModifyPostNoTextModel

			b, e := ioutil.ReadAll(ap.R.Body); err.Check( e )
			e = json.Unmarshal( b, &post ); err.Check( e )

			account, e := GetAccountByAtk( ap.R ); err.Check( e )
			e = orm.ModifyPostNoText( post.Id, post.Title, account.Id, post.TagIds ); err.Check( e )
			return cc.HerOk()
		} )

		a.GET( "", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var (
				id int64
				e error
				p orm.Post
				myPost bool
			)
			myPost = false
			strId := ap.R.FormValue("id")
			id, e = convert.Atoi64( strId ); err.Check( e )
			// 获取文章
			p, e = orm.GetPostById( id ); err.Check( e )

			myId, _ := GetIdByAtk( ap.R ) // 没有权限也可以访问，可以为-1
			// 只有不是本人的私有文章才不返回
			if p.IsPrivate && myId != p.OwnerId {
				return cc.HttpErrReturn{
					ErrCod: err_code.ERR_NO_AUTH,
					Desc:   "target post is private, cannot access",
					Data:   "",
				}, http.StatusOK
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
				Date: p.Date,
				MyPost: myPost,
				IsPrivate: p.IsPrivate,
				ViewedCount: getPostView( convert.I64toa( p.Id ) ),
			}
			return cc.HerOkWithData( tP )
		} )
		return nil
	} )
}



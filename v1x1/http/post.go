
package http

import (
	"../../cc"
	err "../../cc/err"
	"../../cc/err_code"
	orm "../orm"
	"encoding/json"
	"errors"
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
		Path string
	}
	PostReaderModel struct {
		Id int64
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
		Path string
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
		// \brief 创建一篇文章
		// \return 新建Post的Id
		a.POST( "/create", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var (
				post PostModel
				id int64
			)

			b, e := ioutil.ReadAll( ap.R.Body ); err.Assert( e )
			e = json.Unmarshal( b, &post ); err.Assert( e )

			account, e := GetAccountByAtk( ap.R ); err.Assert( e )
			id, e = orm.NewPost( post.Title, post.Text, account.Id, post.TagIds, post.IsPrivate, post.Path ); err.Assert( e )
			return cc.HerOkWithString( convert.I64toa(id) )
		} )

		a.GET("/modify/path", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			path := ap.R.FormValue("path")
			pid := ap.R.FormValue("pid")
			uid, e := GetIdByAtk( ap.R ); err.Assert( e )
			npid, e := convert.Atoi64( pid ); err.Assert( e )

			post, e := orm.GetPostById( npid ); err.Assert( e )
			if post.OwnerId != uid {
				err.Assert( errors.New("cannot modify path of others' post") )
			}
			e = orm.ModifyPostPath( post.Id, path )
			return cc.HerOk()
		} )

		a.SetFreq(100).POST( "/modify", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var post ModifiedPostModel

			b, e := ioutil.ReadAll( ap.R.Body); err.Assert( e )
			e = json.Unmarshal( b, &post ); err.Assert( e )

			account, e := GetAccountByAtk( ap.R ); err.Assert( e )
			e = orm.ModifyPost( post.Id, post.Title, post.Text, account.Id, post.IsPrivate, post.TagIds, post.Path ); err.Assert( e )
			return cc.HerOk()
		} )

		// \brief 仅修改文章内容之外的数据
		// \note 仅修改了文章的标题、作者、标签等信息。此请求中不包含文章内容。
		// （设计用于节省流量）
		a.POST( "/modifyNT", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var post ModifyPostNoTextModel

			b, e := ioutil.ReadAll(ap.R.Body); err.Assert( e )
			e = json.Unmarshal( b, &post ); err.Assert( e )

			account, e := GetAccountByAtk( ap.R ); err.Assert( e )
			e = orm.ModifyPostNoText( post.Id, post.Title, account.Id, post.TagIds ); err.Assert( e )
			return cc.HerOk()
		} )
		// \brief 获取自定义样式
		// \note 如需要创建自定义样式，只需创建一篇名为MyMarkdownStyle为标题的文章即可，若有多篇该标题文章，则使用第一个该标题文章
		a.GET( "/custom/style", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode) {
			id, e := GetIdByAtk( ap.R ); err.Assert( e )
			css := orm.GetPostCustomStyleById( id )
			return cc.HerOkWithData( css )
		} )
		// \brief 获取文章
		// \params[style] 为force则强制在markdown中加载自定义样式
		a.GET( "", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var (
				id int64
				e error
				p orm.Post
				myPost bool
			)
			myPost = false
			strId := ap.R.FormValue("id")
			strStyle := ap.GetFormValue("style")
			id, e = convert.Atoi64( strId ); err.Assert( e )
			// 获取文章
			p, e = orm.GetPostById( id ); err.Assert( e )

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
			a, e := orm.GetAccount( p.OwnerId ); err.Assert( e )
			tags, e := orm.GetTagNames( p.TagIds ); err.Assert( e )

			// 添加Markdown CSS Style
			css := ""
			if strStyle == "force" && p.Title != "MyMarkdownStyle" {
				css = orm.GetPostCustomStyleById( myId )
			}

			tP := &PostReaderModel{
				Id: 	p.Id,
				Title:  p.Title,
				Text:   p.Text + css,
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



// 处理文章摘要的业务逻辑
package http

import (
	"../../cc"
	"../../cc/err"
	orm "../orm"
	"encoding/json"
	"errors"
	"stgogo/comn/convert"
	"stgogo/comn/refactor"
	"strconv"
	"strings"
)

type (
	PostInfoModel struct {
		orm.Post
		Author string
		ViewedCount int64
		Tags []string
	}
)

func init() {
	cc.AddActionGroup( "/v1x1/posts/info", func( a cc.ActionGroup ) error {
		a.GET( "", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			r := ap.R
			user := r.FormValue("user")
			rg := r.FormValue("range")
			var (
				e error
				posts []orm.PostInfo
				a * orm.Account
			)

			// 如果user参数为空，则获取所有人的文章
			if user != "" {
				a, e = orm.GetAccountByName( user )			; err.Assert( e )
				posts, e = orm.GetPostInfosByOwnerPublic( a.Id ); err.Assert( e )
			} else {
				// 如果设定了范围，则取范围
				if rg != "" {
					start, count, e := getRange( rg ); err.Assert( e )
					posts, e = orm.GetAllPublicPostInfosLimited( start, count ); err.Assert( e )
				} else {
					posts, e = orm.GetPostInfosAll(); err.Assert( e )
				}
			}

			epi, e := extendPostInfo(posts); err.Assert( e )

			return cc.HerOkWithData( epi )
		} )

		a.GET( "/by/tag", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var (
				e error
				tgs []string
			)
			r := ap.R

			tags := r.FormValue("tags")
			e = json.Unmarshal( []byte(tags), &tgs ); err.Assert( e )
			pis, e := orm.GetPostInfosByTags( tgs ); err.Assert( e )

			ps, e := extendPostInfo( pis )
			return cc.HerOkWithData( ps )
		} )

		a.GET( "/by/date", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var (
				e error
			)
			r := ap.R
			date := r.FormValue("date")
			id, e := GetIdByAtk( ap.R ); err.Assert( e )
			ps, e := orm.GetPostInfosByCateDate( id, date ); err.Assert( e )

			eps, e := extendPostInfo( ps ); err.Assert( e )
			return cc.HerOkWithData( eps )
		} )

		a.GET( "/self", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var (
				posts  []orm.PostInfo
				e      error
			)
			r := ap.R
			rg := r.FormValue("range")

			a, e := GetAccountByAtk(r); err.Assert(e)

			if rg != "" {
				start, count, e := getRange(rg); err.Assert(e)
				posts, e = orm.GetAllPublicPostInfosLimited(start, count); err.Assert(e)
			} else {
				posts, e = orm.GetPostInfosByOwnerAll(a.Id); err.Assert(e)
			}

			epi, e := extendPostInfo(posts); err.Assert(e)
			return cc.HerOkWithData( epi )
		} )
		return nil
	} )
}

// 将后台的PostInfo模型转化为前台的PostInfo模型
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



// 处理文章摘要的业务逻辑
package http

import (
	err "../err"
	orm "../orm"
	"encoding/json"
	"errors"
	"net/http"
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
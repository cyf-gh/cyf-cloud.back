// 编写频繁操作的api
// @see http/README.md
package http

import (
	"../cache"
	"../err"
	"../orm"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"stgogo/comn/convert"
	"strings"
)

const (
	_postViewPrefix = "$post_view$"
	_postLikeItPrefix = "$post_like_it$"
)

func getPostView( pid string ) int64 {
	k :=  _postViewPrefix + pid
	if ps, e := cache.Get( k ); e != nil {
		// 还没有键，创建键
		cache.Set( k, "0")
		return 0
	} else {
		p, _ := convert.Atoi64( ps )
		return p
	}
}

func doPostView( pid string ) {
	k :=  _postViewPrefix + pid

	_, _ =cache.Set( k, convert.I64toa( getPostView( pid ) + 1) ) // no error
}

func getPostLikeIt( pid string ) []string {
	k :=  _postLikeItPrefix + pid
	if ps, e := cache.Get( k ); e != nil {
		// 还没有键，创建键
		cache.Set( k, "")
		return nil
	} else {
		return strings.Split( ps, "," )
	}
}

func clickPostLikeIt( pid string , uid string ) {
	k :=  _postLikeItPrefix + pid
	removeLike := false

	ps := getPostLikeIt( pid )

	// 如果存在了喜欢，则取消喜欢
	for i, uidLiked := range ps {
		if uidLiked == uid {
			ps = append( ps[:i], ps[i+1:]... )
			removeLike = true; break
		}
	}
	// 如果没有进行取消喜欢操作，则添加喜欢
	if !removeLike {
		ps = append( ps, uid )
	}
	_, _ = cache.Set( k, ps )
}

// Actions

func GetViewCount( w http.ResponseWriter, r *http.Request ) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
    var (
		pid string
	)

	if pid = r.FormValue("id"); pid == "" {
		err.HttpReturnArgInvalid( &w, "id"); return
	}

    err.HttpReturnOkWithData( &w, convert.I64toa( getPostView( pid ) ) )
}

func ViewedPost( w http.ResponseWriter, r *http.Request ) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	var (
		pid string
	)
	if pid = r.FormValue("id"); pid == "" {
		err.HttpReturnArgInvalid( &w, "id"); return
	}
	doPostView( pid )
    err.HttpReturnOk( &w )
}

func AddFav( w http.ResponseWriter, r *http.Request ) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	var (
		e error
		favList []int64
	)
	b, e := ioutil.ReadAll(r.Body); err.Check( e )
	e = json.Unmarshal( b, &favList ); err.Check( e )

	id, e := GetIdByAtk( r ); err.Check( e )
	_, e = orm.UpdateFav( id, favList )

    err.HttpReturnOk( &w )
}
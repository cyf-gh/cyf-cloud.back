// 编写频繁操作的api
// @see http/README.md
package http

import (
	"../cache"
	"../err"
	"../orm"
	"../../cc"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"stgogo/comn/convert"
)

type (
	LikeInfoModel struct {
		Count int
		Liked bool
	}
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

func doPostView( pid string ) error {
	k :=  _postViewPrefix + pid

	_, e := cache.Set( k, convert.I64toa( getPostView( pid ) + 1) ) // no error
	return e
}

func getPostLikeIt( pid string ) ( postLikes []string, e error ) {
	k :=  _postLikeItPrefix + pid
	if ps, e := cache.Get( k ); e != nil {
		empty, _ := json.Marshal([]string{})
		// 还没有键，创建键
		cache.Set( k, empty)
		return nil, nil
	} else {
		e = json.Unmarshal( []byte(ps),&postLikes )
		return postLikes, e
	}
}

func clickPostLikeIt( pid string , uid string ) error {
	k :=  _postLikeItPrefix + pid
	removeLike := false

	if ps, e := getPostLikeIt( pid ); e != nil {
		return e
	} else {
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
		strPs, _  := json.Marshal( ps )
		_, e := cache.Set( k, strPs )
		return e
	}
}

func isLikeIt( pid, uid string ) (bool, error) {
	if ps, e := getPostLikeIt( pid ); e != nil {
		return false, e
	} else {
		// 如果存在了喜欢，则取消喜欢
		for _, uidLiked := range ps {
			if uidLiked == uid {
				return true, nil
			}
		}
		return false, nil
	}

}

func CCPostFreq() {

	cc.GET( "/v1x1/post/view/count", func( w http.ResponseWriter, r *http.Request ) {
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
	} )
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
	e := doPostView( pid ); err.Check( e )
    err.HttpReturnOk( &w )
}

func UpdateFav( w http.ResponseWriter, r *http.Request ) {
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
	_, e = orm.UpdateFav( id, favList ); err.Check( e )

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
		pid string
	)
	if pid = r.FormValue("id"); pid == "" {
		err.HttpReturnArgInvalid( &w, "id"); return
	}
	id, e := GetIdByAtk( r ); err.Check( e )
	npid, e := convert.Atoi64( pid ); err.Check( e )
	_, e = orm.AddFav( id, npid ); err.Check( e )

	err.HttpReturnOk( &w )
}

func RemoveFav( w http.ResponseWriter, r *http.Request ) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	var (
		e error
		pid string
	)
	if pid = r.FormValue("id"); pid == "" {
		err.HttpReturnArgInvalid( &w, "id"); return
	}
	id, e := GetIdByAtk( r ); err.Check( e )
	npid, e := convert.Atoi64( pid ); err.Check( e )

	_, e = orm.RemoveFav( id, npid ); err.Check( e )

    err.HttpReturnOk( &w )
}

func IsPostFav( w http.ResponseWriter, r *http.Request ) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	var (
		e error
		pid string
		isFav bool
	)
	if pid = r.FormValue("id"); pid == "" {
		err.HttpReturnArgInvalid( &w, "id"); return
	}
	id, e := GetIdByAtk( r ); err.Check( e )
	npid, e := convert.Atoi64( pid ); err.Check( e )

	isFav, e = orm.IsPostFav( id, npid ); err.Check( e )


    err.HttpReturnOkWithData( &w, convert.Bool2a( isFav ) )
}


func LikeIt( w http.ResponseWriter, r *http.Request ) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	var (
		e error
		pid string
	)
	if pid = r.FormValue("id"); pid == "" {
		err.HttpReturnArgInvalid( &w, "id"); return
	}
	id, e := GetIdByAtk( r ); err.Check( e )

	e = clickPostLikeIt( pid, convert.I64toa( id ) ); err.Check( e )
    err.HttpReturnOk( &w )
}

func IsLikeIt( w http.ResponseWriter, r *http.Request ) {
	defer func() {
		if r := recover(); r  != nil {
			err.HttpRecoverBasic( &w, r )
		}
	}()
	var (
		e error
		pid string
	)
	if pid = r.FormValue("id"); pid == "" {
		err.HttpReturnArgInvalid( &w, "id"); return
	}
	id, e := GetIdByAtk( r ); err.Check( e )

	isLike, e := isLikeIt( pid, convert.I64toa( id ) ); err.Check( e )

	err.HttpReturnOkWithData( &w, convert.Bool2a( isLike ) )
}

func LikeCount( w http.ResponseWriter, r *http.Request ) {
	defer func() {
		if r := recover(); r != nil {
			err.HttpRecoverBasic(&w, r)
		}
	}()
	var (
		e   error
		pid string

	)
	if pid = r.FormValue("id"); pid == "" {
		err.HttpReturnArgInvalid(&w, "id");
		return
	}
	id, e := GetIdByAtk( r ); err.Check( e )
	likes, e := getPostLikeIt( pid ); err.Check( e )
	nid := convert.I64toa( id )

	li := LikeInfoModel{
		Count: len( likes ),
		Liked: false,
	}
	for _, nnid := range likes {
		if nid == nnid {
			li.Liked = true
		}
	}

	bli, e := json.Marshal( li ); err.Check( e )

	err.HttpReturnOkWithData( &w, string(bli) )
}
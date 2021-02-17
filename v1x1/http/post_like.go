package http

import (
	"../../cc"
	"../../cc/err"
	"../cache"
	"encoding/json"
	"stgogo/comn/convert"
)

func init() {
	cc.AddActionGroup( "/v1x1/post/like", func( a cc.ActionGroup ) error {
		// \brief 点赞
		// \note 开关型API，如果已经点击了like则会取消like
	    a.GET( "", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var (
				e error
				pid string
			)
			if pid = ap.R.FormValue("id"); pid == "" {
				return cc.HerArgInvalid( "id" )
			}
			id, e := GetIdByAtk( ap.R ); err.Assert( e )

			e = clickPostLikeIt( pid, convert.I64toa( id ) ); err.Assert( e )

			return cc.HerOk()
	    } )

	    a.GET( "/count", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var (
				e   error
				pid string
			)
			if pid = ap.R.FormValue("id"); pid == "" {
				return cc.HerArgInvalid( "id" )
			}
			id, e := GetIdByAtk( ap.R ); err.Assert( e )
			likes, e := getPostLikeIt( pid ); err.Assert( e )
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
			return cc.HerOkWithData( li )
	    } )

	    a.GET( "/check", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			var (
				e error
				pid string
			)
			if pid = ap.R.FormValue("id"); pid == "" {
				return cc.HerArgInvalid( "id" )
			}
			id, e := GetIdByAtk( ap.R ); err.Assert( e )
			isLike, e := isLikeIt( pid, convert.I64toa( id ) ); err.Assert( e )
			return cc.HerOkWithString( convert.Bool2a( isLike ) )
	    } )

	    return nil
	} )
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

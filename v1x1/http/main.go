package http

import (
	mwh "../../middleware/helper"
	"../../run_mode"
	"net/http"
)

func Init() {
	run_mode.HttpInit()

	http.HandleFunc("/v1x1/account/register", mwh.WrapPost( Register ) )
	http.HandleFunc("/v1x1/account/login", mwh.WrapPost( Login ) )
	http.HandleFunc("/v1x1/account/logout", mwh.WrapPost( Logout ) )

	http.HandleFunc( "/v1x1/account/private/info", mwh.WrapGet( PrivateUserInfo ) )
	http.HandleFunc( "/v1x1/account/public/info", mwh.WrapGet( PublicUserInfo ) )
	http.HandleFunc( "/v1x1/account/upload/avatar", mwh.WrapPost( UploadAvatar ) )
	http.HandleFunc( "/v1x1/account/update/phone", mwh.WrapGet( UploadPhone ) )
	http.HandleFunc( "/v1x1/account/update/description", mwh.WrapGet( UploadInfo ) )

	http.HandleFunc("/v1x1/post/create", mwh.WrapPost( NewPost ) )
	http.HandleFunc("/v1x1/post/modify", mwh.WrapPost( ModifyPost ) )
	http.HandleFunc("/v1x1/post/modifyNT", mwh.WrapPost( ModifyPostNoText ) )
	http.HandleFunc("/v1x1/post", mwh.WrapGet(  GetPost ) )
	http.HandleFunc("/v1x1/posts/info", mwh.WrapGet( GetPostInfos ) )
	http.HandleFunc("/v1x1/posts/info/self", mwh.WrapGet( GetMyPostInfos ) )
	http.HandleFunc( "/v1x1/posts/info/by/tag", mwh.WrapGet( GetPostInfosByTags ) )
	http.HandleFunc( "/v1x1/tags", mwh.WrapGet( GetAllTags) )

	http.HandleFunc( "/v1x1/post/view", mwh.WrapGet( ViewedPost ) )
	// http.HandleFunc( "/v1x1/post/view/count", mwh.WrapGet( GetViewCount ) )

	http.HandleFunc( "/v1x1/post/fav/add", mwh.WrapGet( AddFav ) )
	http.HandleFunc( "/v1x1/post/fav/update", mwh.WrapPost( UpdateFav ) )
	http.HandleFunc( "/v1x1/post/fav/remove", mwh.WrapGet( RemoveFav ) )
	http.HandleFunc( "/v1x1/post/fav/check", mwh.WrapGet( IsPostFav ) )

	http.HandleFunc( "/v1x1/post/like", mwh.WrapGet( LikeIt ) )
	http.HandleFunc( "/v1x1/post/like/count", mwh.WrapGet( LikeCount ) )
	http.HandleFunc( "/v1x1/post/like/check", mwh.WrapGet( IsLikeIt ) )

	http.HandleFunc( "/v1x1/clipboard/push", mwh.WrapPost( ClipboardPush ) )
	http.HandleFunc( "/v1x1/clipboard/fetch", mwh.WrapGet( ClipboardFetch ) )
}

package http

import (
	mwh "../../middleware/helper"
	"../cache"
	"net/http"
)

func Init() {
	// 方便测试的令牌，cyfhaoshuai
	cache.Set( "cyfhaoshuaicyfhaoshuaicyfhaoshuaicyfhaoshuai", 1 )

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
	http.HandleFunc("/v1x1/post/modifyNT", mwh.WrapPost( ModifiyPostNoText ) )
	http.HandleFunc("/v1x1/posts", mwh.WrapPost(  GetPosts ) )

	http.HandleFunc( "/v1x1/clipboard/push", mwh.WrapPost( ClipboardPush ) )
	http.HandleFunc( "/v1x1/clipboard/fetch", mwh.WrapGet( ClipboardFetch ) )
}

func enableCookies(w *http.ResponseWriter) {
	// (*w).Header().Set("Access-Control-Allow-Origin", "https://localhost:8887")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

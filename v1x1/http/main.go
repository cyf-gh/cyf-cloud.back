package http

import (
	mwh "../../middleware/helper"
	"../../run_mode"
	"net/http"
)

func Init() {
	run_mode.HttpInit()

	http.HandleFunc("/v1x1/post/create", mwh.WrapPost( NewPost ) )
	http.HandleFunc("/v1x1/post/modify", mwh.WrapPost( ModifyPost ) )
	http.HandleFunc("/v1x1/post/modifyNT", mwh.WrapPost( ModifyPostNoText ) )
	http.HandleFunc("/v1x1/post", mwh.WrapGet(  GetPost ) )
	http.HandleFunc("/v1x1/posts/info", mwh.WrapGet( GetPostInfos ) )
	http.HandleFunc("/v1x1/posts/info/self", mwh.WrapGet( GetMyPostInfos ) )
	http.HandleFunc( "/v1x1/posts/info/by/tag", mwh.WrapGet( GetPostInfosByTags ) )
	http.HandleFunc( "/v1x1/tags", mwh.WrapGet( GetAllTags) )
}

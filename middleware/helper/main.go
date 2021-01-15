package helper

import (
	mw "../../middleware"
	mwu "../util"
	"net/http"
)

func WrapPost( handler http.HandlerFunc ) http.HandlerFunc {
	return mw.HandlerWrapFully( handler, mwu.Method( mwu.POST ) )
}

func WrapGet( handler http.HandlerFunc ) http.HandlerFunc {
	return mw.HandlerWrapFully( handler, mwu.Method( mwu.GET ) )
}

func WrapWS( handler http.HandlerFunc ) http.HandlerFunc {
	return mw.HandlerWrapFully( handler, mwu.Method( mwu.WS ) )
}


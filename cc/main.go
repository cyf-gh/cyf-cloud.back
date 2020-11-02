package cc

import (
	mwh "../middleware/helper"
	"github.com/kpango/glg"
	"net/http"
)

var (
	postHandlers map[string] *http.HandlerFunc
	getHandlers map[string] *http.HandlerFunc
)

func init() {
	postHandlers = make( map[string] *http.HandlerFunc )
	getHandlers = make( map[string] *http.HandlerFunc )
}

func POST( path string, handler http.HandlerFunc ) {
	glg.Log( "[cc] POST: ", path )
	http.HandleFunc(path, mwh.WrapPost( handler ) )
	postHandlers[path] = &handler
}

func GET( path string, handler http.HandlerFunc ) {
	glg.Log( "[cc] GET: ", path )
	http.HandleFunc(path, mwh.WrapGet( handler ) )
	getHandlers[path] = &handler
}
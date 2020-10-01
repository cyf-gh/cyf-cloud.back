package server

import (
	mwh "../middleware/helper"
	"net/http"
)

func resp(w* http.ResponseWriter, msg string) {
	(*w).Write([]byte(msg))
}

/// hello
func RootWelcomeGet(w http.ResponseWriter, r *http.Request) {
	resp( &w, string("🌸Welcome to api.cyf-cloud.cn!🌸") )
}

func cyfWelcomeGet(w http.ResponseWriter, r *http.Request) {
	resp( &w, string("<a href=\"https://www.cyf-cloud.cn\">") )
}

func pingGet( w http.ResponseWriter, r *http.Request ) {
	resp( &w, string("pong") )
}

func echoGet(w http.ResponseWriter, r *http.Request) {
	a := r.URL.Query()["a"][0]
	resp( &w, string(a) )
}

func InitBasicRequests() {
	/// ======================= video together ===========================
	http.HandleFunc("/", mwh.WrapGet( RootWelcomeGet ) )
	http.HandleFunc("/ping", mwh.WrapGet( pingGet ) )
	http.HandleFunc("/cyf", mwh.WrapGet( cyfWelcomeGet ) )
	http.HandleFunc("/echo", mwh.WrapGet( echoGet ) )
}
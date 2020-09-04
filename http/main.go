package httpMain

import (
	"net/http"

	ccV1 "../v1"
	vtMain "../v1/vt"
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

func echoGet(w http.ResponseWriter, r *http.Request) {
	a := r.URL.Query()["a"][0]
	resp( &w, string(a) )
}


func RunHttpServer( httpAddr string) {
	/// ======================= video together ===========================
	vtMain.Init()
	http.HandleFunc("/", RootWelcomeGet )
	http.HandleFunc("/cyf", cyfWelcomeGet )
	http.HandleFunc("/echo", echoGet )
	// http.HandleFunc( "/sync/guest",  )

	/// ======================= v1 ===========================
	ccV1.Init()
	http.HandleFunc( "/v1/donate/rank", ccV1.DonateRankGet )
	http.HandleFunc("/v1/util/mcdr/plg/script/generate", ccV1.GenerateScriptPost )
	http.HandleFunc("/v1/util/mcdr/plg/scripts", ccV1.FetchScriptGet )
	http.HandleFunc( "/v1/util/mcdr/plg/feed", ccV1.PluginListGet )
	http.ListenAndServe(httpAddr, nil)
}
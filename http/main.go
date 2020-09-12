package httpMain

import (
	"net/http"

	v1 "../v1"
	v1x1 "../v1x1"
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

// 路由应在Init函数中完成
func makeHttpRouter() {
	/// ======================= video together ===========================
	http.HandleFunc("/", RootWelcomeGet )
	http.HandleFunc("/cyf", cyfWelcomeGet )
	http.HandleFunc("/echo", echoGet )
	// http.HandleFunc( "/sync/guest",  )

	/// ======================= v1 ===========================
	v1.Init()
	http.HandleFunc( "/v1/donate/rank", v1.DonateRankGet )
	http.HandleFunc("/v1/util/mcdr/plg/script/generate", v1.GenerateScriptPost )
	http.HandleFunc("/v1/util/mcdr/plg/scripts", v1.FetchScriptGet )
	http.HandleFunc( "/v1/util/mcdr/plg/feed", v1.PluginListGet )

	/// ======================= v1x1 ===========================
	v1x1.Init()
}

// 创建所有的资源路由路径
// 路由路径为弱restful
func RunHttpServer( httpAddr string) {
	makeHttpRouter()

	http.ListenAndServe(httpAddr, nil)
}
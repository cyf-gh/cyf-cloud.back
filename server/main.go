package server

import (
	"github.com/kpango/glg"
	"net/http"

	v1 "../v1"
	v1x1 "../v1x1"
	v1x1http "../v1x1/http"
	orm "../v1x1/orm"
	security "../v1x1/security"
	mw "../middleware"
	mwu "../middleware/util"
)

// 路由应在Init函数中完成
func makeHttpRouter() {

	// server.HandleFunc( "/sync/guest",  )
	InitBasicRequests()

	/// ======================= v1 ===========================
	v1.Init()
	http.HandleFunc( "/v1/donate/rank", v1.DonateRankGet )
	http.HandleFunc("/v1/util/mcdr/plg/script/generate", v1.GenerateScriptPost )
	http.HandleFunc("/v1/util/mcdr/plg/scripts", v1.FetchScriptGet )
	http.HandleFunc( "/v1/util/mcdr/plg/feed", v1.PluginListGet )

	/// ======================= v1x1 ===========================
	v1x1.Init()
	v1x1http.Init()
}

func InitMiddlewares() {
	glg.Log( "middleware loading..." )
	mw.Register( mwu.LogUsedTime() )
	mw.Register( mwu.EnableCookie() )
	glg.Log( "middleware finished loading" )
}

// 创建所有的资源路由路径
// 路由路径为弱restful
func RunHttpServer( httpAddr string) {
	// 添加所有中间件
	InitMiddlewares()

	makeHttpRouter()
	// 初始化orm层
	orm.InitEngine("./.db/")

	// 初始化安全层
	security.Init()

	http.ListenAndServe(httpAddr, nil)
}
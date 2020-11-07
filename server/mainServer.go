package server

import (
	"../cc"
	"../config"
	mw "../middleware"
	mwu "../middleware/util"
	v1 "../v1"
	v1x1 "../v1x1"
	"../v1x1/cache"
	v1x1http "../v1x1/http"
	orm "../v1x1/orm"
	security "../v1x1/security"
	"github.com/kpango/glg"
	"net/http"
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
	mw.Register( cc.ErrorFetcher() )
	// mw.Register( mwu.EnableCookie() )
	// mw.Register( mwu.EnableAllowOrigin() )
	glg.Log( "middleware finished loading" )
}

// 创建所有的资源路由路径
// 路由路径为弱restful
func RunHttpServer( httpAddr string, rc config.RedisConfig, sqlitePath string ) {
	// 初始化orm层
	orm.InitEngine( sqlitePath )
	// 初始化键值对数据库
	cache.Init( rc )
	// 初始化安全层
	security.Init()
	// 添加所有中间件
	InitMiddlewares()
	// 部署所有路由
	makeHttpRouter()
	// 部署所有路由
	if e := cc.RegisterActions(); e != nil {
		panic( e )
	}
	http.ListenAndServe(httpAddr, nil)
}
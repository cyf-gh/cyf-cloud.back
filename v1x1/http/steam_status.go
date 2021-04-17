package http

import (
	"../../cc"
	err "../../cc/err"
	cfg "../../config"
	"../orm"
	"net/http"
)

func init() {
	cc.AddActionGroup( "/v1x1/steam", func( a cc.ActionGroup ) error {

		a.GET( "/recent", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			// ex := CheckSteamInfo( ap.R )
			return cc.HerOk()
		} )

		a.GET( "/achievement", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			return cc.HerOk()
		} )

		a.GET_DO( "/img", func( ap cc.ActionPackage ) ( cc.HttpErrReturn, cc.StatusCode ) {
			appid := ap.GetFormValue("appid")
			hash := ap.GetFormValue("hash")
			r, e := cc.GetByProxy("http://media.steampowered.com/steamcommunity/public/images/apps/" + appid + "/" + hash + ".jpg", cfg.ProxyAddr ); err.Assert( e )
			return cc.HerData( r )
		} )
		return nil
	} )
}

func CheckSteamInfo( r *http.Request ) *orm.AccountEx {
	ex, e := GetAccountExByAtk( r ); err.Assert( e )
	if ex.SteamId == "" || ex.SteamApi == "" {
		return nil
	}
	return ex
}
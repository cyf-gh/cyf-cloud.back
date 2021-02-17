package file

import (
	"../../cc"
	err "../../cc/err"
	"../../cc/err_code"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/user"
	"strings"
)

func GetFileRaw( w http.ResponseWriter, r *http.Request ) {
	var e error
	defer func() {
		if r := recover(); r != nil {
			cc.HttpReturn(&w, fmt.Sprint( r ), err_code.ERR_SYS, "", cc.MakeHER404 )
		}
	}()

	dir := r.URL.Query()["d"][0]

	if strings.Contains( dir, "..") {
		cc.HttpReturn( &w, "upper directory is not allowed", err_code.ERR_SECURITY, "", cc.MakeHER200 )
		return
	}

	current, e := user.Current()
	err.Assert( e )

	text, e := ioutil.ReadFile( current.HomeDir + "/.raw/" + dir)
	err.Assert( e )

	cc.HttpReturn( &w, "ok", err_code.ERR_OK, string(text), cc.MakeHER200 )
}

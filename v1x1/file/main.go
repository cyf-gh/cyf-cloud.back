package file

import (
	err "../err"
	err_code "../err_code"
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
			err.HttpReturn(&w, fmt.Sprint( r ), err_code.ERR_SYS, "", err_code.MakeHER404 )
		}
	}()

	dir := r.URL.Query()["d"][0]

	if strings.Contains( dir, "..") {
		err.HttpReturn( &w, "upper directory is not allowed", err_code.ERR_SECURITY, "", err_code.MakeHER200 )
		return
	}

	current, e := user.Current()
	err.CheckErr( e )

	text, e := ioutil.ReadFile( current.HomeDir + "/.raw/" + dir)
	err.CheckErr( e )

	err.HttpReturn( &w, "ok", err_code.ERR_OK, string(text), err_code.MakeHER200 )
}

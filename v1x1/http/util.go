package http

import (
	"io/ioutil"
	"net/http"
)

func Body2String( r * http.Request ) (string, error) {
	b, e := ioutil.ReadAll(r.Body)
	return string(b), e
}
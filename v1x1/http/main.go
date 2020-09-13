package http

import "net/http"

func Init() {
	http.HandleFunc("/v1x1/account/register", Register)
}
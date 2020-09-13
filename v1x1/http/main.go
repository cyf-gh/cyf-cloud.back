package http

import "net/http"

func Init() {
	http.HandleFunc("/v1x1/account/register", Register)
	http.HandleFunc("/v1x1/account/login", Login)
}

func enableCookies(w *http.ResponseWriter) {
	// (*w).Header().Set("Access-Control-Allow-Origin", "https://localhost:8887")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}
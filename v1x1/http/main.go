package http

import "net/http"

func Init() {
	AccessTokens = make(map[string]int64)
	// 方便测试的令牌，cyfhaoshuai
	AccessTokens["cyhaoshuaicyfhaoshuaicyfhaoshuaicyfhaoshuai"] = 1

	http.HandleFunc("/v1x1/account/register", Register)
	http.HandleFunc("/v1x1/account/login", Login)

	http.HandleFunc("/v1x1/post/create", NewPost)
	http.HandleFunc("/v1x1/post/modify", ModifyPost)
	http.HandleFunc("/v1x1/post/modifyNT", ModifiyPostNoText)
}

func enableCookies(w *http.ResponseWriter) {
	// (*w).Header().Set("Access-Control-Allow-Origin", "https://localhost:8887")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}
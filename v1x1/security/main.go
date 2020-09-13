package security

import "net/http"

func Init() {
	http.HandleFunc("/v1x1/security/captcha", GenerateCaptcha)
}
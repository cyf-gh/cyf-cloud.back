package security

import (
	"net/http"
	mwh "../../middleware/helper"
)

func Init() {
	http.HandleFunc("/v1x1/security/captcha", mwh.WrapGet( GenerateCaptcha ) )
}
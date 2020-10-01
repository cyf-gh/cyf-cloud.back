package security

import (
	mwh "../../middleware/helper"
	"net/http"
)

func Init() {
	http.HandleFunc("/v1x1/security/captcha", mwh.WrapGet( GenerateCaptcha ) )
}
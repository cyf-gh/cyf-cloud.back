package v1x1

import (
	"net/http"
	"./file"
	mwh "../middleware/helper"
)

func Init() {
	http.HandleFunc( "/v1x1/raw", mwh.WrapGet( file.GetFileRaw ) )
}


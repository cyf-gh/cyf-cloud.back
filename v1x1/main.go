package v1x1

import (
	mwh "../middleware/helper"
	"./file"
	"net/http"
)

func Init() {
	http.HandleFunc( "/v1x1/raw", mwh.WrapGet( file.GetFileRaw ) )
}


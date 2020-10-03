package v1x1

import (
	mwh "../middleware/helper"
	"./cli"
	"./file"
	"net/http"
)

func Init() {
	cli.Init()

	http.HandleFunc( "/v1x1/raw", mwh.WrapGet( file.GetFileRaw ) )
}


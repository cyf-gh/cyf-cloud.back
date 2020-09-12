package v1x1

import (
	"net/http"
	"./file"
)

func Init() {
	http.HandleFunc( "/v1x1/raw", file.GetFileRaw )
}


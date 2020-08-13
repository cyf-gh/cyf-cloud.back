package v1

import (
	"io/ioutil"
	"net/http"
)

func resp(w* http.ResponseWriter, msg string) {
	(*w).Write([]byte(msg))
}

func DonateRankGet(w http.ResponseWriter, r *http.Request) {
	json, err := ioutil.ReadFile("./.static/donate-rank.json")
	if err != nil {
		panic(err)
	}
	resp( &w, string(json) )
}

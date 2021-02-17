package cc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Get( url string ) ( respStr string, e error ) {
	resp, e := http.Get( url ); if e != nil { return }
	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body); if e != nil { return }
	respStr = string( body )
	return
}

func GetJ( url string, v interface{} ) ( e error ) {
	resp, e := http.Get( url ); if e != nil { return }
	defer resp.Body.Close()
	b, e := ioutil.ReadAll( resp.Body ); if e != nil { return e }
	e = json.Unmarshal( b, v ); if e != nil { return e }
	return
}

func PostJ( u string, body interface{},v interface{} ) ( e error ) {
	client := &http.Client{}
	jn, e := json.Marshal( body ); if e != nil { return e }
	b := bytes.NewBuffer(jn)
	req, e := http.NewRequest("POST", u, b ); if e != nil { return e }
	resp, e := client.Do(req); if e != nil { return e }
	defer resp.Body.Close()
	bbb, e := ioutil.ReadAll( resp.Body ); if e != nil { return e }
	e = json.Unmarshal( bbb, v ); if e != nil { return e }
	return
}

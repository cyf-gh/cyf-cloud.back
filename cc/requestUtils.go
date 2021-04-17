package cc

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func GetByProxy( webUrl, proxyUrl string ) ( respStr string, e error ) {
	proxy, e := url.Parse(proxyUrl)
	tr := &http.Transport{
		Proxy: http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 15,
	}
	resp, e := client.Get(webUrl)

	body, e := ioutil.ReadAll(resp.Body)
	respStr = string(body)
	defer resp.Body.Close()
	return
}

func PostByProxy( webUrl string, body interface{}, v interface{}, proxyUrl string) ( e error ) {
	proxy, _ := url.Parse(proxyUrl)
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 15,
	}
	jn, e := json.Marshal( body ); if e != nil { return e }
	b := bytes.NewBuffer(jn)
	req, e := http.NewRequest("POST", webUrl, b ); if e != nil { return e }
	resp, e := client.Do(req); if e != nil { return e }

	defer resp.Body.Close()
	bbb, e := ioutil.ReadAll( resp.Body ); if e != nil { return e }
	e = json.Unmarshal( bbb, v ); if e != nil { return e }
	return
}

// http get request
func Get( url string ) ( respStr string, e error ) {
	resp, e := http.Get( url ); if e != nil { return }
	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body); if e != nil { return }
	respStr = string( body )
	return
}

// http get request with json
func GetJ( url string, v interface{} ) ( e error ) {
	resp, e := http.Get( url ); if e != nil { return }
	defer resp.Body.Close()
	b, e := ioutil.ReadAll( resp.Body ); if e != nil { return e }
	e = json.Unmarshal( b, v ); if e != nil { return e }
	return
}

// http post request with json
func PostJ( u string, body interface{}, v interface{} ) ( e error ) {
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

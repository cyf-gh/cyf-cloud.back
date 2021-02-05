package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

// https://mholt.github.io/json-to-go/
// \sa server_dm1.go
// this struct is auto generated
type DM1Config struct {
	Ignore []struct {
		Name string `json:"name"`
	} `json:"ignore"`
}

var (
	DM1 DM1Config
)

func init() {
	jf, e := os.Open("./server_dm1.json"); defer jf.Close()
	if e != nil {
		panic( e )
	}
	bytes, _ := ioutil.ReadAll( jf )
	json.Unmarshal( bytes, &DM1 )
	println("DM1 config loaded...")
}

// 配置文件中是否标明该文件需要忽略
func IsIgnoreResource( path string ) bool {
	for _, i := range DM1.Ignore {
		if strings.Contains( path, i.Name ) {
			return true
		}
	}
	return false
}
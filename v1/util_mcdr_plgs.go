package v1

import (
	"encoding/json"
	"github.com/kpango/glg"
	"io/ioutil"
	"net/http"
	"strings"
)

func InitUtilMcdrPlgs() {
	PlgPlugins = make(map[string]string)
}

type PlgsModel struct {
	Name string `json:"name"`
	Git string `json:"git"`
	Desc string `json:"desc"`
}
var PlgPlugins map[string]string

func GenerateScriptPost( w http.ResponseWriter, r *http.Request ) {
	uuid := r.URL.Query()["u"][0]
	sav := r.URL.Query()["clear"][0]
	path := r.URL.Query()["path"][0]
	var plgs[] PlgsModel

	b, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal( b, &plgs )
	if err != nil {
		glg.Log(err)
		panic(err)
		resp( &w, "INTERVAL_ERR")
		return
	}
	if path[len(path)-1] == '/' {
		path = strings.TrimRight(path, "/")
	}
	scriptClone := ""
	for  _, plg := range plgs {
		scriptClone += "git clone "+ plg.Git + "\n"
	}
	clearPluginFolder := ""

	if sav == "false" {
		clearPluginFolder = "rm -r $mcdrPath/plugins/*"
	}
	script := `
echo [Info] MCDR Plugin Helper v0.1 by cyf

mcdrPath=`+ path +`
if [ ! -d $mcdrPath/plugins ]; then
    echo "[Error] please check the MCDR directory"
    return
fi
### clone所有插件 ###
if [ ! -d ./.plg ]; then
    mkdir ./plg
fi
cd ./plg
`+ scriptClone + `
cd ../
### 确认删除原来的插件 ###	
echo $mcdrPath
echo [Info] Delete old plugin folder[Require confirm]
`+clearPluginFolder+`
### 开始复制插件 ###	
echo [Info] Copying plugins...
echo ==========
ls ./plg | while read line
do
    plgPath=./plg/$line
    if [ -d $plgPath ]; then
        cp -r $plgPath/*.py  $mcdrPath/plugins/
        echo [Success] $plgPath.py copied!!!
    fi
done
echo ==========
echo [Info] Done.
`
	PlgPlugins[uuid] = script
	resp( &w, "OK" )
}

func FetchScriptGet( w http.ResponseWriter, r *http.Request ) {
	uuid := r.URL.Query()["u"][0]
	resp( &w, PlgPlugins[uuid] )
}

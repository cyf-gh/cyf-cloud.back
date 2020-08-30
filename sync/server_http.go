package vt_sync

import (
	ccV1 "../v1"
	vtLobby "../lobby"
	"encoding/json"
	"github.com/kpango/glg"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)
var Lock *sync.Mutex

func resp(w* http.ResponseWriter, msg string) {
	(*w).Write([]byte(msg))
}

func CheckUserStatus(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query()["username"][0]
	resp( &w, vtLobby.CheckUserStatus( username, vtLobby.Lobbies ) )
}

func ExitLobbyGet(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query()["username"][0]
	var res string
	Lock.Lock()
	// LOBBY_DELETED
	// LOBBY_EXIT
	// NO_SUCH_LOBBY
	res, vtLobby.Lobbies = vtLobby.ExitLobby(username, vtLobby.Lobbies)
	Lock.Unlock()
	resp( &w, res )
}

func QueryLobbiesGet(w http.ResponseWriter, r *http.Request) {
	lobbyNames := ""
	for _, lob := range vtLobby.Lobbies {
		lobbyNames += lob.Name + ","
	}
	lobbyNames = strings.TrimSuffix( lobbyNames, "," )
	glg.Log(lobbyNames)
	resp( &w, lobbyNames )
}

func EnterlobbyGet(w http.ResponseWriter, r *http.Request){
	userName := r.URL.Query()["username"][0]
	lobbyName := r.URL.Query()["lobbyname"][0]
	passwd := r.URL.Query()["passwd"][0]
	for _, lob:= range vtLobby.Lobbies {
		if lob.Name == lobbyName {
			if passwd != lob.Password {
				resp(&w,"PASSWORD_INCORRECT")
				return
			}
			lob.Viewers = append(lob.Viewers, vtLobby.VTViewer{
				Name:     userName,
				Location: "0",
				IsHost:   false,
				IsPause:  false,
			})
			glg.Info("Guest[" + userName + "]has entered lobby[" + lobbyName + "]")
			resp(&w,"OK")
			return
		}
	}
	resp(&w,"NO_SUCH_LOBBY")
}

func CreatelobbyGet(w http.ResponseWriter, r *http.Request){
	hostName := r.URL.Query()["hostname"][0]
	lobbyName := r.URL.Query()["lobbyname"][0]
	passwd := r.URL.Query()["passwd"][0]

	// 添加新房间
	newLobby := &vtLobby.VTLobby{
		Name:      lobbyName,
		Password:  passwd,
		Viewers: []vtLobby.VTViewer{
			// 添加房主
			vtLobby.VTViewer{
				Name:     hostName,
				Location: "",
				IsHost:   true,
				IsPause:  true,
			},
		},
		MaxOffset: 2,
		VideoIndex: 0,
		VideoLs: "",
		LastUpdateTime: time.Now(),
	}
	if vtLobby.IsSameNameLobbyExist( newLobby.Name, vtLobby.Lobbies ) {
		resp(&w,"LOBBY_EXISTED")
		glg.Info(newLobby.Name + " is already exist but someone still wants to borrow one.")
		return
	}
	Lock.Lock()
	vtLobby.Lobbies = append(vtLobby.Lobbies, newLobby)
	Lock.Unlock()

	resp(&w,"OK")

	glg.Info("lobby " + newLobby.Name + " created!")
	for _, l := range vtLobby.Lobbies {
		glg.Info(*l)
	}
}

func PingGet(w http.ResponseWriter, r *http.Request) {
	resp( &w, "OK")
}

func SendVideoInfoPost(w http.ResponseWriter, r *http.Request) {
	hostName := r.URL.Query()["hostname"][0]
	lobby, i, _ := vtLobby.FindLobbyByHost( hostName, vtLobby.Lobbies )
	if lobby == nil {
		resp( &w, "NO_AUTH")
		return
	}
	// 将数据同步至房间
	var videoDesc vtLobby.VTVideoDesc
	err := json.NewDecoder(r.Body).Decode( &videoDesc )
	if err != nil {
		panic(err)
		resp( &w, "INTERVAL_ERR")
		return
	}
	// web 版本ls即为src
	lobby.VideoLs = videoDesc.Ls
	lobby.VideoIndex = videoDesc.Index
	lobby.Md5 = videoDesc.Md5
	vtLobby.Lobbies[i] = lobby
	glg.Info("Video[" + videoDesc.Ls + "]\n P:[" + strconv.Itoa( videoDesc.Index ) + "]")
	resp( &w, "OK")
}
func UserWhereGet (w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query()["username"][0]
	lobby, _, ishost := vtLobby.FindLobbyByUser( username, vtLobby.Lobbies )
	if lobby == nil {
		resp( &w, "IDLE" )
		return
	}
	var host string
	if ishost {
		host = "HOST"
	} else {
		host = "GUEST"
	}
	resp( &w, lobby.Name + "," + host + "," + lobby.Password )
}

/// sync ///
func SendSyncHostGet(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query()["name"][0]
	location := r.URL.Query()["location"][0]
	pause := r.URL.Query()["ispause"][0] // p, s
	part :=  r.URL.Query()["p"][0]
	lobby, i, _ := vtLobby.FindLobbyByHost( name, vtLobby.Lobbies )
	if lobby == nil {
		return
	}
	lobby.IsPause = pause
	lobby.Location = location
	lobby.VideoIndex, _ = strconv.Atoi( part )
	glg.Log("========SYNC========")
	glg.Log("[HOST]"+ name)
	glg.Log("[LOCATION]"+ location)
	glg.Log("[IS PAUSE]" + pause)
	glg.Log("[PART]" + part)
	glg.Log("========SYNC========")
	Lock.Lock()
	vtLobby.Lobbies[i] = lobby
	Lock.Unlock()
}

func SendSyncGuestGet(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query()["name"][0]
	lb, i, ishost := vtLobby.FindLobbyByUser( name, vtLobby.Lobbies )
	if i == -1 || ishost || lb == nil {
		resp( &w, "ERR" )
	}
	// md5,p/s,location,part
	resp( &w, lb.Md5 + "," + lb.IsPause + "," + lb.Location + "," + strconv.Itoa( lb.VideoIndex ) )
}

func GetCurrentVideoDesc(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query()["name"][0]
	t := r.URL.Query()["t"][0]
	lb, i, ishost := vtLobby.FindLobbyByUser( name, vtLobby.Lobbies )
	if i == -1 || ishost || lb == nil {
		resp( &w, "ERR" )
		return
	}
	if t == "web" {
		resp( &w, lb.VideoLs )
		return
	}
	// ls,index
	resp( &w, lb.VideoLs + "`" + strconv.Itoa( lb.VideoIndex ) )
}

func GetUserStatus(w http.ResponseWriter, r *http.Request) {
	lobbyname := r.URL.Query()["lobbyname"][0]
	exists, lob := vtLobby.IsLobbyExist(lobbyname, vtLobby.Lobbies)
	if !exists {
		resp( &w, "NO_SUCH_LOBBY" )
		return
	}
	jsons, errs := json.Marshal(lob.Viewers)
	if errs != nil {
		glg.Log(errs)
		return
	}
	resp( &w, string(jsons) )
}

/// hello
func RootWelcomeGet(w http.ResponseWriter, r *http.Request) {
	resp( &w, string("Welcome to api.cyf-cloud.cn!\ncyf-cloud's backend API interfaces.") )
}

func cyfWelcomeGet(w http.ResponseWriter, r *http.Request) {
	resp( &w, string("Who is cyf?") )
}

func echoGet(w http.ResponseWriter, r *http.Request) {
	a := r.URL.Query()["a"][0]
	resp( &w, string(a) )
}


type PlgsModel struct {
	Name string `json:"name"`
	Git string `json:"git"`
	Desc string `json:"desc"`
}
var PlgPlugins map[string]string

func generateScriptPost( w http.ResponseWriter, r *http.Request ) {
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

func mcdrPlgScriptGet( w http.ResponseWriter, r *http.Request ) {
	uuid := r.URL.Query()["u"][0]
	resp( &w, PlgPlugins[uuid] )
}

func RunHttpSyncServer( httpAddr string, lock *sync.Mutex ) {
	Lock = lock
	PlgPlugins = make(map[string]string)
	http.HandleFunc("/", RootWelcomeGet )
	http.HandleFunc("/v1/cyf", cyfWelcomeGet )
	http.HandleFunc("/v1/echo", echoGet )
	http.HandleFunc("/v1/vt/ping", PingGet)
	http.HandleFunc("/v1/vt/lobby/enter", EnterlobbyGet)
	http.HandleFunc("/v1/vt/lobby/create", CreatelobbyGet)
	http.HandleFunc("/v1/vt/lobby/exit", ExitLobbyGet)
	http.HandleFunc( "/v1/vt/lobby/update/videodesc", SendVideoInfoPost )
	http.HandleFunc("/v1/vt/lobbies", QueryLobbiesGet)
	http.HandleFunc( "/v1/vt/user/status", CheckUserStatus )
	http.HandleFunc( "/v1/vt/user/where", UserWhereGet )
	http.HandleFunc( "/v1/vt/sync/host", SendSyncHostGet )
	http.HandleFunc( "/v1/vt/sync/guest", SendSyncGuestGet )
	http.HandleFunc( "/v1/vt/lobby/users/status", GetUserStatus )
	http.HandleFunc( "/v1/vt/lobby/videodesc", GetCurrentVideoDesc )
	// http.HandleFunc( "/sync/guest",  )

	http.HandleFunc( "/v1/donate/rank", ccV1.DonateRankGet )

	http.HandleFunc("/v1/util/mcdr/plg/script/generate", generateScriptPost )
	http.HandleFunc("/v1/util/mcdr/plg/scripts", mcdrPlgScriptGet )
	http.HandleFunc( "/v1/util/mcdr/plg/feed", ccV1.PluginListGet )
	http.ListenAndServe(httpAddr,nil)
}
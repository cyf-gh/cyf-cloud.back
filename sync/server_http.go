package vt_sync

import (
	ccV1 "../v1"
	vtLobby "../lobby"
	"encoding/json"
	"github.com/kpango/glg"
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
	lb, i, ishost := vtLobby.FindLobbyByUser( name, vtLobby.Lobbies )
	if i == -1 || ishost || lb == nil {
		resp( &w, "ERR" )
	}
	// ls,index
	resp( &w, lb.VideoLs + "`" + strconv.Itoa( lb.VideoIndex ) )
}

func RunHttpSyncServer( httpAddr string, lock *sync.Mutex ) {
	Lock = lock
	http.HandleFunc("/ping", PingGet)
	http.HandleFunc("/lobby/enter", EnterlobbyGet)
	http.HandleFunc("/lobby/create", CreatelobbyGet)
	http.HandleFunc("/lobby/exit", ExitLobbyGet)
	http.HandleFunc("/lobby", ExitLobbyGet)
	http.HandleFunc( "/lobby/update/videodesc", SendVideoInfoPost )
	http.HandleFunc("/lobbies", QueryLobbiesGet)
	http.HandleFunc( "/user/status", CheckUserStatus )
	http.HandleFunc( "/user/where", UserWhereGet )
	http.HandleFunc( "/sync/host", SendSyncHostGet )
	http.HandleFunc( "/sync/guest", SendSyncGuestGet )
	http.HandleFunc( "/lobby/videodesc", GetCurrentVideoDesc )
	// http.HandleFunc( "/sync/guest",  )

	http.HandleFunc( "/v1/donate/rank", ccV1.DonateRankGet )

	http.ListenAndServe(httpAddr,nil)
}
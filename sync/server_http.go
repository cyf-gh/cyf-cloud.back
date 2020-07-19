package vt_sync

import (
	vtLobby "../lobby"
	"github.com/kpango/glg"
	"net/http"
	"strings"
	"sync"
	"time"
)
var Lock *sync.Mutex

func resp(w* http.ResponseWriter, msg string) {
	(*w).Write([]byte(msg))
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
			resp(&w,"OK")
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

func RunHttpSyncServer( httpAddr string, lock *sync.Mutex ) {
	Lock = lock
	http.HandleFunc("/ping", PingGet)
	http.HandleFunc("/lobby/enter", EnterlobbyGet)
	http.HandleFunc("/lobby/create", CreatelobbyGet)
	http.HandleFunc("/lobbies", QueryLobbiesGet)
	http.ListenAndServe(httpAddr,nil)
}
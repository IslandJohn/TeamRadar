package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Login struct {
	User     string
	Password string
	sync     sync.Mutex
}
type User struct {
	Id          string
	DisplayName string
	Url         string
	ImageUrl    string
}
type RoomUser struct {
	RoomId       int
	User         *User
	LastActivity string
	JoinedDate   string
	IsOnline     bool
}
type RoomUsers struct {
	Count int
	Value []*RoomUser
}
type Room struct {
	Id                      int
	Name                    string
	Description             string
	LastActivity            string
	CreatedBy               *User
	CreatedDate             string
	HasAdminPermissions     bool
	HasReadWritePermissions bool
}
type Rooms struct {
	Count int
	Value []*Room
}

var apiClient http.RoundTripper = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
}
var apiBase string = os.Args[1]
var apiEndpoint []string = []string{
	apiBase + "/_apis/chat/rooms?api-version=1.0",
	apiBase + "/_apis/chat/rooms/%d/users?api-version=1.0",
	apiBase + "/_apis/chat/rooms/%d/users/%s?api-version=1.0",
	apiBase + "/_apis/chat/rooms/%d/messages?api-version=1.0",
	apiBase + "/_apis/chat/rooms/%d/messages?$filter=PostedTime ge %s&api-version=1.0",
}

func main() {
	login := Login{
		User:     os.Args[2],
		Password: os.Args[3],
	}
	//log.Printf("main login=%s", login)

	listener := make(chan interface{}) // routines send updates here
	quitter := make(chan interface{})  // close to have all routines return

	// go pollClient() // get input
	go pollRooms(5*time.Second, &login, &listener, &quitter)

	for l := range listener {
		fmt.Println(l)
	}
}

// requests are mutually exclusive on login
// so login can be updated on errors
// without others triggering same
func makeRequest(verb string, url string, body string, code int, login *Login, listener *chan interface{}) ([]byte, error) {
	//log.Printf("makeRequest verb=%s url=%s body=%s code=%d", verb, url, body, code)
	login.sync.Lock()
	defer login.sync.Unlock()
	
	request, err := http.NewRequest(verb, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(login.User, login.Password)

	response, err := apiClient.RoundTrip(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != code {
		return nil, errors.New(response.Status)
	}

	return ioutil.ReadAll(response.Body)
}

func pollRooms(d time.Duration, login *Login, listener *chan interface{}, quitter *chan interface{}) {
	//log.Printf("pollRooms d=%d", d)
	roomMap := make(map[int]*Room)
	roomTicker := time.NewTicker(d)
	defer roomTicker.Stop()

	for {
		select {
		case <-roomTicker.C:
			// api
			body, err := makeRequest("GET", apiEndpoint[0], "", 200, login, listener)
			if err != nil {
				log.Printf("pollRooms err=%s", err)
				continue
			}

			var rooms Rooms // or interface{} for generic
			err = json.Unmarshal(body, &rooms)
			if err != nil {
				log.Printf("pollRooms err=%s", err)
				continue
			}
			//log.Printf("pollRooms rooms=%s", rooms)

			// added
			newRoomMap := make(map[int]*Room)
			for _, room := range rooms.Value {
				newRoomMap[room.Id] = room
				if _, ok := roomMap[room.Id]; ok {
					continue
				}

				json, err := json.Marshal(room)
				if err != nil {
					log.Printf("pollRooms err=%s", err)
					continue
				}

				*listener <- fmt.Sprintf("room add %d %s", room.Id, json)
				roomMap[room.Id] = room

				// need to track quitter for room users polling
				go pollRoomUsers(3*time.Second, room, login, listener, quitter)
			}

			// removed
			for _, room := range roomMap {
				if _, ok := newRoomMap[room.Id]; ok {
					continue
				}

				json, err := json.Marshal(*room)
				if err != nil {
					log.Printf("pollRooms err=%s", err)
					continue
				}

				*listener <- fmt.Sprintf("room remove %d %s", room.Id, json)
				delete(roomMap, room.Id)

				// need to clean up room users polling
			}
		case <-*quitter:
			return
		}

		// if 401 or "auth" header, need U/P from speaker, block everywhere?
		// for added rooms start ticker, for removed rooms stop ticker
		// need to signal "done"?
	}
}

func pollRoomUsers(d time.Duration, room *Room, login *Login, listener *chan interface{}, quitter *chan interface{}) {
	//log.Printf("pollRoomUsers d=%d room=%s", d, *room)
	userMap := make(map[string]*RoomUser)
	userTicker := time.NewTicker(d)
	defer userTicker.Stop()

	for {
		select {
		case <-userTicker.C:
			// api
			body, err := makeRequest("GET", fmt.Sprintf(apiEndpoint[1], room.Id), "", 200, login, listener)
			if err != nil {
				continue
			}

			var users RoomUsers
			err = json.Unmarshal(body, &users)
			if err != nil {
				log.Printf("pollRoomUsers err=%s", err)
				continue
			}

			// added
			newUserMap := make(map[string]*RoomUser)
			for _, user := range users.Value {
				newUserMap[user.User.Id] = user
				if _, ok := userMap[user.User.Id]; ok {
					continue
				}

				json, err := json.Marshal(user)
				if err != nil {
					log.Printf("pollRoomUsers err=%s", err)
					continue
				}

				*listener <- fmt.Sprintf("user add %d %s %s", user.RoomId, user.User.Id, json)
				userMap[user.User.Id] = user
			}

			// removed
			for _, user := range userMap {
				if _, ok := newUserMap[user.User.Id]; ok {
					continue
				}

				json, err := json.Marshal(*user)
				if err != nil {
					log.Printf("pollRoomUsers err=%s", err)
					continue
				}

				*listener <- fmt.Sprintf("user remove %d %s %s", user.RoomId, user.User.Id, json)
				delete(userMap, user.User.Id)
			}

			// changed
			for id, user := range newUserMap {
				if newUserMap[id].IsOnline != userMap[id].IsOnline {
					json, err := json.Marshal(*user)
					if err != nil {
						log.Printf("pollRoomUsers err=%s", err)
						continue
					}

					userMap[id] = newUserMap[id]
					*listener <- fmt.Sprintf("user change %d %s %s", user.RoomId, user.User.Id, json)
				}
			}
		case <-*quitter:
			return
		}

		// if 401 or "auth" header, need U/P from speaker, block everywhere?
		// for added rooms start ticker, for removed rooms stop ticker
		// need to signal "done"?
	}
}

func pollMessages(d time.Duration) {
}

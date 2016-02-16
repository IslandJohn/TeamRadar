package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	//	"k8s.io/kubernetes/pkg/util/sets"
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
type Users struct {
	Count int
	Value []User
}
type Room struct {
	Id                      int
	Name                    string
	Description             string
	LastActivity            string
	CreatedBy               User
	CreatedDate             string
	HasAdminPermissions     bool
	HasReadWritePermissions bool
}
type Rooms struct {
	Count int
	Value []Room
}

var apiClient http.RoundTripper = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
}
var apiBase string = os.Args[1]
var apiEndpoint []string = []string{
	apiBase + "/_apis/chat/rooms?api-version=1.0",
	apiBase + "/_apis/chat/rooms/%s/users?api-version=1.0",
	apiBase + "/_apis/chat/rooms/%s/users/%s?api-version=1.0",
	apiBase + "/_apis/chat/rooms/%s/messages?api-version=1.0",
	apiBase + "/_apis/chat/rooms/%s/messages?$filter=PostedTime ge %s&api-version=1.0",
}

func main() {
	login := Login{
		User:     os.Args[2],
		Password: os.Args[3],
	}
	listener := make(chan string)

	go pollRooms(5*time.Second, &login, &listener)

	for l := range listener {
		fmt.Println(l)
	}
}

// requests are mutually exclusive on login
// so login can be updated on errors
// without others triggering same
func makeRequest(login *Login, verb string, url string, body string, code int, listener *chan string) ([]byte, error) {
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

func pollRooms(d time.Duration, login *Login, listener *chan string) {
	//roomMap := make(map[string]string)
	roomTicker := time.NewTicker(d)
	defer roomTicker.Stop()

	for _ = range roomTicker.C {
		body, err := makeRequest(login, "GET", apiEndpoint[0], "", 200, listener)
		if err != nil {
			*listener <- fmt.Sprintf("error pollRooms %s", err)
			continue
		}

		var rooms Rooms // or interface{} for generic
		err = json.Unmarshal(body, &rooms)
		if err != nil {
			*listener <- fmt.Sprintf("error pollRooms %s", err)
			continue
		}

		// if 401 or "auth" header, need U/P from speaker, block everywhere?
		// determine room added, removed
		// for added rooms start ticker, for removed rooms stop ticker
		// need to signal "done"?
	}
}

func pollUsers(d time.Duration) {
}

func pollMessages(d time.Duration) {
}

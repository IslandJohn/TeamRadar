/*
Copyright 2016 IslandJohn and the TeamRadar Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tfs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IslandJohn/TeamRadar/Go/teamradar/rest"
	"net/http"
	"strings"
)

// TFS JSON
type User struct {
	Id          string
	DisplayName string
	Url         string
	ImageUrl    string
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
type RoomMessage struct {
	Id           int
	Content      string
	MessageType  string
	PostedTime   string
	PostedRoomId int
	PostedBy     *User
}
type RoomMessages struct {
	Count int
	Value []*RoomMessage
}

// TFS URL
var restEndpoint []string = []string{
	"%s/_apis/projects?api-version=1.0",
	"%s/_apis/chat/rooms?api-version=1.0",
	"%s/_apis/chat/rooms/%d/users?api-version=1.0",
	"%s/_apis/chat/rooms/%d/users/%s?api-version={version}",
	"%s/_apis/chat/rooms/%d/messages?api-version=1.0",
	"%s/_apis/chat/rooms/%d/messages?$filter=PostedTime+ge+%s&api-version=1.0",
}

// TFS API
type Account struct {
	UserId    string
	LoginUser string
}
type Api struct {
	restBase     string
	restClient   *rest.Client
	LoginAccount *Account
}

// create a new TFS API instance, verifying access
func NewApi(url string, user string, password string) (*Api, error) {
	api := Api{
		restBase:   url,
		restClient: rest.NewClient(),
	}
	api.restClient.SetLogin(user, password)

	header, _, err := api.restClient.MakeRequest("GET", fmt.Sprintf(restEndpoint[0], api.restBase), "", http.StatusOK)
	if err != nil {
		return nil, err
	}

	userdata, ok := header["X-Vss-Userdata"]
	if !ok || len(user) <= 0 {
		return nil, errors.New("Missing header X-VSS-UserData")
	}

	fields := strings.SplitN(userdata[0], ":", 2)
	if len(fields) != 2 {
		return nil, errors.New("Invalid header X-VSS-UserData")
	}

	api.LoginAccount = &Account{
		fields[0],
		fields[1],
	}
	return &api, nil
}

// get the list of rooms
func (a *Api) GetRooms() (*Rooms, error) {
	_, body, err := a.restClient.MakeRequest("GET", fmt.Sprintf(restEndpoint[1], a.restBase), "", http.StatusOK)
	if err != nil {
		return nil, err
	}

	var rooms Rooms
	err = json.Unmarshal(body, &rooms)
	if err != nil {
		return nil, err
	}

	return &rooms, nil
}

// get the list of users in a room
func (a *Api) GetRoomUsers(room *Room) (*RoomUsers, error) {
	_, body, err := a.restClient.MakeRequest("GET", fmt.Sprintf(restEndpoint[2], a.restBase, room.Id), "", http.StatusOK)
	if err != nil {
		return nil, err
	}

	var users RoomUsers
	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, err
	}

	return &users, nil
}

// join a room
func (a *Api) JoinRoom(room *Room) error {
	json, err := json.Marshal(a.LoginAccount)
	if err != nil {
		return err
	}

	_, _, err = a.restClient.MakeRequest("PUT", fmt.Sprintf(restEndpoint[3], a.restBase, room.Id, a.LoginAccount.UserId), string(json), http.StatusNoContent)
	if err != nil {
		return err
	}

	return nil
}

// leave a room
func (a *Api) LeaveRoom(room *Room) error {
	_, _, err := a.restClient.MakeRequest("DELETE", fmt.Sprintf(restEndpoint[3], a.restBase, room.Id, a.LoginAccount.UserId), "", http.StatusNoContent)
	if err != nil {
		return err
	}

	return nil
}

// create a message in a room
func (a *Api) SendRoomMessage(room *Room, msg string) (*RoomMessage, error) {
	return nil, nil
}

// get the list of messages in a room since date
func (a *Api) GetRoomMessages(room *Room, date string) (*RoomMessages, error) {
	_, body, err := a.restClient.MakeRequest("GET", fmt.Sprintf(restEndpoint[5], a.restBase, room.Id, date), "", http.StatusOK)
	if err != nil {
		return nil, err
	}

	var messages RoomMessages
	err = json.Unmarshal(body, &messages)
	if err != nil {
		return nil, err
	}

	return &messages, nil
}

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
	"fmt"
	"github.com/IslandJohn/TeamRadar/Go/teamradar/rest"
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
	"%s/_apis/chat/rooms?api-version=1.0",
	"%s/_apis/chat/rooms/%d/users?api-version=1.0",
	"%s/_apis/chat/rooms/%d/users/%s?api-version=1.0",
	"%s/_apis/chat/rooms/%d/messages?api-version=1.0",
	"%s/_apis/chat/rooms/%d/messages?$filter=PostedTime+ge+%s&api-version=1.0",
}

// TFS API
type Api struct {
	restBase   string
	restClient *rest.Client
}

func NewApi(url string, user string, password string) *Api {
	api := Api{
		restBase:   url,
		restClient: rest.NewClient(),
	}
	api.SetLogin(user, password)

	return &api
}

// set login information
func (a *Api) SetLogin(user string, password string) {
	a.restClient.SetLogin(user, password)
}

// get the list of rooms
func (a *Api) GetRooms() (*Rooms, error) {
	body, err := a.restClient.MakeRequest("GET", fmt.Sprintf(restEndpoint[0], a.restBase), "", 200)
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
	body, err := a.restClient.MakeRequest("GET", fmt.Sprintf(restEndpoint[1], a.restBase, room.Id), "", 200)
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

// get the list of messages in a room since date
func (a *Api) GetRoomMessages(room *Room, date string) (*RoomMessages, error) {
	body, err := a.restClient.MakeRequest("GET", fmt.Sprintf(restEndpoint[4], a.restBase, room.Id, date), "", 200)
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

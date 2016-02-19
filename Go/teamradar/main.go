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

package main

import (
	"encoding/json"
	"fmt"
	"github.com/IslandJohn/TeamRadar/Go/teamradar/tfs"
	"github.com/IslandJohn/TeamRadar/Go/teamradar/trace"
	"os"
	"time"
)

// main runs the routines and collects
func main() {
	tfsApi := tfs.NewApi(os.Args[1], os.Args[2], os.Args[3])
	recv := make(chan interface{}) // routines send updates here
	quit := make(chan interface{}) // close this to have all routines return

	go pollTfsRooms(tfsApi, 1*time.Second, 60*time.Second, &recv, &quit)

	for l := range recv {
		fmt.Println(l)
	}
}

// routine to poll room information at variable intervals
func pollTfsRooms(tfsApi *tfs.Api, min time.Duration, max time.Duration, send *chan interface{}, recv *chan interface{}) {
	delay := min
	roomMap := make(map[int]*tfs.Room)
	roomQuit := make(map[int]*chan interface{})

	timer := time.NewTimer(delay)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			delay = delay * 2
			if delay > max {
				delay = max
			}

			rooms, err := tfsApi.GetRooms()
			if err != nil {
				trace.Log(err)
				continue
			}

			// added
			newRoomMap := make(map[int]*tfs.Room)
			for _, room := range rooms.Value {
				newRoomMap[room.Id] = room
				if _, ok := roomMap[room.Id]; ok {
					continue
				}

				json, err := json.Marshal(room)
				if err != nil {
					trace.Log(err)
					continue
				}

				*send <- fmt.Sprintf("room add %d %s", room.Id, json)
				roomMap[room.Id] = room
				q := make(chan interface{})
				roomQuit[room.Id] = &q // http://stackoverflow.com/questions/25601802/why-does-inline-instantiation-of-variable-requires-explicitly-taking-the-address
				go pollTfsRoomUsers(room, tfsApi, min, max/2, send, roomQuit[room.Id])
				go pollTfsRoomMessages(room, tfsApi, min, max/4, send, roomQuit[room.Id])
				delay = min
			}

			// removed
			for _, room := range roomMap {
				if _, ok := newRoomMap[room.Id]; ok {
					continue
				}

				json, err := json.Marshal(*room)
				if err != nil {
					trace.Log(err)
					continue
				}

				close(*roomQuit[room.Id])
				delete(roomMap, room.Id)
				delete(roomQuit, room.Id)
				*send <- fmt.Sprintf("room remove %d %s", room.Id, json)
				delay = min
			}

			timer.Reset(delay)
		case <-*recv:
			return
		}
	}
}

// routine to poll user information for a given room at variable intervals
func pollTfsRoomUsers(room *tfs.Room, tfsApi *tfs.Api, min time.Duration, max time.Duration, send *chan interface{}, recv *chan interface{}) {
	delay := min
	userMap := make(map[string]*tfs.RoomUser)

	timer := time.NewTimer(delay)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			delay = delay * 2
			if delay > max {
				delay = max
			}

			users, err := tfsApi.GetRoomUsers(room)
			if err != nil {
				trace.Log(err)
				continue
			}

			// added
			newUserMap := make(map[string]*tfs.RoomUser)
			for _, user := range users.Value {
				newUserMap[user.User.Id] = user
				if _, ok := userMap[user.User.Id]; ok {
					continue
				}

				json, err := json.Marshal(user)
				if err != nil {
					trace.Log(err)
					continue
				}

				*send <- fmt.Sprintf("user add %d %s %s", user.RoomId, user.User.Id, json)
				userMap[user.User.Id] = user
				delay = min
			}

			// removed
			for _, user := range userMap {
				if _, ok := newUserMap[user.User.Id]; ok {
					continue
				}

				json, err := json.Marshal(*user)
				if err != nil {
					trace.Log(err)
					continue
				}

				*send <- fmt.Sprintf("user remove %d %s %s", user.RoomId, user.User.Id, json)
				delete(userMap, user.User.Id)
				delay = min
			}

			// changed
			for id, user := range newUserMap {
				if newUserMap[id].IsOnline != userMap[id].IsOnline {
					json, err := json.Marshal(*user)
					if err != nil {
						trace.Log(err)
						continue
					}

					userMap[id] = newUserMap[id]
					*send <- fmt.Sprintf("user change %d %s %s", user.RoomId, user.User.Id, json)
					delay = min
				}
			}

			timer.Reset(delay)
		case <-*recv:
			return
		}
	}
}

// routine to poll message information for a given room at variable intervals
func pollTfsRoomMessages(room *tfs.Room, tfsApi *tfs.Api, min time.Duration, max time.Duration, send *chan interface{}, recv *chan interface{}) {
	delay := min
	messageMap := make(map[int]*tfs.RoomMessage)
	messageLast := room.LastActivity

	timer := time.NewTimer(delay)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			delay = delay * 2
			if delay > max {
				delay = max
			}

			messages, err := tfsApi.GetRoomMessages(room, messageLast)
			if err != nil {
				trace.Log(err)
				continue
			}

			// new
			for _, message := range messages.Value {
				if _, ok := messageMap[message.Id]; ok {
					continue
				}

				json, err := json.Marshal(message)
				if err != nil {
					trace.Log(err)
					continue
				}

				*send <- fmt.Sprintf("message new %d %s %d %s", message.PostedRoomId, message.PostedBy.Id, message.Id, json)
				messageMap[message.Id] = message
				messageLast = message.PostedTime
				delay = min
			}

			timer.Reset(delay)
		case <-*recv:
			return
		}
	}
}

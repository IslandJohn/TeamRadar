/*
Copyright 2016 IslandJohn and the TeamRadar Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IslandJohn/TeamRadar/Go/teamradar/tfs"
	"github.com/IslandJohn/TeamRadar/Go/teamradar/trace"
	"os"
	"strconv"
	"strings"
	"time"
)

// main runs the routines and collects
func main() {
	tfsApi, err := tfs.NewApi(os.Args[1], os.Args[2], os.Args[3])
	if err != nil {
		trace.Log(err)
		return
	}

	event, err := pullTfsAccount(tfsApi)
	if err != nil {
		trace.Log(err)
		return
	}
	fmt.Println(event)

	recv := make(chan interface{}) // routines send updates here
	quit := make(chan interface{}) // close this to have all routines return

	go pollTfsRooms(tfsApi, 1*time.Second, 60*time.Second, &recv, &quit)
	go pollCommandInterface(&recv, &quit)

	for event := range recv {
		if command := event.(string); strings.HasPrefix(command, "interface") {
			_, action, _, err := tokenizeEvent(command)
			if err != nil {
				trace.Log(err)
			}
			if action == "error" || action == "exit" || action == "logout" || action == "quit" {
				close(quit)
				return
			}
			
		} else {
			fmt.Println(event)
		}
	}
}

// return the routine, action, room, user, message of an event
func tokenizeEvent(event string) (string, string, int, error) {
	routine := ""
	action := ""
	room := 0
	err := error(nil)
	fields := strings.SplitN(event, " ", 3)

	if len(fields) >= 1 {
		routine = fields[0]
	}

	if len(fields) >= 2 {
		action = fields[1]
	}

	if len(fields) >= 3 {
		if action == "error" {
			err = errors.New(fields[2])
		} else {
			room, _ = strconv.Atoi(fields[2])
		}
	}

	return routine, action, room, err
}

// pulls the account information of the user
func pullTfsAccount(tfsApi *tfs.Api) (interface{}, error) {
	json, err := json.Marshal(tfsApi.LoginAccount)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("accounts login %s", json), nil
}

// routine to poll room information at variable intervals
func pollTfsRooms(tfsApi *tfs.Api, min time.Duration, max time.Duration, send *chan interface{}, recv *chan interface{}) {
	delay := min
	numErrors := 0
	roomRecv := make(chan interface{})          // routines started here we'll be proxied
	roomQuit := make(map[int]*chan interface{}) // we'll close this to have routines we started return
	roomMap := make(map[int]*tfs.Room)
	defer close(roomRecv)

	timer := time.NewTimer(delay)
	defer timer.Stop()
	for {
		select {
		case <-timer.C: // tick tock
			delay = delay * 2
			if delay > max {
				delay = max
			}

			rooms, err := tfsApi.GetRooms()
			if err != nil {
				trace.Log(err)
				numErrors++
				if numErrors >= 3 {
					for _, quit := range roomQuit { // clean up on error
						close(*quit) // routines should return on this being closed
					}
					*send <- fmt.Sprintf("rooms error %s", err)
					return
				}
				time.Sleep(time.Duration(numErrors) * time.Second)
				continue
			}
			numErrors = 0

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

				*send <- fmt.Sprintf("rooms add %d %s", room.Id, json)
				roomMap[room.Id] = room
				q := make(chan interface{})
				roomQuit[room.Id] = &q // http://stackoverflow.com/questions/25601802/why-does-inline-instantiation-of-variable-requires-explicitly-taking-the-address
				go pollTfsRoomUsers(room, tfsApi, min, max/2, &roomRecv, roomQuit[room.Id])
				go pollTfsRoomMessages(room, tfsApi, min, max/4, &roomRecv, roomQuit[room.Id])
				delay = min
			}

			// removed
			for _, room := range roomMap {
				if _, ok := newRoomMap[room.Id]; ok {
					continue
				}

				close(*roomQuit[room.Id])
				delete(roomMap, room.Id)
				delete(roomQuit, room.Id)
				*send <- fmt.Sprintf("rooms remove %d", room.Id)
				delay = min
			}

			timer.Reset(delay)
		case event := <-roomRecv: // relay send or clean up from a routine that errored
			_, action, room, err := tokenizeEvent(event.(string))
			if err != nil {
				trace.Log(err)
			}
			if action == "error" { // routine error
				quit, ok := roomQuit[room]
				if ok { // cleanup
					close(*quit)
					delete(roomMap, room)
					delete(roomQuit, room)
					*send <- fmt.Sprintf("rooms remove %d", room)
				}
			} else {
				*send <- event // relay
			}
		case <-*recv: // we need to quit
			for _, quit := range roomQuit { // clean up routines
				close(*quit) // routines should return on this being closed
			}
			return
		}
	}
}

// routine to poll user information for a given room at variable intervals
func pollTfsRoomUsers(room *tfs.Room, tfsApi *tfs.Api, min time.Duration, max time.Duration, send *chan interface{}, recv *chan interface{}) {
	delay := min
	numErrors := 0
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
				numErrors++
				if numErrors >= 3 {
					*send <- fmt.Sprintf("users error %d %s", room.Id, err)
					return
				}
				time.Sleep(time.Duration(numErrors) * time.Second)
				continue
			}
			numErrors = 0

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

				*send <- fmt.Sprintf("users add %d %s %s", user.RoomId, user.User.Id, json)
				userMap[user.User.Id] = user
				delay = min
			}

			// removed
			for _, user := range userMap {
				if _, ok := newUserMap[user.User.Id]; ok {
					continue
				}

				*send <- fmt.Sprintf("users remove %d %s", user.RoomId, user.User.Id)
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
					*send <- fmt.Sprintf("users change %d %s %s", user.RoomId, user.User.Id, json)
					delay = min
				}
			}

			timer.Reset(delay)
		case <-*recv: // quit
			return
		}
	}
}

// routine to poll message information for a given room at variable intervals
func pollTfsRoomMessages(room *tfs.Room, tfsApi *tfs.Api, min time.Duration, max time.Duration, send *chan interface{}, recv *chan interface{}) {
	delay := min
	numErrors := 0
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
				numErrors++
				if numErrors >= 3 {
					*send <- fmt.Sprintf("messages error %d %s", room.Id, err)
					return
				}
				time.Sleep(time.Duration(numErrors) * time.Second)
				continue
			}
			numErrors = 0

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

				*send <- fmt.Sprintf("messages new %d %s %d %s", message.PostedRoomId, message.PostedBy.Id, message.Id, json)
				messageMap[message.Id] = message
				messageLast = message.PostedTime
				delay = min
			}

			timer.Reset(delay)
		case <-*recv: // quit
			return
		}
	}
}

// read standard input for commands
func pollCommandInterface(send *chan interface{}, recv *chan interface{}) {
	numErrors := 0
	input := make(chan interface{})
	defer close(input)

	// this will read lines from stdin and send it over a channel
	go func() {
		stdin := bufio.NewReader(os.Stdin)

		for {
			line, _, err := stdin.ReadLine()

			if err != nil {
				numErrors++
				if numErrors >= 3 {
					input <- err
					return
				}
				time.Sleep(time.Duration(numErrors) * time.Second)
				continue
			}
			numErrors = 0

			input <- string(line)
		}
	}()

	for {
		select {
		case event := <-input:
			line, ok := event.(string)
			if ok {
				*send <- fmt.Sprintf("interface %s", line)
			} else {
				*send <- fmt.Sprintf("interface error %s", fmt.Sprint(event))
				return
			}
		case <-*recv: // we need to quit
			return
		}
	}
}

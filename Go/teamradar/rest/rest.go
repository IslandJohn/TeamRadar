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

package rest

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// REST client
type Client struct {
	transport http.RoundTripper
	tries     int
	user      string
	password  string
}

func NewClient(errs int) *Client {
	return &Client{
		transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
		tries: errs,
	}
}

// set the login information
func (c *Client) SetLogin(u string, p string) {
	c.user = u
	c.password = p
}

// make a request, expect a code, and return the body or error
func (c *Client) MakeRequest(verb string, url string, body string, code int) (map[string][]string, []byte, error) {
	numErrors := 0
	for {
		request, err := http.NewRequest(verb, url, strings.NewReader(body))
		if err != nil {
			return nil, nil, err
		}
		if c.user != "" {
			request.SetBasicAuth(c.user, c.password)
		}
		request.Header.Set("Content-Type", "application/json")
		
		response, err := c.transport.RoundTrip(request)
		if err != nil {
			numErrors++
			if numErrors >= c.tries {
				return nil, nil, err
			}
			time.Sleep(time.Duration(numErrors) * time.Second)
			continue
		}
		numErrors = 0

		defer response.Body.Close()
		if response.StatusCode != code {
			return nil, nil, errors.New(fmt.Sprintf("%d %s", response.StatusCode, response.Status))
		}

		ret, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, nil, err
		}

		return response.Header, ret, err
	}
}

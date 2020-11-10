/*
 * Copyright 2020 Unisys Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package helpers contains common functions used by the rest of the Gateway
package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

var log = GetLogger("HttpClient")

// Timeout for a HTTP Request to the agent.
// Currently unenforced
const timeout = 20000

// ErrNotFound is an error when an asset is not found
var ErrNotFound = errors.New("not found")

// ErrAlreadyExistsOnAgent is an error when a conflict occurs
var ErrAlreadyExistsOnAgent = errors.New("conflict, already exists")

// ErrUnauthorized is an error when the agent is unauthorized to perform this action
var ErrUnauthorized = errors.New("the entity that this agent is authenticated as is not authorized to perform this operation")

// PostJSONRequest executes a POST request to a url with a Text/JSON body and optionally additional headers
func PostJSONRequest(base string, path string, headerExtra map[string]string, bytesRepresentation []byte) (result map[string]interface{}, err error) {

	client := &http.Client{
		//Timeout:timeout,
	}
	req, err := http.NewRequest("POST", base+path, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Err(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	if headerExtra != nil {
		for key, value := range headerExtra {
			req.Header.Add(key, value)
		}
	}
	resp, err := client.Do(req)

	if err != nil {
		log.Err(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusConflict {
			err = ErrAlreadyExistsOnAgent
		} else if resp.StatusCode == http.StatusUnauthorized {
			err = ErrUnauthorized
		} else {
			err = errors.New("Got non OK statuscode: " + strconv.Itoa(resp.StatusCode))
		}
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// GetRequest executes a GET request to a URL, with optional headers
func GetRequest(url string, path string, headerExtra map[string]string) (result io.ReadCloser, err error) {

	req, err := http.NewRequest("GET", url+path, nil)

	if err != nil {
		log.Err(err)
		return
	}

	client := &http.Client{
		//Timeout:timeout,
	}
	req.Header.Add("Content-Type", "application/json")
	if headerExtra != nil {
		for key, value := range headerExtra {
			req.Header.Add(key, value)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		} else if resp.StatusCode == http.StatusUnauthorized {
			err = ErrUnauthorized
		} else {
			err = errors.New("Got non OK statuscode: " + strconv.Itoa(resp.StatusCode))
		}
		return
	}
	result = resp.Body

	return
}

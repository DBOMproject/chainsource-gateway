/*
 * Copyright 2023 Unisys Corporation
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

package helpers

import (
	"testing"

	"gopkg.in/h2non/gock.v1"
)

func setupMockRemoteEP(statusCode int) {
	gock.New("http://mock-addr").
		Post("/").
		Reply(statusCode).
		JSON(map[string]string{"foo": "bar"})
	gock.New("http://mock-addr").
		Get("/").
		Reply(statusCode).
		JSON(map[string]string{"foo": "bar"})
}

func setupMockRemoteEPInvalidJSON(statusCode int) {
	gock.New("http://mock-addr").
		Post("/").
		Reply(statusCode).
		BodyString("invalid.json")
	gock.New("http://mock-addr").
		Get("/").
		Reply(statusCode).
		BodyString("invalid.json")
}

func TestPostJSONRequest(t *testing.T) {
	// TODO: Test cases for PostJSONRequest function
}

func TestGetRequest(t *testing.T) {
	// TODO: Test cases for GetRequest function
}

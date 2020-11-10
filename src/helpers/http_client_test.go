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

package helpers

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"
)

// setupMockRemoteEP sets up a gock based remote handler that simulates a http endpoint
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

// setupMockRemoteEP sets up a gock based remote handler that simulates a http endpoint that returns bad JSON
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

// TestPostJSONRequest tests the POST JSON function
func TestPostJSONRequest(t *testing.T) {

	t.Run("When_Remote_Responds_OK_No_Headers", func(t *testing.T) {
		setupMockRemoteEP(http.StatusOK)
		defer gock.Off()

		_, err := PostJSONRequest("http://mock-addr", "/", nil,  []byte(""))
		assert.NoError(t, err, "Request is successful")
	})
	t.Run("When_Bad_URL", func(t *testing.T) {
		_, err := PostJSONRequest(":\\*.", "/", nil,  []byte(""))
		assert.Error(t, err, "URL Parse Error is handled")
	})
	t.Run("When_Lookup_Fail", func(t *testing.T) {
		_, err := PostJSONRequest("http://null", "/", nil,  []byte(""))
		assert.Error(t, err, "Domain lookup error is handled")
	})
	t.Run("When_Remote_Responds_OK_With_Headers", func(t *testing.T) {
		setupMockRemoteEP(http.StatusOK)
		defer gock.Off()
		headerMap := map[string]string{}
		headerMap["foo"] = "bar"

		_, err := PostJSONRequest("http://mock-addr", "/", headerMap,  []byte(""))
		assert.NoError(t, err, "Request is successful")
	})
	t.Run("When_Remote_Responds_Conflict", func(t *testing.T) {
		setupMockRemoteEP(http.StatusConflict)
		defer gock.Off()

		_, err := PostJSONRequest("http://mock-addr", "/", nil,  []byte(""))
		assert.EqualError(t, err, "conflict, already exists", "Conflict error is returned")
	})
	t.Run("When_Remote_Responds_Unauthorized", func(t *testing.T) {
		setupMockRemoteEP(http.StatusUnauthorized)
		defer gock.Off()

		_, err := PostJSONRequest("http://mock-addr", "/", nil,  []byte(""))
		assert.EqualError(t, err, "the entity that this agent is authenticated as is not authorized to perform this operation", "Unauthorized error is returned")
	})
	t.Run("When_Remote_Responds_InternalServerError", func(t *testing.T) {
		setupMockRemoteEP(http.StatusInternalServerError)
		defer gock.Off()

		_, err := PostJSONRequest("http://mock-addr", "/", nil,  []byte(""))
		assert.Error(t, err, "Error is handled")
	})
	t.Run("When_Remote_Responds_InvalidJSON", func(t *testing.T) {
		setupMockRemoteEPInvalidJSON(http.StatusOK)
		defer gock.Off()

		_, err := PostJSONRequest("http://mock-addr", "/", nil,  []byte(""))
		assert.Error(t, err, "JSON error is handled")
	})
}

// TestGetRequest tests the GET function
func TestGetRequest(t *testing.T) {

	t.Run("When_Remote_Responds_OK_No_Headers", func(t *testing.T) {
		setupMockRemoteEP(http.StatusOK)
		defer gock.Off()

		_, err := GetRequest("http://mock-addr", "/", nil)
		assert.NoError(t, err, "Request is successful")
	})
	t.Run("When_Bad_URL", func(t *testing.T) {
		_, err := GetRequest(":\\*.", "/", nil)
		assert.Error(t, err, "URL Parse Error is handled")
	})
	t.Run("When_Lookup_Fail", func(t *testing.T) {
		_, err := GetRequest("http://null", "/", nil)
		assert.Error(t, err, "Domain lookup error is handled")
	})
	t.Run("When_Remote_Responds_OK_With_Headers", func(t *testing.T) {
		setupMockRemoteEP(http.StatusOK)
		defer gock.Off()
		headerMap := map[string]string{}
		headerMap["foo"] = "bar"

		_, err := GetRequest("http://mock-addr", "/", headerMap)
		assert.NoError(t, err, "Request is successful")
	})
	t.Run("When_Remote_Responds_NotFound", func(t *testing.T) {
		setupMockRemoteEP(http.StatusNotFound)
		defer gock.Off()

		_, err := GetRequest("http://mock-addr", "/", nil)
		assert.EqualError(t, err, "not found", "Not Found error is returned")
	})
	t.Run("When_Remote_Responds_Unauthorized", func(t *testing.T) {
		setupMockRemoteEP(http.StatusUnauthorized)
		defer gock.Off()

		_, err := GetRequest("http://mock-addr", "/", nil)
		assert.EqualError(t, err, "the entity that this agent is authenticated as is not authorized to perform this operation", "Unauthorized error is returned")
	})
	t.Run("When_Remote_Responds_InternalServerError", func(t *testing.T) {
		setupMockRemoteEP(http.StatusInternalServerError)
		defer gock.Off()

		_, err := GetRequest("http://mock-addr", "/", nil)
		assert.Error(t, err, "Error is handled")
	})
}

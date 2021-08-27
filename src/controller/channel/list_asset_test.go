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

package channel

import (
	"chainsource-gateway/helpers"
	"chainsource-gateway/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var assetsPath = "../../testdata/asset_controller_tests/listAssets/listAssets.json"

// TestListAssets is the happy path for listing assets
func TestListAssets(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	var expected interface{}
	json.NewDecoder(openTestJSON(assetsPath)).Decode(&expected)

	mockAgent.EXPECT().ListAssets(gomock.Any(),
		mocks.AgentQueryFor("C1", "")).
		Return(openTestJSON(assetsPath), nil)
	mockRequest := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "", mockAgent,
		mocks.NewMockAssetSchemaAlwaysValid(ctrl))
	handler := http.HandlerFunc(ListAssets)
	handler.ServeHTTP(responseRecorder, mockRequest)

	var actual interface{}
	json.NewDecoder(responseRecorder.Result().Body).Decode(&actual)
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode, "Response Should be 200 OK")
	assert.Equal(t, expected, actual,
		"Stored asset is exactly equal to the sent asset when all fields are filled")
}

// TestListAssetsUInauthorizedErrorConditions checks unauthorized handling
func TestListAssetsUInauthorizedErrorConditions(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	assetC1A1 := openTestJSON(assetsPath)
	var expected map[string]interface{}
	json.NewDecoder(assetC1A1).Decode(&expected)

	var empty interface{}
	mockAgent.EXPECT().ListAssets(gomock.Any(),
		mocks.AgentQueryFor("C1", "")).
		Return(nil, helpers.ErrUnauthorized)
	mockRequest := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "", mockAgent,
		mocks.NewMockAssetSchemaAlwaysValid(ctrl))
	mockRequest = injectMockQueryAssetContext(mockRequest, empty, empty, 0, 10)
	handler := http.HandlerFunc(ListAssets)
	handler.ServeHTTP(responseRecorder, mockRequest)

	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 not authorized")
}

// TestListAssetsAgentErrorConditions checks bad gateway handling
func TestListAssetsAgentErrorConditions(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	assetC1A1 := openTestJSON(assetsPath)
	var expected map[string]interface{}
	json.NewDecoder(assetC1A1).Decode(&expected)

	var empty interface{}
	mockAgent.EXPECT().ListAssets(gomock.Any(),
		mocks.AgentQueryFor("C1", "")).
		Return(nil, errors.New("ANY"))
	mockRequest := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "", mockAgent,
		mocks.NewMockAssetSchemaAlwaysValid(ctrl))
	mockRequest = injectMockQueryAssetContext(mockRequest, empty, empty, 0, 10)
	handler := http.HandlerFunc(ListAssets)
	handler.ServeHTTP(responseRecorder, mockRequest)

	assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 bad gateway")
}

// TestListAssetsErrorConditions checks not found handling
func TestListAssetsErrorConditions(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	var empty interface{}
	mockAgent.EXPECT().ListAssets(gomock.Any(),
		mocks.AgentQueryFor("C1", "")).
		Return(nil, helpers.ErrNotFound)
	mockRequest := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "", mockAgent,
		mocks.NewMockAssetSchemaAlwaysValid(ctrl))
	mockRequest = injectMockQueryAssetContext(mockRequest, empty, empty, 0, 10)
	handler := http.HandlerFunc(ListAssets)
	handler.ServeHTTP(responseRecorder, mockRequest)

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode, "Response Should be 404 not found")
}

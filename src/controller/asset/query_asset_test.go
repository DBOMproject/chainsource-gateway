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

package asset

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

// TestQueryAllAssets is the happy path for querying all assets
func TestQueryAllAssets(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	var queryPath = "../../testdata/asset_controller_tests/query/queryAsset.json"

	assetC1A1 := openTestJSON(queryPath)
	var expected map[string]interface{}
	json.NewDecoder(assetC1A1).Decode(&expected)

	var empty interface{}
	var emptyArray []string
	mockAgent.EXPECT().QueryAssets(gomock.Any(),
		mocks.AgentQueryFor("C1", ""), mocks.AgentRichQueryFor(nil, emptyArray, 0, 10)).
		Return(expected, nil)
	mockRequest := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "", mockAgent,
		mocks.NewMockAssetSchemaAlwaysValid(ctrl))
	mockRequest = injectMockQueryAssetContext(mockRequest, empty, empty, 0, 10)
	handler := http.HandlerFunc(QueryAsset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	var actual map[string]interface{}
	json.NewDecoder(responseRecorder.Result().Body).Decode(&actual)
	assert.Equal(t, expected, actual,
		"Stored asset is exactly equal to the sent asset when all fields are filled")
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode, "Response Should be 200 OK")
}

// TestQueryInvalidQueryAssets checks invalid query handling
func TestQueryInvalidQueryAssets(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	var empty interface{}
	mockRequest := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "", mockAgent,
		mocks.NewMockAssetSchemaAlwaysInvalid(ctrl))
	mockRequest = injectMockQueryAssetContext(mockRequest, "test", empty, 0, 10)
	handler := http.HandlerFunc(QueryAsset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode, "Response Should be 400 BAD REQUEST")
}

// TestQueryErrorQueryAssets checks query error handling
func TestQueryErrorQueryAssets(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	var empty interface{}
	mockRequest := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "", mockAgent,
		mocks.NewMockAssetSchemaAlwaysFailure(ctrl))
	mockRequest = injectMockQueryAssetContext(mockRequest, "test", empty, 0, 10)
	handler := http.HandlerFunc(QueryAsset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode, "Response Should be 500 INTERNAL SERVER ERROR")
}

// TestQueryAllAssetsErrorConditions checks not found handling
func TestQueryAllAssetsErrorConditions(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	var empty interface{}
	var emptyArray []string
	mockAgent.EXPECT().QueryAssets(gomock.Any(),
		mocks.AgentQueryFor("C1", ""), mocks.AgentRichQueryFor(nil, emptyArray, 0, 10)).
		Return(nil, helpers.ErrNotFound)
	mockRequest := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "", mockAgent,
		mocks.NewMockAssetSchemaAlwaysValid(ctrl))
	mockRequest = injectMockQueryAssetContext(mockRequest, empty, empty, 0, 10)
	handler := http.HandlerFunc(QueryAsset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode, "Response Should be 404 not found")
}

// TestQueryAllAssetsUnauthorizedErrorConditions checks unauthorized handling
func TestQueryAllAssetsUnauthorizedErrorConditions(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	var queryPath = "../../testdata/asset_controller_tests/query/queryAsset.json"

	assetC1A1 := openTestJSON(queryPath)
	var expected map[string]interface{}
	json.NewDecoder(assetC1A1).Decode(&expected)

	var empty interface{}
	var emptyArray []string
	mockAgent.EXPECT().QueryAssets(gomock.Any(),
		mocks.AgentQueryFor("C1", ""), mocks.AgentRichQueryFor(nil, emptyArray, 0, 10)).
		Return(nil, helpers.ErrUnauthorized)
	mockRequest := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "", mockAgent,
		mocks.NewMockAssetSchemaAlwaysValid(ctrl))
	mockRequest = injectMockQueryAssetContext(mockRequest, empty, empty, 0, 10)
	handler := http.HandlerFunc(QueryAsset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 not authorized")
}

// TestQueryAllAssetsAgentErrorConditions checks bad gateway handling
func TestQueryAllAssetsAgentErrorConditions(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	var queryPath = "../../testdata/asset_controller_tests/query/queryAsset.json"

	assetC1A1 := openTestJSON(queryPath)
	var expected map[string]interface{}
	json.NewDecoder(assetC1A1).Decode(&expected)

	var empty interface{}
	var emptyArray []string
	mockAgent.EXPECT().QueryAssets(gomock.Any(),
		mocks.AgentQueryFor("C1", ""), mocks.AgentRichQueryFor(nil, emptyArray, 0, 10)).
		Return(nil, errors.New("ANY"))
	mockRequest := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "", mockAgent,
		mocks.NewMockAssetSchemaAlwaysValid(ctrl))
	mockRequest = injectMockQueryAssetContext(mockRequest, empty, empty, 0, 10)
	handler := http.HandlerFunc(QueryAsset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 bad gateway")
}

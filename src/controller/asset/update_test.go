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
	"chainsource-gateway/agent"
	"chainsource-gateway/helpers"
	"chainsource-gateway/mocks"
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var assetToBeUpdatedPath = "../../testdata/asset_controller_tests/update/assetToBeUpdated.json"
var assetToBeUpdatedReadOnlyPath = "../../testdata/asset_controller_tests/update/assetToBeUpdatedReadOnly.json"
var assetUpdatePayload = "../../testdata/asset_controller_tests/update/assetUpdatePayload.json"

// TestUpdateWithAssetIDThatExists is the happy path for update
func TestUpdateWithAssetIDThatExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()
	var finalAsset helpers.Asset

	assetC1A1 := openTestJSON(assetToBeUpdatedPath)
	mockAgent.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C1", "A1")).
		Return(assetC1A1, nil)
	mockAgent.EXPECT().Commit(gomock.Any(),
		mocks.AgentCommitTo("C1", "A1", "UPDATE")).
		Do(func(ctx context.Context, args agent.CommitArgs) {
			finalAsset = args.Payload
		}).
		Return(getAgentSuccessResponse(), nil)

	updatePayload := openTestJSON(assetUpdatePayload)
	mockRequest := httptest.NewRequest("POST", "/", updatePayload)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
		mocks.NewMockAssetSchemaAlwaysValid(ctrl))
	handler := http.HandlerFunc(UpdateAsset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode, "Response Should be 200 OK")
	assert.Equal(t, 1, len(finalAsset.AttachedChildren), "Child asset links are not changed by update")
	assert.Equal(t, "DB1", finalAsset.ParentAsset.RepoID, "Parent asset link update is ignored")
}

// TestUpdateWithAssetErrorConditions contains the tests that simulate error conditions of the asset
func TestUpdateWithAssetErrorConditions(t *testing.T){
	t.Run("Invalid_AssetID",func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()
		assetC1A1 := openTestJSON(assetToBeUpdatedPath)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(nil, helpers.ErrNotFound)

		mockRequest := httptest.NewRequest("POST", "/", assetC1A1)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		handler := http.HandlerFunc(UpdateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	} )
	t.Run("Read_Only_Asset",func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		assetC1A1 := openTestJSON(assetToBeUpdatedReadOnlyPath)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)

		mockRequest := httptest.NewRequest("PUT", "/", openTestJSON(assetUpdatePayload))
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		handler := http.HandlerFunc(UpdateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		var result map[string]interface{}
		json.NewDecoder(responseRecorder.Result().Body).Decode(&result)

		assert.Equal(t, http.StatusConflict, responseRecorder.Result().StatusCode, "Response Should be 409 CONFLICT")
	} )
}

// TestUpdateWithRequestErrorConditions contains the tests that simulate user request error conditions
func TestUpdateWithRequestErrorConditions(t *testing.T){
	t.Run("Invalid_Payload",func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()
		assetC1A1 := openTestJSON(assetToBeUpdatedPath)

		mockRequest := httptest.NewRequest("POST", "/", assetC1A1)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysInvalid(ctrl))
		handler := http.HandlerFunc(UpdateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		var result map[string]interface{}
		json.NewDecoder(responseRecorder.Result().Body).Decode(&result)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode, "Response Should be 400 BAD REQUEST")
		assert.Contains(t, result, "error", "Response Contains Error String")
	} )
	t.Run("Validator_Failure",func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()
		assetC1A1 := openTestJSON(assetToBeUpdatedPath)

		mockRequest := httptest.NewRequest("POST", "/", assetC1A1)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysFailure(ctrl))
		handler := http.HandlerFunc(UpdateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		var result map[string]interface{}
		json.NewDecoder(responseRecorder.Result().Body).Decode(&result)

		assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode, "Response Should be 500 INTERNAL SERVER")
	} )
	t.Run("Non_JSON",func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		mockRequest := httptest.NewRequest("POST", "/", strings.NewReader(invalidJSON))
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		handler := http.HandlerFunc(UpdateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		var result map[string]interface{}
		json.NewDecoder(responseRecorder.Result().Body).Decode(&result)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode, "Response Should be 400 INTERNAL SERVER")
	} )
}

// TestUpdateWithAgentFailureConditions contains the test for agent failures
func TestUpdateWithAgentFailureConditions(t *testing.T) {
	t.Run("Query_Failure",func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		assetC1A1 := openTestJSON(assetToBeUpdatedPath)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(nil, errors.New(""))

		mockRequest := httptest.NewRequest("POST", "/", assetC1A1)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		handler := http.HandlerFunc(UpdateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		var result map[string]interface{}
		json.NewDecoder(responseRecorder.Result().Body).Decode(&result)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
	t.Run("Query_Unauthorized",func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		assetC1A1 := openTestJSON(assetToBeUpdatedPath)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(nil, helpers.ErrUnauthorized)

		mockRequest := httptest.NewRequest("POST", "/", assetC1A1)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		handler := http.HandlerFunc(UpdateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		var result map[string]interface{}
		json.NewDecoder(responseRecorder.Result().Body).Decode(&result)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 UNAUTHORIZED")
	})
	t.Run("Commit_Failure",func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		assetC1A1 := openTestJSON(assetToBeUpdatedPath)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C1", "A1", "UPDATE")).
			Return(nil, errors.New(""))

		mockRequest := httptest.NewRequest("POST", "/", assetC1A1)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		handler := http.HandlerFunc(UpdateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		var result map[string]interface{}
		json.NewDecoder(responseRecorder.Result().Body).Decode(&result)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	} )
	t.Run("Commit_Unauthorized",func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		assetC1A1 := openTestJSON(assetToBeUpdatedPath)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C1", "A1", "UPDATE")).
			Return(nil, helpers.ErrUnauthorized)

		mockRequest := httptest.NewRequest("POST", "/", assetC1A1)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		handler := http.HandlerFunc(UpdateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		var result map[string]interface{}
		json.NewDecoder(responseRecorder.Result().Body).Decode(&result)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 UNAUTHORIZED")
	} )
}

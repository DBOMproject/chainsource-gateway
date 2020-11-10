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

var assetToBeCreatedPath = "../../testdata/asset_controller_tests/create/assetToBeCreated.json"

// TestCreateWithAssetIDThatDoesNotExist contains the happy path for creation of an asset
func TestCreateWithAssetIDThatDoesNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()
	var finalAsset helpers.Asset

	assetC1A1 := openTestJSON(assetToBeCreatedPath)
	mockAgent.EXPECT().Commit(gomock.Any(),
		mocks.AgentCommitTo("C1", "A1", "CREATE")).
		Do(func(ctx context.Context, args agent.CommitArgs) {
			finalAsset = args.Payload
		}).
		Return(getAgentSuccessResponse(), nil)

	mockRequest := httptest.NewRequest("POST", "/", assetC1A1)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
		mocks.NewMockAssetSchemaAlwaysValid(ctrl))
	handler := http.HandlerFunc(CreateAsset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	assert.Equal(t, http.StatusCreated, responseRecorder.Result().StatusCode, "Response Should be 201 CREATED")
	assert.Equal(t, 0, len(finalAsset.AttachedChildren), "All child asset links are stripped")
	assert.Equal(t, "", finalAsset.ParentAsset.RepoID, "Parent asset link is stripped")
}

// TestCreateWithAssetErrorConditions contains the tests that simulate error conditions of the asset
func TestCreateWithAssetErrorConditions(t *testing.T) {
	t.Run("Existing_AssetID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()
		assetC1A1 := openTestJSON(assetToBeCreatedPath)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C1", "A1", "CREATE")).
			Return(nil, helpers.ErrAlreadyExistsOnAgent)

		mockRequest := httptest.NewRequest("POST", "/", assetC1A1)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		handler := http.HandlerFunc(CreateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusConflict, responseRecorder.Result().StatusCode, "Response Should be 409 CONFLICT")
	})
}

// TestCreateWithRequestErrorConditions contains the tests that simulate user request error conditions
func TestCreateWithRequestErrorConditions(t *testing.T) {
	t.Run("Invalid_Payload", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()
		assetC1A1 := openTestJSON(assetToBeCreatedPath)

		mockRequest := httptest.NewRequest("POST", "/", assetC1A1)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysInvalid(ctrl))
		handler := http.HandlerFunc(CreateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		var result map[string]interface{}
		json.NewDecoder(responseRecorder.Result().Body).Decode(&result)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode, "Response Should be 400 BAD REQUEST")
		assert.Contains(t, result, "error", "Response Contains Error String")
	})
	t.Run("Validator_Failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		assetC1A1 := openTestJSON(assetToBeCreatedPath)

		mockRequest := httptest.NewRequest("POST", "/", assetC1A1)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysFailure(ctrl))
		handler := http.HandlerFunc(CreateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		var result map[string]interface{}
		json.NewDecoder(responseRecorder.Result().Body).Decode(&result)

		assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode, "Response Should be 500 INTERNAL SERVER")
	})
	t.Run("Non_JSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		mockRequest := httptest.NewRequest("POST", "/", strings.NewReader(invalidJSON))
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		handler := http.HandlerFunc(CreateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		var result map[string]interface{}
		json.NewDecoder(responseRecorder.Result().Body).Decode(&result)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode, "Response Should be 400 INTERNAL SERVER")
	})
}

// TestCreateWithAgentErrorConditions contains the test for agent failures
func TestCreateWithAgentErrorConditions(t *testing.T) {
	t.Run("Commit_Failure", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		assetC1A1 := openTestJSON(assetToBeCreatedPath)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C1", "A1", "CREATE")).
			Return(nil, errors.New(""))

		mockRequest := httptest.NewRequest("POST", "/", assetC1A1)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		handler := http.HandlerFunc(CreateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		var result map[string]interface{}
		json.NewDecoder(responseRecorder.Result().Body).Decode(&result)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
	t.Run("Commit_Unauthorized", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		assetC1A1 := openTestJSON(assetToBeCreatedPath)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C1", "A1", "CREATE")).
			Return(nil, helpers.ErrUnauthorized)

		mockRequest := httptest.NewRequest("POST", "/", assetC1A1)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		handler := http.HandlerFunc(CreateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		var result map[string]interface{}
		json.NewDecoder(responseRecorder.Result().Body).Decode(&result)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 UNAUTHORIZED")
	})
}

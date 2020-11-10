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
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var assetC1A1Path = "../../testdata/asset_controller_tests/retrieve/retrievableAsset.json"

// TestRetrieveWithAssetThatExists is the happy path for retrieve
func TestRetrieveWithAssetThatExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	assetC1A1 := openTestJSON(assetC1A1Path)
	mockAgent.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C1", "A1")).
		Return(assetC1A1, nil)
	mockRequest := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
		mocks.NewMockAssetSchemaAlwaysValid(ctrl))
	handler := http.HandlerFunc(RetrieveAsset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	assert.Equal(t, decodeAsset(openTestJSON(assetC1A1Path)), decodeAsset(responseRecorder.Result().Body),
		"Stored asset is exactly equal to the sent asset when all fields are filled")
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode, "Response Should be 200 OK")
}

// TestRetrieveWithAssetErrorConditions contains the tests that simulate error conditions of the asset
func TestRetrieveWithAssetErrorConditions(t *testing.T) {
	t.Run("Invalid_AssetID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			gomock.Any()).
			Return(nil, helpers.ErrNotFound)
		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		handler := http.HandlerFunc(RetrieveAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode, "Response Should be 404 NOT FOUND")
	})
}

// TestRetrieveWithAgentErrorConditions contains the test for agent failures
func TestRetrieveWithAgentErrorConditions(t *testing.T) {
	t.Run("Query_Failure", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			gomock.Any()).
			Return(nil, errors.New(""))
		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		handler := http.HandlerFunc(RetrieveAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
	t.Run("Query_Unauthorized", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			gomock.Any()).
			Return(nil, helpers.ErrUnauthorized)
		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		handler := http.HandlerFunc(RetrieveAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 UNAUTHORIZED")
	})
}

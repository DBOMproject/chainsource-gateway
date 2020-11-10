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
	"chainsource-gateway/pgp"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const signedAssetLocation = "../../testdata/asset_controller_tests/validate/assetToBeValidated.json"
const signedAssetNoFPLocation = "../../testdata/asset_controller_tests/validate/assetToBeValidatedNoFP.json"
const signedAssetEmptyFPLocation = "../../testdata/asset_controller_tests/validate/assetToBeValidatedEmptyFP.json"
const signedAssetEmptySLocation = "../../testdata/asset_controller_tests/validate/assetToBeValidatedEmptySignature.json"

// injectValidateContext adds a mock validator into a request
func injectValidateContext(r *http.Request, validator pgp.SignatureValidator) (reqWithContext *http.Request) {
	ctx := context.WithValue(r.Context(), "signatureValidator", validator)
	reqWithContext = r.WithContext(ctx)
	return
}

// signingServiceSuccessReturn returns expected success body from the signing service
func signingServiceSuccessReturn() map[string]interface{} {
	result := make(map[string]interface{})
	result["success"] = true
	result["valid"] = true
	return result
}

// TestValidateWithSignedAssetThatExists is the happy path for validate
func TestValidateWithSignedAssetThatExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	mockValidator := mocks.NewMockSignatureValidator(ctrl)
	mockValidator.EXPECT().Validate(gomock.Any(), gomock.Any()).
		Return(signingServiceSuccessReturn(), nil)
	defer ctrl.Finish()

	assetC1A1 := openTestJSON(signedAssetLocation)
	mockAgent.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C1", "A1")).
		Return(assetC1A1, nil)
	mockRequest := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
		mocks.NewMockAssetSchemaAlwaysValid(ctrl))
	mockRequest = injectValidateContext(mockRequest, mockValidator)
	handler := http.HandlerFunc(ValidateAsset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode, "Response Should be 200 OK")
}

// TestValidateWithAssetErrorConditions contains the tests that simulate error conditions of the asset
func TestValidateWithAssetErrorConditions(t *testing.T) {
	t.Run("Asset_DoesNotExist", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		mockValidator := mocks.NewMockSignatureValidator(ctrl)
		defer ctrl.Finish()

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(nil, helpers.ErrNotFound)
		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = injectValidateContext(mockRequest, mockValidator)
		handler := http.HandlerFunc(ValidateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode, "Response Should be 404 NOT FOUND")
	})
	t.Run("Asset_Unauthorized", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		mockValidator := mocks.NewMockSignatureValidator(ctrl)
		defer ctrl.Finish()

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(nil, helpers.ErrUnauthorized)
		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = injectValidateContext(mockRequest, mockValidator)
		handler := http.HandlerFunc(ValidateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 UNAUTHORIZED")
	})
	t.Run("Asset_DoesNotHaveManufactureFingerprint", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		mockValidator := mocks.NewMockSignatureValidator(ctrl)
		defer ctrl.Finish()

		assetC1A1 := openTestJSON(signedAssetNoFPLocation)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1,nil)
		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = injectValidateContext(mockRequest, mockValidator)
		handler := http.HandlerFunc(ValidateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode, "Response Should be 404 NOT FOUND")
	})
	t.Run("Asset_HasEmptyManufactureFingerprint", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		mockValidator := mocks.NewMockSignatureValidator(ctrl)
		defer ctrl.Finish()

		assetC1A1 := openTestJSON(signedAssetEmptyFPLocation)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1,nil)
		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = injectValidateContext(mockRequest, mockValidator)
		handler := http.HandlerFunc(ValidateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode, "Response Should be 404 NOT FOUND")
	})
	t.Run("Asset_HasEmptySignature", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		mockValidator := mocks.NewMockSignatureValidator(ctrl)
		defer ctrl.Finish()

		assetC1A1 := openTestJSON(signedAssetEmptySLocation)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1,nil)
		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = injectValidateContext(mockRequest, mockValidator)
		handler := http.HandlerFunc(ValidateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode, "Response Should be 404 NOT FOUND")
	})
}

// TestValidateWithAgentFailureConditions contains the test for agent failures
func TestValidateWithAgentFailureConditions(t *testing.T) {
	t.Run("Failure", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		mockValidator := mocks.NewMockSignatureValidator(ctrl)
		defer ctrl.Finish()

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(nil, errors.New(""))
		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = injectValidateContext(mockRequest, mockValidator)
		handler := http.HandlerFunc(ValidateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
}

// TestValidateWithSigningFailureConditions contains the test for signing service failures
func TestValidateWithSigningFailureConditions(t *testing.T) {
	t.Run("Failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		mockValidator := mocks.NewMockSignatureValidator(ctrl)
		mockValidator.EXPECT().Validate(gomock.Any(), gomock.Any()).
			Return(nil, errors.New(""))
		defer ctrl.Finish()

		assetC1A1 := openTestJSON(signedAssetLocation)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = injectValidateContext(mockRequest, mockValidator)
		handler := http.HandlerFunc(ValidateAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode, "Response Should be 400 BAD REQUEST")
	})
}

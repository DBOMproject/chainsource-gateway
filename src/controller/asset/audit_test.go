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
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var auditPath = "../../testdata/asset_controller_tests/audit/audit.json"

// TestAuditWithAssetThatExists contains the happy path test for audit
func TestAuditWithAssetThatExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	audit := openTestJSON(auditPath)
	var auditFinal  map[string]interface{}
	json.NewDecoder(audit).Decode(&auditFinal)

	mockAgent.EXPECT().QueryAuditTrail(gomock.Any(),
		mocks.AgentQueryFor("C1", "A1")).
		Return(auditFinal, nil)
	mockRequest := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
		mocks.NewMockAssetSchemaAlwaysValid(ctrl))
	handler := http.HandlerFunc(AuditAsset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	var auditResponse  map[string]interface{}
	json.NewDecoder(responseRecorder.Result().Body).Decode(&auditResponse)

	assert.Equal(t,auditFinal, auditResponse,"Retrieved audit trail is received with no changes")
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode, "Response Should be 200 OK")
}

// TestAuditWithAssetErrorConditions contains the tests that simulate error conditions of the asset
func TestAuditWithAssetErrorConditions(t *testing.T) {
	t.Run("Does_Not_Exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		mockAgent.EXPECT().QueryAuditTrail(gomock.Any(),
			gomock.Any()).
			Return(nil, helpers.ErrNotFound)
		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		handler := http.HandlerFunc(AuditAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
}

// TestAuditWithAgentErrorConditions contains the test for agent failures
func TestAuditWithAgentErrorConditions(t *testing.T) {
	t.Run("Audit_Query_Failure", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		mockAgent.EXPECT().QueryAuditTrail(gomock.Any(),
			gomock.Any()).
			Return(nil, errors.New("ANY"))
		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		handler := http.HandlerFunc(AuditAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
	t.Run("Audit_Query_Unauthorized", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		mockAgent.EXPECT().QueryAuditTrail(gomock.Any(),
			gomock.Any()).
			Return(nil, helpers.ErrUnauthorized)
		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		handler := http.HandlerFunc(AuditAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 UNAUTHORIZED")
	})
}

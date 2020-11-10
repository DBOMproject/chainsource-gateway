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
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const assetL0Location = "../../testdata/asset_controller_tests/export/levelZeroAsset.json"
const assetL0LoopLocation = "../../testdata/asset_controller_tests/export/levelZeroAssetWithPLoop.json"
const assetL1Location = "../../testdata/asset_controller_tests/export/levelOneAsset.json"
const assetL2A1Location = "../../testdata/asset_controller_tests/export/levelTwoAssetOne.json"
const assetL2A1LoopLocation = "../../testdata/asset_controller_tests/export/levelTwoAssetOneWithCLoop.json"
const assetL2A2Location = "../../testdata/asset_controller_tests/export/levelTwoAssetTwo.json"
const resultLocation = "../../testdata/asset_controller_tests/export/expectedResult.json"
const PLoopResultLocation = "../../testdata/asset_controller_tests/export/expectedPLoopResult.json"

// injectExportContext adds export routing variables to the request context
func injectExportContext(r *http.Request, vars helpers.ExportRoutingVars) (reqWithContext *http.Request) {
	ctx := context.WithValue(r.Context(), "exportVars", vars)
	reqWithContext = r.WithContext(ctx)
	return
}

// injectNullExportContext adds null export routing variables to the request context
func injectNullExportContext(r *http.Request) *http.Request {
	return injectExportContext(r, helpers.ExportRoutingVars{
		FileName:       "",
		InlineResponse: "",
	})
}

// TestExportOnHappyPath is the happy path test for export
func TestExportOnHappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	assetC0A0 := openTestJSON(assetL0Location)
	assetC1A1 := openTestJSON(assetL1Location)
	assetC2A2 := openTestJSON(assetL2A1Location)
	assetC3A3 := openTestJSON(assetL2A2Location)

	mockAgent.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C0", "A0")).
		Return(assetC0A0, nil)
	mockAgent.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C1", "A1")).
		Return(assetC1A1, nil)
	mockAgent.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C2", "A2")).
		Return(assetC2A2, nil)
	mockAgent.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C3", "A3")).
		Return(assetC3A3, nil)

	providerMap := make(map[string]agent.Agent)
	providerMap["T1"] = mockAgent

	mockRequest := httptest.NewRequest("GET", "/", nil)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
		mocks.NewMockAssetSchemaAlwaysValid(ctrl))
	mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
	mockRequest = injectNullExportContext(mockRequest)
	handler := http.HandlerFunc(ExportAsset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	buf := new(strings.Builder)
	io.Copy(buf, openTestJSON(resultLocation))

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode, "Response Should be 200 OK")
	assert.JSONEq(t, buf.String(), responseRecorder.Body.String(), "Export must match expected result")
	assert.Contains(t, responseRecorder.Header().Get("Content-Disposition"), "A1", "Export must contain assetID in content disposition")
}

// TestExportWithExportArgVariants contains happy path tests for various combinations of the export vars
func TestExportWithExportArgVariants(t *testing.T) {
	t.Run("With_InlineResponse", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		assetC0A0 := openTestJSON(assetL0Location)
		assetC1A1 := openTestJSON(assetL1Location)
		assetC2A2 := openTestJSON(assetL2A1Location)
		assetC3A3 := openTestJSON(assetL2A2Location)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C0", "A0")).
			Return(assetC0A0, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(assetC2A2, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C3", "A3")).
			Return(assetC3A3, nil)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		mockRequest = injectExportContext(mockRequest, helpers.ExportRoutingVars{
			FileName:       "",
			InlineResponse: "true",
		})
		handler := http.HandlerFunc(ExportAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		buf := new(strings.Builder)
		io.Copy(buf, openTestJSON(resultLocation))

		assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode, "Response Should be 200 OK")
		assert.JSONEq(t, buf.String(), responseRecorder.Body.String(), "Export must match expected result")
		assert.Equal(t, "" , responseRecorder.Header().Get("Content-Disposition"), "Export must not have a Content-Disposition Header")
	})
	t.Run("With_CustomFilename", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		assetC0A0 := openTestJSON(assetL0Location)
		assetC1A1 := openTestJSON(assetL1Location)
		assetC2A2 := openTestJSON(assetL2A1Location)
		assetC3A3 := openTestJSON(assetL2A2Location)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C0", "A0")).
			Return(assetC0A0, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(assetC2A2, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C3", "A3")).
			Return(assetC3A3, nil)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		mockRequest = injectExportContext(mockRequest, helpers.ExportRoutingVars{
			FileName:       "TestFilename",
			InlineResponse: "",
		})
		handler := http.HandlerFunc(ExportAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		buf := new(strings.Builder)
		io.Copy(buf, openTestJSON(resultLocation))

		assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode, "Response Should be 200 OK")
		assert.JSONEq(t, buf.String(), responseRecorder.Body.String(), "Export must match expected result")
		assert.Contains(t, responseRecorder.Header().Get("Content-Disposition"), "TestFilename", "Export must contain custom filename in content disposition")
	})
}

// TestExportWithAssetErrorConditions contains the tests that simulate error conditions of the asset
func TestExportWithAssetErrorConditions(t *testing.T) {
	t.Run("Asset_Loop_Parent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()


		assetC0A0 := openTestJSON(assetL0LoopLocation)
		assetC1A1 := openTestJSON(assetL1Location)
		assetC2A2 := openTestJSON(assetL2A1Location)
		assetC3A3 := openTestJSON(assetL2A2Location)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C0", "A0")).
			Return(assetC0A0, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(assetC2A2, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C3", "A3")).
			Return(assetC3A3, nil)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		mockRequest = injectNullExportContext(mockRequest)
		handler := http.HandlerFunc(ExportAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
	t.Run("Asset_Loop_Child", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		assetC0A0 := openTestJSON(assetL0Location)
		assetC1A1 := openTestJSON(assetL1Location)
		assetC2A2 := openTestJSON(assetL2A1LoopLocation)
		assetC3A3 := openTestJSON(assetL2A2Location)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C0", "A0")).
			Return(assetC0A0, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(assetC2A2, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C3", "A3")).
			Return(assetC3A3, nil)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		mockRequest = injectNullExportContext(mockRequest)
		handler := http.HandlerFunc(ExportAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")

	})
	t.Run("Asset_NotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(nil, helpers.ErrNotFound)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		mockRequest = injectNullExportContext(mockRequest)
		handler := http.HandlerFunc(ExportAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusNotFound, responseRecorder.Result().StatusCode, "Response Should be 404 NOT FOUND")

	})

	t.Run("Asset_Bad_JSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(ioutil.NopCloser(strings.NewReader(invalidJSON)), nil)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		mockRequest = injectNullExportContext(mockRequest)
		handler := http.HandlerFunc(ExportAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode, "Response Should be 500 INTERNAL SERVER ERROR")

	})
	t.Run("ParentAsset_NotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		assetC1A1 := openTestJSON(assetL1Location)
		assetC2A2 := openTestJSON(assetL2A1Location)
		assetC3A3 := openTestJSON(assetL2A2Location)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C0", "A0")).
			Return(nil, helpers.ErrNotFound)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(assetC2A2, nil).MaxTimes(1)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C3", "A3")).
			Return(assetC3A3, nil).MaxTimes(1)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		mockRequest = injectNullExportContext(mockRequest)
		handler := http.HandlerFunc(ExportAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")

	})
	t.Run("ParentAsset_Bad_JSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		assetC1A1 := openTestJSON(assetL1Location)
		assetC2A2 := openTestJSON(assetL2A1Location)
		assetC3A3 := openTestJSON(assetL2A2Location)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C0", "A0")).
			Return(ioutil.NopCloser(strings.NewReader(invalidJSON)), nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(assetC2A2, nil).MaxTimes(1)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C3", "A3")).
			Return(assetC3A3, nil).MaxTimes(1)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		mockRequest = injectNullExportContext(mockRequest)
		handler := http.HandlerFunc(ExportAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")

	})
	t.Run("ChildAsset_NotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		assetC0A0 := openTestJSON(assetL0Location)
		assetC1A1 := openTestJSON(assetL1Location)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C0", "A0")).
			Return(assetC0A0, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(nil, helpers.ErrNotFound)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		mockRequest = injectNullExportContext(mockRequest)
		handler := http.HandlerFunc(ExportAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)
		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
	t.Run("ChildAsset_Bad_JSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		assetC0A0 := openTestJSON(assetL0Location)
		assetC1A1 := openTestJSON(assetL1Location)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C0", "A0")).
			Return(assetC0A0, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(ioutil.NopCloser(strings.NewReader(invalidJSON)), nil)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		mockRequest = injectNullExportContext(mockRequest)
		handler := http.HandlerFunc(ExportAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)
		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
}

// TestExportWithAgentErrorConditions contains the test for agent failures
func TestExportWithAgentErrorConditions(t *testing.T) {
	t.Run("Agent_Failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(nil, errors.New(""))

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		mockRequest = injectNullExportContext(mockRequest)
		handler := http.HandlerFunc(ExportAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")

	})
	t.Run("Agent_Unauthorized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(nil, helpers.ErrUnauthorized)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		mockRequest = injectNullExportContext(mockRequest)
		handler := http.HandlerFunc(ExportAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 UNAUTHORIZED")

	})
	t.Run("ParentAgent_Failure", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		assetC1A1 := openTestJSON(assetL1Location)
		assetC2A2 := openTestJSON(assetL2A1Location)
		assetC3A3 := openTestJSON(assetL2A2Location)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C0", "A0")).
			Return(nil, errors.New(""))
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(assetC2A2, nil).MaxTimes(1)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C3", "A3")).
			Return(assetC3A3, nil).MaxTimes(1)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		mockRequest = injectNullExportContext(mockRequest)
		handler := http.HandlerFunc(ExportAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")

	})
	t.Run("ChildAgent_Failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		assetC0A0 := openTestJSON(assetL0Location)
		assetC1A1 := openTestJSON(assetL1Location)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C0", "A0")).
			Return(assetC0A0, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(nil, errors.New(""))

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		mockRequest = injectNullExportContext(mockRequest)
		handler := http.HandlerFunc(ExportAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
}

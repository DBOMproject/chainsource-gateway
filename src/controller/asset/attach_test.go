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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var attachRequestLocation = "../../testdata/asset_controller_tests/attach/attachRequest.json"
var attachRequestSelfLocation = "../../testdata/asset_controller_tests/attach/attachRequestToSelf.json"
var attachRequestMRLocation = "../../testdata/asset_controller_tests/attach/attachRequestMultiRepo.json"
var attachParentLocation = "../../testdata/asset_controller_tests/attach/parentAsset.json"
var attachChildLocation = "../../testdata/asset_controller_tests/attach/childAsset.json"
var attachParentDoneLocation = "../../testdata/asset_controller_tests/attach/parentAssetDone.json"

// TestAttachOnHappyPath contains the single repo test for attach
func TestAttachOnHappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	requestBody := openTestJSON(attachRequestLocation)

	assetC1A1 := openTestJSON(attachParentLocation)
	assetC2A2 := openTestJSON(attachChildLocation)

	var finalParentAsset helpers.Asset
	var finalChildAsset helpers.Asset

	mockAgent.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C1", "A1")).
		Return(assetC1A1, nil)
	mockAgent.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C2", "A2")).
		Return(assetC2A2, nil)
	mockAgent.EXPECT().Commit(gomock.Any(),
		mocks.AgentCommitTo("C1", "A1", "ATTACH")).
		Do(func(ctx context.Context, args agent.CommitArgs) {
			finalParentAsset = args.Payload
		}).
		Return(getAgentSuccessResponse(), nil)
	mockAgent.EXPECT().Commit(gomock.Any(),
		mocks.AgentCommitTo("C2", "A2", "ATTACH")).
		Do(func(ctx context.Context, args agent.CommitArgs) {
			finalChildAsset = args.Payload
		}).
		Return(getAgentSuccessResponse(), nil)

	providerMap := make(map[string]agent.Agent)
	providerMap["T1"] = mockAgent

	mockRequest := httptest.NewRequest("POST", "/", requestBody)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
		mocks.NewMockAssetSchemaAlwaysValid(ctrl))
	mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
	handler := http.HandlerFunc(AttachSubasset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	expectedChildLinkElement := helpers.AssetLinkElement{
		AssetElement: helpers.AssetElement{
			RepoID:    "T1",
			ChannelID: "C2",
			AssetID:   "A2",
			Asset:     nil,
		},
		Role:    "A_Role",
		SubRole: "A_SubRole",
	}
	expectedParentLinkElement := helpers.AssetLinkElement{
		AssetElement: helpers.AssetElement{
			RepoID:    "T1",
			ChannelID: "C1",
			AssetID:   "A1",
			Asset:     nil,
		},
		Role:    "A_Role",
		SubRole: "A_SubRole",
	}
	assert.Equal(t, 1, len(finalParentAsset.AttachedChildren), "Parent should have 1 child")
	assert.Equal(t, expectedChildLinkElement, finalParentAsset.AttachedChildren[0], "Parent should have correct child")
	assert.Equal(t, expectedParentLinkElement, *finalChildAsset.ParentAsset, "Child should have correct parent")
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode, "Response Should be 200 OK")
}

// TestAttachOnHappyPathOverMultipleRepos contains the multiple repo test for attach
func TestAttachOnHappyPathOverMultipleRepos(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgentT1 := mocks.NewMockAgent(ctrl)
	mockAgentT2 := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	requestBody := openTestJSON(attachRequestMRLocation)

	assetC1A1 := openTestJSON(attachParentLocation)
	assetC2A2 := openTestJSON(attachChildLocation)

	var finalParentAsset helpers.Asset
	var finalChildAsset helpers.Asset

	mockAgentT1.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C1", "A1")).
		Return(assetC1A1, nil)
	mockAgentT2.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C2", "A2")).
		Return(assetC2A2, nil)
	mockAgentT1.EXPECT().Commit(gomock.Any(),
		mocks.AgentCommitTo("C1", "A1", "ATTACH")).
		Do(func(ctx context.Context, args agent.CommitArgs) {
			finalParentAsset = args.Payload
		}).
		Return(getAgentSuccessResponse(), nil)
	mockAgentT2.EXPECT().Commit(gomock.Any(),
		mocks.AgentCommitTo("C2", "A2", "ATTACH")).
		Do(func(ctx context.Context, args agent.CommitArgs) {
			finalChildAsset = args.Payload
		}).
		Return(getAgentSuccessResponse(), nil)

	providerMap := make(map[string]agent.Agent)
	providerMap["T1"] = mockAgentT1
	providerMap["T2"] = mockAgentT2

	mockRequest := httptest.NewRequest("POST", "/", requestBody)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgentT1,
		mocks.NewMockAssetSchemaAlwaysValid(ctrl))
	mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
	handler := http.HandlerFunc(AttachSubasset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	expectedChildLinkElement := helpers.AssetLinkElement{
		AssetElement: helpers.AssetElement{
			RepoID:    "T2",
			ChannelID: "C2",
			AssetID:   "A2",
			Asset:     nil,
		},
		Role:    "A_Role",
		SubRole: "A_SubRole",
	}
	expectedParentLinkElement := helpers.AssetLinkElement{
		AssetElement: helpers.AssetElement{
			RepoID:    "T1",
			ChannelID: "C1",
			AssetID:   "A1",
			Asset:     nil,
		},
		Role:    "A_Role",
		SubRole: "A_SubRole",
	}
	assert.Equal(t, len(finalParentAsset.AttachedChildren), 1, "Parent should have 1 child")
	assert.Equal(t, expectedChildLinkElement, finalParentAsset.AttachedChildren[0], "Parent should have correct child")
	assert.Equal(t, expectedParentLinkElement, *finalChildAsset.ParentAsset, "Child should have correct parent")
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode, "Response Should be 200 OK")
}

// TestAttachWithAssetErrorConditions contains the tests that simulate error conditions of the asset
func TestAttachWithAssetErrorConditions(t *testing.T) {
	t.Run("Already_Attached", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(attachRequestLocation)

		assetC1A1 := openTestJSON(attachParentDoneLocation)
		assetC2A2 := openTestJSON(attachChildLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(assetC2A2, nil)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(AttachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusForbidden, responseRecorder.Result().StatusCode, "Response Should be 403 FORBIDDEN")
	})
	t.Run("Parent_Does_Not_Exist", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(attachRequestLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(nil, helpers.ErrNotFound)


		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(AttachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
	t.Run("Parent_Invalid_Format_On_Repo", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()
		assetC2A2 := openTestJSON(attachChildLocation)

		requestBody := openTestJSON(attachRequestLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(ioutil.NopCloser(strings.NewReader(invalidJSON)), nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(assetC2A2, nil).MaxTimes(1)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(AttachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
	t.Run("Child_Does_Not_Exist", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(attachRequestLocation)

		assetC1A1 := openTestJSON(attachParentLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(nil, helpers.ErrNotFound)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(AttachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
	t.Run("Child_Invalid_Format_On_Repo", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(attachRequestLocation)

		assetC1A1 := openTestJSON(attachParentLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(ioutil.NopCloser(strings.NewReader(invalidJSON)), nil)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(AttachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
}

// TestAttachWithRequestErrorConditions contains the tests that simulate user request error conditions
func TestAttachWithRequestErrorConditions(t *testing.T) {
	t.Run("Invalid_Request", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(attachRequestLocation)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent
		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysInvalid(ctrl))
		handler := http.HandlerFunc(AttachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode, "Response Should be 400 BAD REQUEST")
	})
	t.Run("Validator_Failure", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(attachRequestLocation)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent
		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysFailure(ctrl))
		handler := http.HandlerFunc(AttachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode, "Response Should be 500 INTERNAL SERVER ERROR")
	})
	t.Run("Self_Attach", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(attachRequestSelfLocation)

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		handler := http.HandlerFunc(AttachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusConflict, responseRecorder.Result().StatusCode, "Response Should be 409 CONFLICT")
	})
	t.Run("Non_JSON", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent
		mockRequest := httptest.NewRequest("POST", "/", strings.NewReader(invalidJSON))
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		handler := http.HandlerFunc(AttachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode, "Response Should be 400 BAD REQUEST")
	})
}

// TestAttachWithAgentErrorConditions contains the test for agent failures
func TestAttachWithAgentErrorConditions(t *testing.T) {
	t.Run("Parent_Query_Failure", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(attachRequestLocation)

		// assetC1A1 := openTestJSON(attachParentLocation)
		//assetC2A2 := openTestJSON(attachChildLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(nil, errors.New(""))
		//mockAgent.EXPECT().QueryStream(gomock.Any(),
		//	mocks.AgentQueryFor("C2", "A2")).
		//	Return(assetC2A2, nil)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(AttachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
	t.Run("Parent_Query_Unauthorized", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(attachRequestLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(nil, helpers.ErrUnauthorized)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(AttachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 UNAUTHORIZED")
	})
	t.Run("Parent_Commit_Failure", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(attachRequestLocation)

		assetC1A1 := openTestJSON(attachParentLocation)
		assetC2A2 := openTestJSON(attachChildLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(assetC2A2, nil)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C1", "A1", "ATTACH")).
			Return(nil, errors.New(""))

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(AttachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
	t.Run("Parent_Commit_Unauthorized", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(attachRequestLocation)

		assetC1A1 := openTestJSON(attachParentLocation)
		assetC2A2 := openTestJSON(attachChildLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(assetC2A2, nil)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C1", "A1", "ATTACH")).
			Return(nil, helpers.ErrUnauthorized)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(AttachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 UNAUTHORIZED")
	})
	t.Run("Child_Query_Failure",func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(attachRequestLocation)

		assetC1A1 := openTestJSON(attachParentLocation)
		//assetC2A2 := openTestJSON(attachChildLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(nil, errors.New(""))

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(AttachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
	t.Run("Child_Query_Unauthorized",func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(attachRequestLocation)

		assetC1A1 := openTestJSON(attachParentLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(nil, helpers.ErrUnauthorized)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(AttachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 UNAUTHORIZED")
	})
	t.Run("Child_Commit_Failure",func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(attachRequestLocation)

		assetC1A1 := openTestJSON(attachParentLocation)
		assetC2A2 := openTestJSON(attachChildLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(assetC2A2, nil)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C1", "A1", "ATTACH")).
			Return(getAgentSuccessResponse(), nil)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C2", "A2", "ATTACH")).
			Return(nil, errors.New(""))

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(AttachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
	t.Run("Child_Commit_Unauthorized",func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(attachRequestLocation)

		assetC1A1 := openTestJSON(attachParentLocation)
		assetC2A2 := openTestJSON(attachChildLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(assetC2A2, nil)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C1", "A1", "ATTACH")).
			Return(getAgentSuccessResponse(), nil)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C2", "A2", "ATTACH")).
			Return(nil, helpers.ErrUnauthorized)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(AttachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 UNAUTHORIZED")
	})
}

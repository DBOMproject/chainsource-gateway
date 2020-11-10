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

const detachRequestLocation = "../../testdata/asset_controller_tests/detach/detachRequest.json"
const detachRequestNLLocation = "../../testdata/asset_controller_tests/detach/detachRequestNotLinked.json"
const detachRequestMRLocation = "../../testdata/asset_controller_tests/detach/detachRequestMultiRepo.json"
const detachParentLocation = "../../testdata/asset_controller_tests/detach/parentAsset.json"
const detachParentChildlessLocation = "../../testdata/asset_controller_tests/detach/parentAssetChildless.json"
const detachParentNoChildAssetLinkLocation = "../../testdata/asset_controller_tests/detach/parentAssetWithoutSpecificChild.json"
const detachChildLocation = "../../testdata/asset_controller_tests/detach/childAsset.json"

// TestDetachOnHappyPath contains the single repo test for detach
func TestDetachOnHappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	requestBody := openTestJSON(detachRequestLocation)
	assetC1A1 := openTestJSON(detachParentLocation)
	assetC2A2 := openTestJSON(detachChildLocation)

	var finalParentAsset helpers.Asset
	var finalChildAsset helpers.Asset

	mockAgent.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C1", "A1")).
		Return(assetC1A1, nil).AnyTimes()
	mockAgent.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C2", "A2")).
		Return(assetC2A2, nil).AnyTimes()
	mockAgent.EXPECT().Commit(gomock.Any(),
		mocks.AgentCommitTo("C1", "A1", "DETACH")).
		Do(func(ctx context.Context, args agent.CommitArgs) {
			finalParentAsset = args.Payload
		}).
		Return(getAgentSuccessResponse(), nil).AnyTimes()
	mockAgent.EXPECT().Commit(gomock.Any(),
		mocks.AgentCommitTo("C2", "A2", "DETACH")).
		Do(func(ctx context.Context, args agent.CommitArgs) {
			finalChildAsset = args.Payload
		}).
		Return(getAgentSuccessResponse(), nil).AnyTimes()

	providerMap := make(map[string]agent.Agent)
	providerMap["T1"] = mockAgent

	mockRequest := httptest.NewRequest("POST", "/", requestBody)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
		mocks.NewMockAssetSchemaAlwaysValid(ctrl))
	mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
	handler := http.HandlerFunc(DetachSubasset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	expectedChildLinkElement := helpers.AssetLinkElement{
		AssetElement: helpers.AssetElement{
			RepoID:    "T1",
			ChannelID: "C2",
			AssetID:   "A3",
			Asset:     nil,
		},
		Role:    "A_Role2",
		SubRole: "A_SubRole2",
	}
	expectedParentLinkElement := helpers.AssetLinkElement{
		AssetElement: helpers.AssetElement{
			RepoID:    "",
			ChannelID: "",
			AssetID:   "",
			Asset:     nil,
		},
		Role:    "",
		SubRole: "",
	}
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode, "Response Should be 200 OK")
	assert.Equal(t, 1, len(finalParentAsset.AttachedChildren), "Parent should have 1 child")
	assert.Equal(t, expectedChildLinkElement, finalParentAsset.AttachedChildren[0], "Parent should have correct child")
	assert.Equal(t, expectedParentLinkElement, *finalChildAsset.ParentAsset, "Child should have correct parent")

}

// TestDetachOnHappyPathOverMultiRepo contains the multi repo test for detach
func TestDetachOnHappyPathOverMultiRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgentT1 := mocks.NewMockAgent(ctrl)
	mockAgentT2 := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	requestBody := openTestJSON(detachRequestMRLocation)
	assetC1A1 := openTestJSON(detachParentLocation)
	assetC2A2 := openTestJSON(detachChildLocation)

	var finalParentAsset helpers.Asset
	var finalChildAsset helpers.Asset

	mockAgentT1.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C1", "A1")).
		Return(assetC1A1, nil).AnyTimes()
	mockAgentT2.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C2", "A2")).
		Return(assetC2A2, nil).AnyTimes()
	mockAgentT1.EXPECT().Commit(gomock.Any(),
		mocks.AgentCommitTo("C1", "A1", "DETACH")).
		Do(func(ctx context.Context, args agent.CommitArgs) {
			finalParentAsset = args.Payload
		}).
		Return(getAgentSuccessResponse(), nil).AnyTimes()
	mockAgentT2.EXPECT().Commit(gomock.Any(),
		mocks.AgentCommitTo("C2", "A2", "DETACH")).
		Do(func(ctx context.Context, args agent.CommitArgs) {
			finalChildAsset = args.Payload
		}).
		Return(getAgentSuccessResponse(), nil).AnyTimes()

	providerMap := make(map[string]agent.Agent)
	providerMap["T1"] = mockAgentT1
	providerMap["T2"] = mockAgentT2

	mockRequest := httptest.NewRequest("POST", "/", requestBody)
	responseRecorder := httptest.NewRecorder()
	mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgentT1,
		mocks.NewMockAssetSchemaAlwaysValid(ctrl))
	mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
	handler := http.HandlerFunc(DetachSubasset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	expectedChildLinkElement := helpers.AssetLinkElement{
		AssetElement: helpers.AssetElement{
			RepoID:    "T1",
			ChannelID: "C2",
			AssetID:   "A3",
			Asset:     nil,
		},
		Role:    "A_Role2",
		SubRole: "A_SubRole2",
	}
	expectedParentLinkElement := helpers.AssetLinkElement{
		AssetElement: helpers.AssetElement{
			RepoID:    "",
			ChannelID: "",
			AssetID:   "",
			Asset:     nil,
		},
		Role:    "",
		SubRole: "",
	}
	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode, "Response Should be 200 OK")
	assert.Equal(t, 1, len(finalParentAsset.AttachedChildren), "Parent should have 1 child")
	assert.Equal(t, expectedChildLinkElement, finalParentAsset.AttachedChildren[0], "Parent should have correct child")
	assert.Equal(t, expectedParentLinkElement, *finalChildAsset.ParentAsset, "Child should have correct parent")

}

// TestDetachWithAssetErrorConditions contains the tests that simulate error conditions of the asset
func TestDetachWithAssetErrorConditions(t *testing.T) {
	t.Run("Childless", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(detachRequestLocation)
		assetC1A1 := openTestJSON(detachParentChildlessLocation)
		assetC2A2 := openTestJSON(detachChildLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
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
		handler := http.HandlerFunc(DetachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusForbidden, responseRecorder.Result().StatusCode, "Response Should be 403 FORBIDDEN")
	})
	t.Run("Child_Not_Linked", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(detachRequestNLLocation)
		assetC1A1 := openTestJSON(detachParentLocation)
		assetC2A6 := openTestJSON(detachChildLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A6")).
			Return(assetC2A6, nil).MaxTimes(1)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(DetachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusForbidden, responseRecorder.Result().StatusCode, "Response Should be 403 FORBIDDEN")
	})
	t.Run("Parent_Does_Not_Exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(detachRequestLocation)
		assetC2A2 := openTestJSON(detachChildLocation)

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
		handler := http.HandlerFunc(DetachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")

	})
	t.Run("Parent_Invalid_Format_On_Repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(detachRequestLocation)
		assetC2A2 := openTestJSON(detachChildLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(nil, helpers.ErrNotFound)
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
		handler := http.HandlerFunc(DetachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")

	})
	t.Run("Child_Does_Not_Exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(detachRequestLocation)
		assetC1A1 := openTestJSON(detachParentLocation)

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
		handler := http.HandlerFunc(DetachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")

	})
	t.Run("Child_Invalid_Format_On_Repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(detachRequestLocation)
		assetC1A1 := openTestJSON(detachParentNoChildAssetLinkLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil).AnyTimes()
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(ioutil.NopCloser(strings.NewReader(invalidJSON)), nil).AnyTimes()

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(DetachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")

	})
}

// TestDetachWithRequestErrorConditions contains the tests that simulate user request error conditions
func TestDetachWithRequestErrorConditions(t *testing.T) {
	t.Run("Invalid_Request", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(detachRequestLocation)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysInvalid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(DetachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode, "Response Should be 400 BAD REQUEST")

	})
	t.Run("Validator_Failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(detachRequestLocation)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysFailure(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(DetachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode, "Response Should be 500 INTERNAL SERVER ERROR")

	})
	t.Run("Non_JSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", strings.NewReader(invalidJSON))
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(DetachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode, "Response Should be 400 BAD REQUEST")
	})
}

// TestDetachWithAgentErrorConditions contains the test for agent failures
func TestDetachWithAgentErrorConditions(t *testing.T) {
	t.Run("Parent_Query_Failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(detachRequestLocation)
		assetC1A1 := openTestJSON(detachParentLocation)
		//assetC2A2 := openTestJSON(detachChildLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, errors.New(""))


		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(DetachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")

	})
	t.Run("Parent_Query_Unauthorized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(detachRequestLocation)
		assetC1A1 := openTestJSON(detachParentLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, helpers.ErrUnauthorized)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(DetachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 UNAUTHORIZED")

	})
	t.Run("Parent_Commit_Failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(detachRequestLocation)

		assetC1A1 := openTestJSON(detachParentLocation)
		assetC2A2 := openTestJSON(detachChildLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(assetC2A2, nil)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C1", "A1", "DETACH")).
			Return(nil, errors.New(""))

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(DetachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
	t.Run("Parent_Commit_Unauthorized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(detachRequestLocation)

		assetC1A1 := openTestJSON(detachParentLocation)
		assetC2A2 := openTestJSON(detachChildLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(assetC2A2, nil)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C1", "A1", "DETACH")).
			Return(nil, helpers.ErrUnauthorized)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(DetachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 UNAUTHORIZED")
	})
	t.Run("Child_Query_Failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(detachRequestLocation)
		assetC1A1 := openTestJSON(detachParentLocation)

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
		handler := http.HandlerFunc(DetachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")

	})
	t.Run("Child_Query_Unauthorized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(detachRequestLocation)
		assetC1A1 := openTestJSON(detachParentLocation)

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
		handler := http.HandlerFunc(DetachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 UNAUTHORIZED")

	})
	t.Run("Child_Commit_Failure", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(detachRequestLocation)
		assetC1A1 := openTestJSON(detachParentLocation)
		assetC2A2 := openTestJSON(detachChildLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(assetC2A2, nil)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C1", "A1", "DETACH")).
			Return(getAgentSuccessResponse(), nil)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C2", "A2", "DETACH")).
			Return(nil, errors.New(""))

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(DetachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")

	})
	t.Run("Child_Commit_Unauthorized", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(detachRequestLocation)
		assetC1A1 := openTestJSON(detachParentLocation)
		assetC2A2 := openTestJSON(detachChildLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(assetC2A2, nil)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C1", "A1", "DETACH")).
			Return(getAgentSuccessResponse(), nil)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C2", "A2", "DETACH")).
			Return(nil, helpers.ErrUnauthorized)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(DetachSubasset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 UNAUTHORIZED")

	})
}

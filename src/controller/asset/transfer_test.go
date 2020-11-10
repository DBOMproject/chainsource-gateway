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
const transferRequestLocation = "../../testdata/asset_controller_tests/transfer/transferRequest.json"
const transferredNeverAssetLocation = "../../testdata/asset_controller_tests/transfer/transferredNever.json"
const transferRequestMRLocation = "../../testdata/asset_controller_tests/transfer/transferRequestMultiRepo.json"
const transferredOnceAssetLocation = "../../testdata/asset_controller_tests/transfer/transferredOnce.json"

// TestTransferOnHappyPath is the single repo happy path for transfer
func TestTransferOnHappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	requestBody := openTestJSON(transferRequestLocation)

	assetC1A1 := openTestJSON(transferredOnceAssetLocation)

	var finalParentAsset helpers.Asset
	var finalChildAsset helpers.Asset

	mockAgent.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C1", "A1")).
		Return(assetC1A1, nil)
	mockAgent.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C2", "A2")).
		Return(nil, helpers.ErrNotFound)
	mockAgent.EXPECT().Commit(gomock.Any(),
		mocks.AgentCommitTo("C1", "A1", "TRANSFER-OUT")).
		Do(func(ctx context.Context, args agent.CommitArgs) {
			finalParentAsset = args.Payload
		}).
		Return(getAgentSuccessResponse(), nil)
	mockAgent.EXPECT().Commit(gomock.Any(),
		mocks.AgentCommitTo("C2", "A2", "TRANSFER-IN")).
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
	handler := http.HandlerFunc(TransferAsset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	expectedCustodyTransferEvent := helpers.CustodyTransferEvent{
		Timestamp:            "",
		TransferDescription:  "sold",
		SourceRepoID:         "T1",
		SourceChannelID:      "C1",
		SourceAssetID:        "A1",
		DestinationRepoID:    "T1",
		DestinationChannelID: "C2",
		DestinationAssetID:   "A2",
	}

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode, "Response Should be 200 OK")

	assert.Equal(t, 2, len(finalParentAsset.CustodyTransferEvents), "Parent should have 2 Custody Transfer Events")
	assert.NotEqual(t, "", finalParentAsset.CustodyTransferEvents[1].Timestamp, "Parent Custody Transfer Event with Timestamp")
	finalParentAsset.CustodyTransferEvents[1].Timestamp = ""
	assert.Equal(t, expectedCustodyTransferEvent, finalParentAsset.CustodyTransferEvents[1], "Parent must have expected Custody Transfer Events")
	assert.True(t, finalParentAsset.ReadOnly, "Parent asset to be marked ReadOnly")

	assert.Equal(t, 2, len(finalChildAsset.CustodyTransferEvents), "Child should have 2 Custody Transfer Events")
	assert.NotEqual(t, "", finalChildAsset.CustodyTransferEvents[0].Timestamp, "Child Custody Transfer Event with Timestamp")
	finalChildAsset.CustodyTransferEvents[1].Timestamp = ""
	assert.Equal(t, expectedCustodyTransferEvent, finalChildAsset.CustodyTransferEvents[1], "Child must have expected Custody Transfer Events")
	assert.False(t, finalChildAsset.ReadOnly, "Child asset not to be marked ReadOnly")
}

// TestTransferOnHappyPath is the single repo happy path for transfer with an asset that has never been transferred
func TestTransferOnHappyPathWithNeverTransferredAsset(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgent := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	requestBody := openTestJSON(transferRequestLocation)

	assetC1A1 := openTestJSON(transferredNeverAssetLocation)

	var finalParentAsset helpers.Asset
	var finalChildAsset helpers.Asset

	mockAgent.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C1", "A1")).
		Return(assetC1A1, nil)
	mockAgent.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C2", "A2")).
		Return(nil, helpers.ErrNotFound)
	mockAgent.EXPECT().Commit(gomock.Any(),
		mocks.AgentCommitTo("C1", "A1", "TRANSFER-OUT")).
		Do(func(ctx context.Context, args agent.CommitArgs) {
			finalParentAsset = args.Payload
		}).
		Return(getAgentSuccessResponse(), nil)
	mockAgent.EXPECT().Commit(gomock.Any(),
		mocks.AgentCommitTo("C2", "A2", "TRANSFER-IN")).
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
	handler := http.HandlerFunc(TransferAsset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	expectedCustodyTransferEvent := helpers.CustodyTransferEvent{
		Timestamp:            "",
		TransferDescription:  "sold",
		SourceRepoID:         "T1",
		SourceChannelID:      "C1",
		SourceAssetID:        "A1",
		DestinationRepoID:    "T1",
		DestinationChannelID: "C2",
		DestinationAssetID:   "A2",
	}

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode, "Response Should be 200 OK")

	assert.Equal(t, 1, len(finalParentAsset.CustodyTransferEvents), "Parent should have 1 Custody Transfer Event")
	finalParentAsset.CustodyTransferEvents[0].Timestamp = ""
	assert.Equal(t, expectedCustodyTransferEvent, finalParentAsset.CustodyTransferEvents[0], "Parent must have expected Custody Transfer Events")
	assert.True(t, finalParentAsset.ReadOnly, "Parent asset to be marked ReadOnly")

	assert.Equal(t, 1, len(finalChildAsset.CustodyTransferEvents), "Child should have 1 Custody Transfer Event")
	finalChildAsset.CustodyTransferEvents[0].Timestamp = ""
	assert.Equal(t, expectedCustodyTransferEvent, finalChildAsset.CustodyTransferEvents[0], "Child must have expected Custody Transfer Events")
	assert.False(t, finalChildAsset.ReadOnly, "Child asset not to be marked ReadOnly")
}

// TestTransferOnHappyPath is the multi repo happy path for transfer
func TestTransferOnHappyPathOverMultipleRepos(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAgentT1 := mocks.NewMockAgent(ctrl)
	mockAgentT2 := mocks.NewMockAgent(ctrl)
	defer ctrl.Finish()

	requestBody := openTestJSON(transferRequestMRLocation)

	assetC1A1 := openTestJSON(transferredOnceAssetLocation)

	var finalParentAsset helpers.Asset
	var finalChildAsset helpers.Asset

	mockAgentT1.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C1", "A1")).
		Return(assetC1A1, nil)
	mockAgentT2.EXPECT().QueryStream(gomock.Any(),
		mocks.AgentQueryFor("C2", "A2")).
		Return(nil, helpers.ErrNotFound)
	mockAgentT1.EXPECT().Commit(gomock.Any(),
		mocks.AgentCommitTo("C1", "A1", "TRANSFER-OUT")).
		Do(func(ctx context.Context, args agent.CommitArgs) {
			finalParentAsset = args.Payload
		}).
		Return(getAgentSuccessResponse(), nil)
	mockAgentT2.EXPECT().Commit(gomock.Any(),
		mocks.AgentCommitTo("C2", "A2", "TRANSFER-IN")).
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
	handler := http.HandlerFunc(TransferAsset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	expectedCustodyTransferEvent := helpers.CustodyTransferEvent{
		Timestamp:            "",
		TransferDescription:  "sold",
		SourceRepoID:         "T1",
		SourceChannelID:      "C1",
		SourceAssetID:        "A1",
		DestinationRepoID:    "T2",
		DestinationChannelID: "C2",
		DestinationAssetID:   "A2",
	}

	assert.Equal(t, http.StatusOK, responseRecorder.Result().StatusCode, "Response Should be 200 OK")

	assert.Equal(t, 2, len(finalParentAsset.CustodyTransferEvents), "Parent should have 2 Custody Transfer Events")
	assert.NotEqual(t, "", finalParentAsset.CustodyTransferEvents[1].Timestamp, "Parent Custody Transfer Event with Timestamp")
	finalParentAsset.CustodyTransferEvents[1].Timestamp = ""
	assert.Equal(t, expectedCustodyTransferEvent, finalParentAsset.CustodyTransferEvents[1], "Parent must have expected Custody Transfer Events")
	assert.True(t, finalParentAsset.ReadOnly, "Parent asset to be marked ReadOnly")

	assert.Equal(t, 2, len(finalChildAsset.CustodyTransferEvents), "Child should have 2 Custody Transfer Events")
	assert.NotEqual(t, "", finalChildAsset.CustodyTransferEvents[0].Timestamp, "Child Custody Transfer Event with Timestamp")
	finalChildAsset.CustodyTransferEvents[1].Timestamp = ""
	assert.Equal(t, expectedCustodyTransferEvent, finalChildAsset.CustodyTransferEvents[1], "Child must have expected Custody Transfer Events")
	assert.False(t, finalChildAsset.ReadOnly, "Child asset not to be marked ReadOnly")
}

// TestTransferWithAssetErrorConditions contains the tests that simulate error conditions of the asset
func TestTransferWithAssetErrorConditions(t *testing.T) {
	t.Run("Asset_Exists_At_Destination", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(transferRequestLocation)

		assetC1A1 := openTestJSON(transferredOnceAssetLocation)
		assetC2A2 := openTestJSON(transferredOnceAssetLocation)

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
		handler := http.HandlerFunc(TransferAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusConflict, responseRecorder.Result().StatusCode, "Response Should be 409 CONFLICT")

	})
	t.Run("Source_Does_Not_Exist", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(transferRequestLocation)

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
		handler := http.HandlerFunc(TransferAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
	t.Run("Source_Invalid_JSON", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(transferRequestLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(ioutil.NopCloser(strings.NewReader(invalidJSON)), nil)
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
		handler := http.HandlerFunc(TransferAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
}

// TestTransferWithRequestErrorConditions contains the tests that simulate user request error conditions
func TestTransferWithRequestErrorConditions(t *testing.T) {
	t.Run("Invalid_Request", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(transferRequestLocation)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysInvalid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(TransferAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode, "Response Should be 400 BAD REQUEST")
	})
	t.Run("Validator_Failure", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(transferRequestLocation)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysFailure(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(TransferAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusInternalServerError, responseRecorder.Result().StatusCode, "Response Should be 500 INTERNAL SERVER ERROR")
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
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(TransferAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Result().StatusCode, "Response Should be 400 BAD REQUEST")
	})
}

// TestTransferWithAgentErrorConditions contains the test for agent failures
func TestTransferWithAgentErrorConditions(t *testing.T) {
	t.Run("Source_Query_Failure", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(transferRequestLocation)

		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(nil, errors.New(""))

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(TransferAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
	t.Run("Source_Query_Unauthorized", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(transferRequestLocation)

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
		handler := http.HandlerFunc(TransferAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 UNAUTHORIZED")
	})
	t.Run("Destination_Query_Failure",func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(transferRequestLocation)

		assetC1A1 := openTestJSON(transferredOnceAssetLocation)

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
		handler := http.HandlerFunc(TransferAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
	t.Run("Destination_Query_Unauthorized",func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(transferRequestLocation)

		assetC1A1 := openTestJSON(transferredOnceAssetLocation)

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
		handler := http.HandlerFunc(TransferAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 UNAUTHORIZED")
	})
	t.Run("Destination_Commit_Failure", func (t *testing.T) {	ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(transferRequestLocation)

		assetC1A1 := openTestJSON(transferredOnceAssetLocation)


		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(nil, helpers.ErrNotFound)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C1", "A1", "TRANSFER-OUT")).
			Return(nil, errors.New(""))
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C2", "A2", "TRANSFER-IN")).
			Return(getAgentSuccessResponse(), nil).MaxTimes(1)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(TransferAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
	t.Run("Destination_Commit_Unauthorized", func (t *testing.T) {	ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(transferRequestLocation)

		assetC1A1 := openTestJSON(transferredOnceAssetLocation)


		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(nil, helpers.ErrNotFound)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C1", "A1", "TRANSFER-OUT")).
			Return(nil, helpers.ErrUnauthorized)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C2", "A2", "TRANSFER-IN")).
			Return(getAgentSuccessResponse(), nil).MaxTimes(1)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(TransferAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 UNAUTHORIZED")
	})
	t.Run("Destination_Commit_Failure",func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(transferRequestLocation)

		assetC1A1 := openTestJSON(transferredOnceAssetLocation)


		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(nil, helpers.ErrNotFound)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C1", "A1", "TRANSFER-OUT")).
			Return(getAgentSuccessResponse(), nil)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C2", "A2", "TRANSFER-IN")).
			Return(nil, errors.New(""))

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(TransferAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusBadGateway, responseRecorder.Result().StatusCode, "Response Should be 502 BAD GATEWAY")
	})
	t.Run("Destination_Commit_Unauthorized",func (t *testing.T) {
		ctrl := gomock.NewController(t)
		mockAgent := mocks.NewMockAgent(ctrl)
		defer ctrl.Finish()

		requestBody := openTestJSON(transferRequestLocation)

		assetC1A1 := openTestJSON(transferredOnceAssetLocation)


		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C1", "A1")).
			Return(assetC1A1, nil)
		mockAgent.EXPECT().QueryStream(gomock.Any(),
			mocks.AgentQueryFor("C2", "A2")).
			Return(nil, helpers.ErrNotFound)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C1", "A1", "TRANSFER-OUT")).
			Return(getAgentSuccessResponse(), nil)
		mockAgent.EXPECT().Commit(gomock.Any(),
			mocks.AgentCommitTo("C2", "A2", "TRANSFER-IN")).
			Return(nil, helpers.ErrUnauthorized)

		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent

		mockRequest := httptest.NewRequest("POST", "/", requestBody)
		responseRecorder := httptest.NewRecorder()
		mockRequest = injectMockAssetContext(mockRequest, "T1", "C1", "A1", mockAgent,
			mocks.NewMockAssetSchemaAlwaysValid(ctrl))
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)
		handler := http.HandlerFunc(TransferAsset)
		handler.ServeHTTP(responseRecorder, mockRequest)

		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Result().StatusCode, "Response Should be 401 UNAUTHORIZED")
	})
}

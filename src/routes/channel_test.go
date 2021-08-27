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

package routes

import (
	"chainsource-gateway/agent"
	"chainsource-gateway/helpers"
	"chainsource-gateway/mocks"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Test_ChannelSubRouting tests if the channel sub router mounts successfully
func Test_ChannelSubRouting(t *testing.T) {
	assert.NotPanics(t, func() {
		channelSubRouting(chi.NewRouter())
	}, "Router mounts without panic")

}

// Test_ChannelContext tests if the channelContext is injected
func Test_ChannelContext(t *testing.T) {

	t.Run("When_Request_OK_And_Repo_Exists", func(t *testing.T) {
		mockRequest := httptest.NewRequest("POST", "/", strings.NewReader(""))
		ctrl := gomock.NewController(t)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("repoID", "T1")
		rctx.URLParams.Add("channelID", "C1")

		mockAgent := mocks.NewMockAgent(ctrl)
		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)

		mockRequest = mockRequest.WithContext(context.WithValue(mockRequest.Context(), chi.RouteCtxKey, rctx))
		responseRecorder := httptest.NewRecorder()
		channelContext(getContextAssertionMiddleware(func(ctx context.Context) {
			val := ctx.Value("assetVars")
			assert.NotNil(t, val, "assetVars must be injected")
			assert.IsType(t, helpers.AssetRoutingVars{}, val, "Is of type AssetRoutingVars")
		})).ServeHTTP(responseRecorder, mockRequest)
		assert.Equal(t, http.StatusOK, responseRecorder.Code, "A 200 OK is returned")
	})
	t.Run("When_Request_OK_And_Repo_Does_Not_Exist", func(t *testing.T) {
		mockRequest := httptest.NewRequest("GET", "/", strings.NewReader(""))
		ctrl := gomock.NewController(t)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("repoID", "NIL_REPO")
		rctx.URLParams.Add("channelID", "C1")

		providerMap := make(map[string]agent.Agent)
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)

		mockRequest = mockRequest.WithContext(context.WithValue(mockRequest.Context(), chi.RouteCtxKey, rctx))
		responseRecorder := httptest.NewRecorder()
		channelContext(nil).ServeHTTP(responseRecorder, mockRequest)
		assert.Equal(t, http.StatusNotFound, responseRecorder.Code, "A 404 NOT FOUND is returned")
	})
	t.Run("When_Request_Not_OK", func(t *testing.T) {
		mockRequest := httptest.NewRequest("GET", "/", strings.NewReader(""))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("repoID", "")
		rctx.URLParams.Add("channelID", "")

		mockRequest = mockRequest.WithContext(context.WithValue(mockRequest.Context(), chi.RouteCtxKey, rctx))
		responseRecorder := httptest.NewRecorder()
		channelContext(nil).ServeHTTP(responseRecorder, mockRequest)
		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code, "A 400 BAD REQUEST is returned")
	})

}

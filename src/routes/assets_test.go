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
	"chainsource-gateway/pgp"
	"chainsource-gateway/schema"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	jaeger "github.com/uber/jaeger-client-go"
)

// TestAssetRouter tests if the asset router initializes successfully
func TestAssetRouter(t *testing.T) {
	assert.NotPanics(t, func() {
		AssetRouter()
	}, "Router initializes without panic")
}

// Test_assetSubRouting tests if the asset sub router mounts successfully
func Test_assetSubRouting(t *testing.T) {
	assert.NotPanics(t, func() {
		assetSubRouting(chi.NewRouter())
	}, "Router mounts without panic")

}

// getContextAssertionMiddleware takes in a function that can access the context of a http request for assertions
func getContextAssertionMiddleware(assertionFn func(ctx context.Context)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertionFn(r.Context())
		w.WriteHeader(http.StatusOK)
	})

}

// Test_agentProvider tests if the agentProvider is injected
func Test_agentProvider(t *testing.T) {
	mockRequest := httptest.NewRequest("POST", "/", strings.NewReader(""))
	responseRecorder := httptest.NewRecorder()
	agentProvider(getContextAssertionMiddleware(func(ctx context.Context) {
		val := ctx.Value("agentProvider")
		assert.NotNil(t, val, "agentProvider must be injected")
		assert.IsType(t, &agent.HttpAgentProvider{}, val, "Is a HTTP Agent Provider")
	})).ServeHTTP(responseRecorder, mockRequest)
	assert.Equal(t, http.StatusOK, responseRecorder.Code, "A 200 OK is returned")
}

// Test_assetContext tests if the assetContext is injected
func Test_assetContext(t *testing.T) {

	t.Run("When_Request_OK_And_Repo_Exists", func(t *testing.T) {
		mockRequest := httptest.NewRequest("POST", "/", strings.NewReader(""))
		ctrl := gomock.NewController(t)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("repoID", "T1")
		rctx.URLParams.Add("channelID", "C1")
		rctx.URLParams.Add("assetID", "A1")

		mockAgent := mocks.NewMockAgent(ctrl)
		providerMap := make(map[string]agent.Agent)
		providerMap["T1"] = mockAgent
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)

		mockRequest = mockRequest.WithContext(context.WithValue(mockRequest.Context(), chi.RouteCtxKey, rctx))
		responseRecorder := httptest.NewRecorder()
		assetContext(getContextAssertionMiddleware(func(ctx context.Context) {
			val := ctx.Value("assetVars")
			assert.NotNil(t, val, "assetVars must be injected")
			assert.IsType(t, helpers.AssetRoutingVars{}, val, "Is of type AssetRoutingVars")
		})).ServeHTTP(responseRecorder, mockRequest)
		assert.Equal(t, http.StatusOK, responseRecorder.Code, "A 200 OK is returned")
	})
	t.Run("When_Request_OK_And_Repo_Does_Not_Exist", func(t *testing.T) {
		mockRequest := httptest.NewRequest("POST", "/", strings.NewReader(""))
		ctrl := gomock.NewController(t)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("repoID", "NIL_REPO")
		rctx.URLParams.Add("channelID", "C1")
		rctx.URLParams.Add("assetID", "A1")

		providerMap := make(map[string]agent.Agent)
		mockRequest = mocks.InjectAgentProviderIntoRequest(mockRequest, ctrl, providerMap)

		mockRequest = mockRequest.WithContext(context.WithValue(mockRequest.Context(), chi.RouteCtxKey, rctx))
		responseRecorder := httptest.NewRecorder()
		assetContext(nil).ServeHTTP(responseRecorder, mockRequest)
		assert.Equal(t, http.StatusNotFound, responseRecorder.Code, "A 404 NOT FOUND is returned")
	})
	t.Run("When_Request_Not_OK", func(t *testing.T) {
		mockRequest := httptest.NewRequest("POST", "/", strings.NewReader(""))

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("repoID", "")
		rctx.URLParams.Add("channelID", "")
		rctx.URLParams.Add("assetID", "")

		mockRequest = mockRequest.WithContext(context.WithValue(mockRequest.Context(), chi.RouteCtxKey, rctx))
		responseRecorder := httptest.NewRecorder()
		assetContext(nil).ServeHTTP(responseRecorder, mockRequest)
		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code, "A 400 BAD REQUEST is returned")
	})

}

// Test_assetSchemaValidator tests if the schemaValidator is injected
func Test_assetSchemaValidator(t *testing.T) {
	mockRequest := httptest.NewRequest("POST", "/", strings.NewReader(""))
	responseRecorder := httptest.NewRecorder()
	assetSchemaValidator(getContextAssertionMiddleware(func(ctx context.Context) {
		val := ctx.Value("schemaValidator")
		assert.NotNil(t, val, "schemaValidator must be injected")
		assert.Implements(t, (*schema.AssetSchema)(nil), val, "Implements AssetSchema")
	})).ServeHTTP(responseRecorder, mockRequest)
	assert.Equal(t, http.StatusOK, responseRecorder.Code, "A 200 OK is returned")
}

// Test_exportContext tests if the exportVars is injected
func Test_exportContext(t *testing.T) {
	mockRequest := httptest.NewRequest("POST", "/", strings.NewReader(""))
	responseRecorder := httptest.NewRecorder()
	exportContext(getContextAssertionMiddleware(func(ctx context.Context) {
		val := ctx.Value("exportVars")
		assert.NotNil(t, val, "exportVars must be injected")
		assert.IsType(t, helpers.ExportRoutingVars{}, val, "Is an ExportRoutingVars instance")
	})).ServeHTTP(responseRecorder, mockRequest)
	assert.Equal(t, http.StatusOK, responseRecorder.Code, "A 200 OK is returned")
}

// Test_queryContext tests if the assetQueryVars is injected
func Test_queryContext(t *testing.T) {
	mockRequest := httptest.NewRequest("POST", "/", strings.NewReader(""))
	responseRecorder := httptest.NewRecorder()
	queryContext(getContextAssertionMiddleware(func(ctx context.Context) {
		val := ctx.Value("assetQueryVars")
		assert.NotNil(t, val, "assetQueryVars must be injected")
		assert.IsType(t, map[string]interface{}{}, val, "Is an AssetQueryVars instance")
	})).ServeHTTP(responseRecorder, mockRequest)
	assert.Equal(t, http.StatusOK, responseRecorder.Code, "A 200 OK is returned")
}

// Test_injectSpanMiddleware tests if the Span is injected
func Test_injectSpanMiddleware(t *testing.T) {

	t.Run("When_No_Span_On_Wire", func(t *testing.T) {
		mockRequest := httptest.NewRequest("POST", "/", strings.NewReader(""))

		responseRecorder := httptest.NewRecorder()
		injectSpanMiddleware(getContextAssertionMiddleware(func(ctx context.Context) {
			assert.NotNil(t, opentracing.SpanFromContext(ctx), "The context must have an injected span")
		})).ServeHTTP(responseRecorder, mockRequest)
		assert.Equal(t, http.StatusOK, responseRecorder.Code, "A 200 OK is returned")
	})
	t.Run("When_Span_On_Wire", func(t *testing.T) {
		mockRequest := httptest.NewRequest("POST", "/", strings.NewReader(""))
		tracer, _ := jaeger.NewTracer("Test Service", jaeger.NewConstSampler(true), jaeger.NewNullReporter())
		noOpTracer := opentracing.GlobalTracer()
		opentracing.SetGlobalTracer(tracer)
		defer opentracing.SetGlobalTracer(noOpTracer)
		tracer.Inject(
			tracer.StartSpan("TestSpan").Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(mockRequest.Header))
		responseRecorder := httptest.NewRecorder()
		injectSpanMiddleware(getContextAssertionMiddleware(func(ctx context.Context) {
			assert.NotNil(t, opentracing.SpanFromContext(ctx), "The context must have an injected span")
		})).ServeHTTP(responseRecorder, mockRequest)
		assert.Equal(t, http.StatusOK, responseRecorder.Code, "A 200 OK is returned")
	})
}

// Test_signingServiceProvider tests if the signatureValidator is injected
func Test_signingServiceProvider(t *testing.T) {
	mockRequest := httptest.NewRequest("POST", "/", strings.NewReader(""))
	responseRecorder := httptest.NewRecorder()
	signingServiceProvider(getContextAssertionMiddleware(func(ctx context.Context) {
		val := ctx.Value("signatureValidator")
		assert.NotNil(t, val, "signatureValidator must be injected")
		assert.Implements(t, (*pgp.SignatureValidator)(nil), val, "Implements SignatureValidator")
	})).ServeHTTP(responseRecorder, mockRequest)
	assert.Equal(t, http.StatusOK, responseRecorder.Code, "A 200 OK is returned")
}

// Test_unmarshalBody tests if the JSONBody is injected
func Test_unmarshalBody(t *testing.T) {
	t.Run("When_POST_Or_PUT_And_Valid", func(t *testing.T) {
		mockRequest := httptest.NewRequest("POST", "/", strings.NewReader("{\"foo\": \"bar\"}"))
		responseRecorder := httptest.NewRecorder()
		unmarshalBody(getContextAssertionMiddleware(func(ctx context.Context) {
			val := ctx.Value("JSONBody")
			assert.NotNil(t, val, "JSONBody must be injected")
			assert.IsType(t, make(map[string]interface{}), val, "Is an interface map")
		})).ServeHTTP(responseRecorder, mockRequest)
		assert.Equal(t, http.StatusOK, responseRecorder.Code, "A 200 OK is returned")
	})
	t.Run("When_GET_And_Valid", func(t *testing.T) {
		mockRequest := httptest.NewRequest("GET", "/", nil)
		responseRecorder := httptest.NewRecorder()
		unmarshalBody(getContextAssertionMiddleware(func(ctx context.Context) {
			val := ctx.Value("JSONBody")
			assert.Nil(t, val, "JSONBody must not be injected")
		})).ServeHTTP(responseRecorder, mockRequest)
		assert.Equal(t, http.StatusOK, responseRecorder.Code, "A 200 OK is returned")
	})
	t.Run("When_POST_Or_PUT_And_Invalid_JSON", func(t *testing.T) {
		mockRequest := httptest.NewRequest("POST", "/", strings.NewReader("invalid.JSON"))
		responseRecorder := httptest.NewRecorder()
		unmarshalBody(nil).ServeHTTP(responseRecorder, mockRequest)
		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code, "A 400 BAD REQUEST is returned")
	})
	t.Run("When_POST_Or_PUT_And_No_Body", func(t *testing.T) {
		mockRequest := httptest.NewRequest("POST", "/", nil)
		responseRecorder := httptest.NewRecorder()
		unmarshalBody(nil).ServeHTTP(responseRecorder, mockRequest)
		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code, "A 400 BAD REQUEST is returned")
	})

}

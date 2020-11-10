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

// Package routes contains all the routes for the Gateway API
package routes

import (
	"bytes"
	"chainsource-gateway/agent"
	"chainsource-gateway/controller/asset"
	"chainsource-gateway/helpers"
	"chainsource-gateway/pgp"
	"chainsource-gateway/responses"
	"chainsource-gateway/schema"
	"chainsource-gateway/tracing"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"
)

var log = helpers.GetLogger("AssetRouter")

// AssetRouter defines the main routes for the Gateway API
func AssetRouter() (r chi.Router) {
	r = chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Route("/repo/{repoID}/chan/{channelID}/asset/{assetID}", assetSubRouting)
	return
}

// assetSubRouting defines the sub routes for the asset APIs
func assetSubRouting(r chi.Router) {
	r.Use(injectSpanMiddleware)
	r.Use(agentProvider)
	r.Use(signingServiceProvider)
	r.Use(assetContext)
	r.Use(assetSchemaValidator)
	r.Use(unmarshalBody)

	// Base CRUD
	r.Post("/", asset.CreateAsset)
	r.Get("/", asset.RetrieveAsset)
	r.Put("/", asset.UpdateAsset)
	r.Delete("/", asset.DeleteAsset)

	// Link/Unlink APIs
	r.Post("/attach", asset.AttachSubasset)
	r.Post("/detach", asset.DetachSubasset)
	r.Get("/trail", asset.AuditAsset)
	r.Post("/transfer", asset.TransferAsset)
	r.Get("/validate", asset.ValidateAsset)

	// Export API
	r.Group(func(r chi.Router) {
		r.Use(exportContext)
		r.Get("/export", asset.ExportAsset)
		r.Get("/export/{fileName}", asset.ExportAsset)
	})
}

// agentProvider injects an agent provider into the request context
func agentProvider(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, ctx := opentracing.StartSpanFromContext(r.Context(), "Embedding Agent Provider")
		provider := agent.NewHTTPAgentProvider()
		ctx = context.WithValue(r.Context(), "agentProvider", provider)
		span.Finish()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// signingServiceProvider injects a "SignatureValidator" that works based on the signing service
func signingServiceProvider(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, ctx := opentracing.StartSpanFromContext(r.Context(), "Embedding Signing Service Provider")
		provider := pgp.NewSigningServiceValidator()
		ctx = context.WithValue(r.Context(), "signatureValidator", provider)
		span.Finish()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// assetContext parses the context for the asset APIs
func assetContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, ctx := opentracing.StartSpanFromContext(r.Context(), "Embedding Asset Context")

		// Append context with routing details

		vars := helpers.AssetRoutingVars{
			RepoID:    chi.URLParam(r, "repoID"),
			ChannelID: chi.URLParam(r, "channelID"),
			AssetID:   chi.URLParam(r, "assetID"),
		}

		err := errors.New("Invalid URL")

		if vars.AssetID == "" || vars.ChannelID == "" || vars.RepoID == "" {
			render.Render(w, r, responses.ErrInvalidRequest(err))
			tracing.LogAndTraceErr(log, span, err, "Empty URL parameters found")
			return
		}

		agentProvider := r.Context().Value("agentProvider").(agent.Provider)
		log.Debug().Msgf("ctx: %+v ", vars)
		ctx = context.WithValue(r.Context(), "assetVars", vars)

		// Getting the agent-config for the repo we're trying to access
		agentConfig, err := agentProvider.GetAgentConfigForRepo(vars.RepoID)
		if err != nil {
			tracing.LogAndTraceErr(log, span, err, "No agent configured for request")
			_ = render.Render(w, r, responses.ErrNoAgent(err))
			return
		}

		requestAgent := agentProvider.NewAgent(&agentConfig)
		ctx = context.WithValue(ctx, "agent", requestAgent)
		span.Finish()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// exportContext parses the context for the export API
func exportContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		span, ctx := opentracing.StartSpanFromContext(r.Context(), "Embedding Export Context")

		// Append export file
		vars := helpers.ExportRoutingVars{
			FileName:       chi.URLParam(r, "fileName"),
			InlineResponse: r.URL.Query().Get("inlineResponse"),
		}

		log.Debug().Msgf("ctx: %+v ", vars)
		ctx = context.WithValue(r.Context(), "exportVars", vars)

		span.Finish()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// unmarshalBody unmarshals the body for the asset APIs
func unmarshalBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method == "POST" || r.Method == "PUT" {
			b, err := ioutil.ReadAll(r.Body)

			if err != nil {
				log.Err(err).Msgf("Failed to parse body of %s request", r.Method)
				_ = render.Render(w, r, responses.ErrInvalidRequest(err))
				return
			}
			r.Body = ioutil.NopCloser(bytes.NewBuffer(b))

			var result map[string]interface{}
			err = json.Unmarshal(b, &result)
			if err != nil {
				log.Err(err).Msgf("Invalid JSON in body")
				_ = render.Render(w, r, responses.ErrInvalidRequest(err))
				return
			}
			ctx = context.WithValue(ctx, "JSONBody", result)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// assetSchemaValidator injects an schema validator into the request context
func assetSchemaValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, ctx := opentracing.StartSpanFromContext(r.Context(), "Embedding Schema Validator")
		schemaValidator := schema.NewAssetSchemaImpl()
		ctx = context.WithValue(r.Context(), "schemaValidator", schemaValidator)
		span.Finish()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// injectSpanMiddleware injects an opentracing span for the asset APIs
func injectSpanMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var span opentracing.Span

		httpContext, err := opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			span = opentracing.GlobalTracer().StartSpan("Asset HTTP Request")
		} else {
			span = opentracing.GlobalTracer().StartSpan("Asset HTTP Request", opentracing.ChildOf(httpContext))
		}

		defer span.Finish()
		span.SetTag("method", r.Method)
		span.SetTag("go.request.id", middleware.GetReqID(r.Context()))
		span.SetTag("url", r.URL)

		ctx = opentracing.ContextWithSpan(ctx, span)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

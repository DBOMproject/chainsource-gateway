/*
 * Copyright 2023 Unisys Corporation
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
	"chainsource-gateway/controller/asset"
	"chainsource-gateway/helpers"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"
)

var assetLog = helpers.GetLogger("AssetRouter")

// AssetRouter defines the main routes for the asset API
func AssetRouter() (r chi.Router) {
	r = chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Route("/", assetSubRouting)
	return
}

// assetSubRouting defines the sub routes for the asset APIs
func assetSubRouting(r chi.Router) {
	r.Use(injectSpanMiddleware)

	// Base CRUD
	r.Get("/", asset.ListAssets)
	r.Get("/{asset_id}", asset.ListOneAsset)
	r.Post("/{asset_id}", asset.CreateAsset)
	r.Put("/{asset_id}", asset.UpdateAsset)
	r.Get("/_query", asset.RichQueryAsset)
	r.Post("/_query", asset.QueryAsset)
	r.Get("/{asset_id}/audit-trail", asset.AuditAsset)
	r.Post("/{asset_id}/validate", asset.ValidateAsset)

	// Link/Unlink APIs
	r.Group(func(r chi.Router) {
		r.Post("/{asset_id}/links", asset.LinkAsset)
		r.Delete("/{asset_id}/links/{link_id}", asset.UnlinkAsset)
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

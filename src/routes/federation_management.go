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
	"chainsource-gateway/controller/federation"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// FederationManagementRouter defines the main routes for the channel API
func FederationManagementRouter() (r chi.Router) {
	r = chi.NewRouter()

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Route("/", fedManagementSubRouting)
	return
}

// fedManagementSubRouting defines the sub routes for the channel APIs
func fedManagementSubRouting(r chi.Router) {
	r.Use(injectSpanMiddleware)
	r.Post("/requests", federation.CreateRequest)
	r.Post("/requests/nodes/update", federation.UpdateNodeDetails)
	r.Get("/requests/nodes/{node_uri}/channels/{channel_id}", federation.FedListOneChannel)
	r.Get("/requests/nodes/{node_uri}/channels/{channel_id}/assets", federation.FedListAssets)
	r.Get("/requests/nodes/{node_uri}/channels/{channel_id}/assets/{asset_id}", federation.FedListOneAsset)
}

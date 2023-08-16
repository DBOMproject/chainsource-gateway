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

// GatewayRouter defines the main routes for the Gateway API
func GatewayRouter() (r chi.Router) {
	r = chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Mount("/nodes/{node_uri}", NodeRouter())
	r.Mount("/nodes/{node_uri}/notaries", NotaryRouter())
	r.Mount("/nodes/{node_uri}/channels", ChannelRouter())
	r.Mount("/nodes/{node_uri}/channels/{channel_id}/assets", AssetRouter())

	r.Get("/federation/requests/all", federation.GetAllRequest)
	r.Get("/federation/requests/{request_id}", federation.GetOneRequest)
	r.Post("/federation/requests/{request_id}/accept", federation.AcceptRequest)
	r.Post("/federation/requests/{request_id}/reject", federation.RejectRequest)
	r.Post("/federation/revoke", federation.RevokeRequest)
	return
}

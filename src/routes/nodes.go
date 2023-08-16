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
	"chainsource-gateway/controller/node"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// NodeRouter defines the main routes for the Node API
func NodeRouter() (r chi.Router) {
	r = chi.NewRouter()

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Route("/", nodeSubRouting)
	return
}

// nodeSubRouting defines the sub routes for the Node APIs
func nodeSubRouting(r chi.Router) {
	r.Use(injectSpanMiddleware)
	// Get node details
	r.Get("/_metadata", node.GetNodeDetails)
}

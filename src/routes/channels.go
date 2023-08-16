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
	"chainsource-gateway/controller/channel"
	"chainsource-gateway/helpers"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

var channelLog = helpers.GetLogger("ChannelRouter")

// ChannelRouter defines the main routes for the channel API
func ChannelRouter() (r chi.Router) {
	r = chi.NewRouter()

	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Route("/", channelSubRouting)
	return
}

// channelSubRouting defines the sub routes for the channel APIs
func channelSubRouting(r chi.Router) {
	r.Use(injectSpanMiddleware)
	// Base CRUD
	r.Get("/", channel.ListChannels)
	r.Post("/", channel.CreateChannel)
	r.Get("/{channel_id}", channel.ListOneChannel)

	// Notary add/remove
	r.Group(func(r chi.Router) {
		r.Put("/{channel_id}/notary", channel.UpdateChannel)
		r.Delete("/{channel_id}/notary/{notary_id}", channel.DeleteChannel)
	})
}

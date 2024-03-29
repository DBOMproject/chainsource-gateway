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

package repo

import (
	"chainsource-gateway/agent"
	"chainsource-gateway/helpers"
	"chainsource-gateway/responses"
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"
)

// ListChannels is a controller function to retrieve channel ids
func ListChannels(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "List Channels")
	defer span.Finish()
	assetVars := r.Context().Value("assetVars").(helpers.AssetRoutingVars)
	requestAgent := r.Context().Value("agent").(agent.Agent)
	log.Debug().Msgf("Getting channels from agent at %s:%d",
		requestAgent.GetHost(), requestAgent.GetPort())

	var result interface{}
	resultStream, err := requestAgent.ListChannels(ctx, agent.QueryArgs{
		ChannelID: assetVars.ChannelID,
	})
	if err != nil {
		if err == helpers.ErrNotFound {
			render.Render(w, r, responses.ErrDoesNotExist(err))
		} else if err == helpers.ErrUnauthorized {
			render.Render(w, r, responses.ErrAgentUnauthorized(err))
		} else {
			render.Render(w, r, responses.ErrAgent(err))
		}
		return
	}
	err = json.NewDecoder(resultStream).Decode(&result)
	render.JSON(w, r, result)
}

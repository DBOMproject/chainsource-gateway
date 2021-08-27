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
	"chainsource-gateway/responses"
	"chainsource-gateway/schema"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"
)

// QueryAsset is a controller function to query assets from a channel on the repository
func QueryAsset(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Query Asset")
	defer span.Finish()
	assetVars := r.Context().Value("assetVars").(helpers.AssetRoutingVars)
	query := r.Context().Value("assetQueryVars")
	requestAgent := r.Context().Value("agent").(agent.Agent)
	log.Debug().Msgf("Querying assets from channel %s from agent at %s:%d", assetVars.ChannelID,
		requestAgent.GetHost(), requestAgent.GetPort())
	assetSchema := r.Context().Value("schemaValidator").(schema.AssetSchema)

	var js = r.Context().Value("JSONBody")

	if js != nil {
		query = MergeJSONMaps(query.(map[string]interface{}), js.(map[string]interface{}))
	}

	errStr, isValid, err := assetSchema.ValidateQueryAsset(ctx, query.(map[string]interface{}))
	if err != nil {
		render.Render(w, r, responses.ErrInternalServer(err))
		return
	}
	if !isValid {
		render.Render(w, r, responses.ErrInvalidRequest(errors.New(errStr)))
		return
	}

	var queryArgs agent.RichQueryArgs
	jsonString, _ := json.Marshal(query)
	err = json.Unmarshal(jsonString, &queryArgs)
	result, err := requestAgent.QueryAssets(ctx, agent.QueryArgs{
		ChannelID: assetVars.ChannelID,
	}, queryArgs)
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
	// err = json.NewDecoder(result).Decode(&result)
	render.JSON(w, r, result)
}

// MergeJSONMaps merges json maps together
func MergeJSONMaps(maps ...map[string]interface{}) (result map[string]interface{}) {
	result = make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			if v != nil {
				result[k] = v
			}
		}
	}
	return result
}

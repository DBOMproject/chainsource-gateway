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
	"chainsource-gateway/tracing"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"
)

// UpdateAsset is a controller function to update an asset on a channel the repository with a given assetID
func UpdateAsset(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Update Asset")
	defer span.Finish()

	assetVars := r.Context().Value("assetVars").(helpers.AssetRoutingVars)
	assetSchema := r.Context().Value("schemaValidator").(schema.AssetSchema)
	requestAgent := r.Context().Value("agent").(agent.Agent)

	log.Info().Msgf("Updating %s/%s to agent at %s:%d", assetVars.ChannelID, assetVars.AssetID,
		requestAgent.GetHost(), requestAgent.GetPort())

	errStr, isValid, err := assetSchema.ValidateAsset(ctx, r.Context().Value("JSONBody").(map[string]interface{}))
	if err != nil {
		render.Render(w, r, responses.ErrInternalServer(err))
		return
	}
	if !isValid {
		render.Render(w, r, responses.ErrInvalidRequest(errors.New(errStr)))
		return
	}

	var requestAsset helpers.Asset
	err = json.NewDecoder(r.Body).Decode(&requestAsset)
	if err != nil {
		tracing.LogAndTraceErr(log, span, err, "Failed to unmarshal, invalid format")
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	// Get current asset state from agent
	result, err := requestAgent.QueryStream(ctx, agent.QueryArgs{
		ChannelID: assetVars.ChannelID,
		AssetID:   assetVars.AssetID,
	})

	if err != nil {
		if err == helpers.ErrUnauthorized{
			render.Render(w, r, responses.ErrAgentUnauthorized(err))
		} else {
			render.Render(w, r, responses.ErrAgent(err))
		}
		return
	}

	// Update asset state with new subassets

	var assetOnAgent helpers.Asset
	err = json.NewDecoder(result).Decode(&assetOnAgent)

	if assetOnAgent.ReadOnly {
		render.Render(w, r, responses.ErrReadOnly())
		return
	}

	requestAsset.AttachedChildren = assetOnAgent.AttachedChildren
	requestAsset.ParentAsset = assetOnAgent.ParentAsset

	// Commit
	res, err := requestAgent.Commit(ctx, agent.CommitArgs{
		ChannelID:  assetVars.ChannelID,
		AssetID:    assetVars.AssetID,
		CommitType: "UPDATE",
		Payload:    requestAsset,
	})

	if err != nil {
		if err == helpers.ErrUnauthorized{
			render.Render(w, r, responses.ErrAgentUnauthorized(err))
		} else {
			render.Render(w, r, responses.ErrAgent(err))
		}
		return
	}

	log.Info().Interface("agentResponse", res).Msg("Updation on agent successful")
	render.Render(w, r, responses.SuccessfulUpdateResponse())
}

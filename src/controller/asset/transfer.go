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
	"time"

	"github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"
)

// TransferAsset is a controller function to manage the a transfer of an asset
// The origin asset must exist while the destination asset must not exist
// Performs the following checks:
//
// 1. If the origin asset exists
//
// 2. If the destination asset does not exist
//
// After the transfer is complete, the origin asset becomes read only
// The origin asset can still be transferred again
func TransferAsset(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Transfer asset")
	defer span.Finish()
	assetVars := r.Context().Value("assetVars").(helpers.AssetRoutingVars)
	assetSchema := r.Context().Value("schemaValidator").(schema.AssetSchema)
	requestAgent := r.Context().Value("agent").(agent.Agent)

	log.Info().Msgf("Transferring asset to %s/%s from agent at %s:%d", assetVars.ChannelID, assetVars.AssetID,
		requestAgent.GetHost(), requestAgent.GetPort())

	errStr, isValid, err := assetSchema.ValidateTransferAsset(ctx, r.Context().Value("JSONBody").(map[string]interface{}))
	if err != nil {
		render.Render(w, r, responses.ErrInternalServer(err))
		return
	}
	if !isValid {
		render.Render(w, r, responses.ErrInvalidRequest(errors.New(errStr)))
		return
	}

	var transferDestinationElement helpers.AssetTransferElement
	err = json.NewDecoder(r.Body).Decode(&transferDestinationElement)
	if err != nil {
		tracing.LogAndTraceErr(log, span, err, "Failed to decode JSON")
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	// Get origin asset state from agent
	log.Debug().Msg("Getting origin asset state")
	childSpan := opentracing.StartSpan("Get current origin asset state", opentracing.ChildOf(span.Context()))
	ctx = opentracing.ContextWithSpan(ctx, childSpan)

	result, err := requestAgent.QueryStream(ctx, agent.QueryArgs{
		ChannelID: assetVars.ChannelID,
		AssetID:   assetVars.AssetID,
	})

	if err != nil {
		if err == helpers.ErrUnauthorized {
			render.Render(w, r, responses.ErrAgentUnauthorized(err))
		} else {
			render.Render(w, r, responses.ErrAgent(err))
		}
		childSpan.Finish()
		return
	}

	childSpan.Finish()

	// Get details of destination asset from an agents
	log.Debug().Msg("Getting current destination asset state")
	childSpan = opentracing.StartSpan("Get destination asset state", opentracing.ChildOf(span.Context()))
	ctx = opentracing.ContextWithSpan(ctx, childSpan)

	destinationAssetAgent, _, err := getChildAssetContextFromAssetElement(ctx, transferDestinationElement.AssetElement)

	if err == nil {
		err = errors.New("Destination asset already exists")
		tracing.LogAndTraceErr(log, childSpan, err, "Invalid transfer operation")
		render.Render(w, r, responses.ErrAlreadyExists(err))
		childSpan.Finish()
		return
	}
	if err != helpers.ErrNotFound {
		if err == helpers.ErrUnauthorized {
			render.Render(w, r, responses.ErrUnauthorizedQueryDestination(err))
		} else {
			render.Render(w, r, responses.ErrFailedQueryDestination(err))
		}
		childSpan.Finish()
		return
	}
	log.Debug().Msgf("Destination asset is %s/%s from agent at %s:%d", transferDestinationElement.ChannelID, transferDestinationElement.AssetID,
		destinationAssetAgent.GetHost(), destinationAssetAgent.GetPort())

	childSpan.Finish()

	// Update origin asset state with new transfer event
	log.Debug().Msg("Committing transfer reference to origin")
	childSpan = opentracing.StartSpan("Committing transfer reference to origin", opentracing.ChildOf(span.Context()))
	ctx = opentracing.ContextWithSpan(ctx, childSpan)

	var requestAsset helpers.Asset
	err = json.NewDecoder(result).Decode(&requestAsset)

	if err != nil {
		tracing.LogAndTraceErr(log, childSpan, err, "Schema invalid (Asset). This should not be possible(!)")
		render.Render(w, r, responses.ErrAgent(err))
		childSpan.Finish()
		return
	}
	if requestAsset.CustodyTransferEvents == nil {
		var attachedArray []helpers.CustodyTransferEvent
		requestAsset.CustodyTransferEvents = attachedArray
	}

	transferEvent := helpers.CustodyTransferEvent{
		Timestamp:            time.Now().UTC().Format("2006-01-02T15:04:05.999Z"),
		TransferDescription:  transferDestinationElement.TransferDescription,
		SourceRepoID:         assetVars.RepoID,
		SourceChannelID:      assetVars.ChannelID,
		SourceAssetID:        assetVars.AssetID,
		DestinationRepoID:    transferDestinationElement.AssetElement.RepoID,
		DestinationChannelID: transferDestinationElement.AssetElement.ChannelID,
		DestinationAssetID:   transferDestinationElement.AssetElement.AssetID,
	}

	requestAsset.CustodyTransferEvents = append(requestAsset.CustodyTransferEvents, transferEvent)
	destinationRequestAsset := requestAsset
	//Make origin asset read only
	requestAsset.ReadOnly = true
	//Make destination asset writeable
	destinationRequestAsset.ReadOnly = false

	_, err = requestAgent.Commit(ctx, agent.CommitArgs{
		ChannelID:  assetVars.ChannelID,
		AssetID:    assetVars.AssetID,
		CommitType: "TRANSFER-OUT",
		Payload:    requestAsset,
	})

	if err != nil {
		if err == helpers.ErrUnauthorized {
			render.Render(w, r, responses.ErrUnauthorizedModifyOrigin(err))
		} else {
			render.Render(w, r, responses.ErrFailedModifyOrigin(err))
		}
		childSpan.Finish()
		return
	}

	childSpan.Finish()

	// Commit updated destination asset with transfer event
	log.Debug().Msg("Committing destination asset")
	childSpan = opentracing.StartSpan("Committing destination asset", opentracing.ChildOf(span.Context()))
	ctx = opentracing.ContextWithSpan(ctx, childSpan)

	_, err = destinationAssetAgent.Commit(ctx, agent.CommitArgs{
		ChannelID:  transferDestinationElement.ChannelID,
		AssetID:    transferDestinationElement.AssetID,
		CommitType: "TRANSFER-IN",
		Payload:    destinationRequestAsset,
	})

	if err != nil {
		if err == helpers.ErrUnauthorized {
			render.Render(w, r, responses.ErrUnauthorizedModifyDestination(err))
		} else {
			render.Render(w, r, responses.ErrFailedModifyDestination(err))
		}
		childSpan.Finish()
		return
	}

	childSpan.Finish()

	render.Render(w, r, responses.SuccessfulTransferResponse())

}

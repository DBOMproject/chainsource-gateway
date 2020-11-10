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
	"github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"
	"net/http"
)

// DetachSubasset is a controller function to manage the detachment of an asset from another asset
// Both assets must exist and can be on different repositories/channels
// Performs the following checks:
//
// 1. If the parent and child asset exists
//
// 2. If any assets are attached to the parent
//
// 2. If (2) is true, if the child asset is  attached to the parent asset
//
// It updates the parent asset with the removal of the child and blanks out the child asset's parent reference
func DetachSubasset(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Detach Asset")
	defer span.Finish()
	assetVars := r.Context().Value("assetVars").(helpers.AssetRoutingVars)
	assetSchema := r.Context().Value("schemaValidator").(schema.AssetSchema)
	requestAgent := r.Context().Value("agent").(agent.Agent)

	log.Info().Msgf("Detaching subasset from %s/%s from agent at %s:%d", assetVars.ChannelID, assetVars.AssetID,
		requestAgent.GetHost(), requestAgent.GetPort())

	errStr, isValid, err := assetSchema.ValidateDetachSubasset(ctx, r.Context().Value("JSONBody").(map[string]interface{}))
	if err != nil {
		render.Render(w, r, responses.ErrInternalServer(err))
		return
	}
	if !isValid {
		render.Render(w, r, responses.ErrInvalidRequest(errors.New(errStr)))
		return
	}

	var childAssetLinkElement helpers.AssetLinkElement
	err = json.NewDecoder(r.Body).Decode(&childAssetLinkElement)
	if err != nil {
		tracing.LogAndTraceErr(log, span, err, "Failed to decode JSON")
		render.Render(w, r, responses.ErrInvalidRequest(err))
		return
	}

	// Get parent asset state from agent
	log.Debug().Msg("Getting current parent state")
	childSpan := opentracing.StartSpan("Get current parent state", opentracing.ChildOf(span.Context()))
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

	// Get details of child asset from an agent
	log.Debug().Msg("Getting current child state")
	childSpan = opentracing.StartSpan("Get current child state", opentracing.ChildOf(span.Context()))
	ctx = opentracing.ContextWithSpan(ctx, childSpan)

	childAssetAgent, childAsset, err := getChildAssetContextFromAssetElement(ctx, childAssetLinkElement.AssetElement)
	if err != nil {
		if err == helpers.ErrUnauthorized {
			render.Render(w, r, responses.ErrUnauthorizedQueryChild(err))
		} else {
			render.Render(w, r, responses.ErrFailedQueryChild(err))
		}
		childSpan.Finish()
		return
	}
	log.Debug().Msgf("Subasset is %s/%s from agent at %s:%d", childAssetLinkElement.ChannelID, childAssetLinkElement.AssetID,
		childAssetAgent.GetHost(), childAssetAgent.GetPort())
	log.Debug().Msg("Child asset does exist and is retrievable")

	childSpan.Finish()

	// Update parent asset state with new subasset references
	log.Debug().Msg("Committing remove child-reference from parent")
	childSpan = opentracing.StartSpan("Commit remove child-reference from parent", opentracing.ChildOf(span.Context()))
	ctx = opentracing.ContextWithSpan(ctx, childSpan)

	var requestAsset helpers.Asset
	err = json.NewDecoder(result).Decode(&requestAsset)

	if err != nil {
		tracing.LogAndTraceErr(log, childSpan, err, "Schema invalid (Asset). This should not be possible(!)")
		render.Render(w, r, responses.ErrAgent(err))
		childSpan.Finish()
		return
	}
	if requestAsset.AttachedChildren == nil {
		err = errors.New("no subassets")
		tracing.LogAndTraceErr(log, childSpan, err, "Invalid detach operation")
		render.Render(w, r, responses.ErrNotAttached(err))
		childSpan.Finish()
		return
	}

	var subassetIsLinked bool

	for i, curChild := range requestAsset.AttachedChildren {
		if curChild.AssetID == childAssetLinkElement.AssetID {
			subassetIsLinked = true
			requestAsset.AttachedChildren =
				append(requestAsset.AttachedChildren[:i],
					requestAsset.AttachedChildren[i+1:]...)
			break
		}
	}

	if !subassetIsLinked {
		err = errors.New("Subasset is not linked to asset")
		tracing.LogAndTraceErr(log, childSpan, err, "Invalid Detach")
		render.Render(w, r, responses.ErrNotAttached(err))
		childSpan.Finish()
		return
	}

	_, err = requestAgent.Commit(ctx, agent.CommitArgs{
		ChannelID:  assetVars.ChannelID,
		AssetID:    assetVars.AssetID,
		CommitType: "DETACH",
		Payload:    requestAsset,
	})

	if err != nil {
		if err == helpers.ErrUnauthorized {
			render.Render(w, r, responses.ErrUnauthorizedModifyParent(err))
		} else {
			render.Render(w, r, responses.ErrFailedModifyParent(err))
		}
		childSpan.Finish()
		return
	}

	childSpan.Finish()

	// Commit modified child asset to channel (Remove Parent)
	log.Debug().Msg("Committing remove parent-reference from child")
	childSpan = opentracing.StartSpan("Commit remove parent-reference from child", opentracing.ChildOf(span.Context()))
	ctx = opentracing.ContextWithSpan(ctx, childSpan)

	childAsset.ParentAsset = &helpers.AssetLinkElement{}
	_, err = childAssetAgent.Commit(ctx, agent.CommitArgs{
		ChannelID:  childAssetLinkElement.ChannelID,
		AssetID:    childAssetLinkElement.AssetID,
		CommitType: "DETACH",
		Payload:    childAsset,
	})

	if err != nil {
		if err == helpers.ErrUnauthorized {
			render.Render(w, r, responses.ErrUnauthorizedModifyChild(err))
		} else {
			render.Render(w, r, responses.ErrFailedModifyChild(err))
		}
		childSpan.Finish()
		return
	}
	childSpan.Finish()

	render.Render(w, r, responses.SuccessfulDetachResponse())

}

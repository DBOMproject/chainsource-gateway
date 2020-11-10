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

// AttachSubasset is a controller function to manage the attachment of an asset to another asset
// Both assets must exist and can be on different repositories/channels
// Performs the following checks:
//
// 1. If the parent and child asset exists
//
// 2. If the asset that is being attached is the same as the asset it is being attached to
//
// 3. If the asset is already attached to the asset
//
// It updates the parent asset with the new child and the child asset with a new parent.
// If the child already has a parent, it is replaced.
func AttachSubasset(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Attach subasset")
	defer span.Finish()
	assetVars := r.Context().Value("assetVars").(helpers.AssetRoutingVars)
	requestAgent := r.Context().Value("agent").(agent.Agent)
	assetSchema := r.Context().Value("schemaValidator").(schema.AssetSchema)

	log.Info().Msgf("Attaching a subasset to %s/%s from agent at %s:%d", assetVars.ChannelID, assetVars.AssetID,
		requestAgent.GetHost(), requestAgent.GetPort())

	errStr, isValid, err := assetSchema.ValidateAttachSubasset(ctx, r.Context().Value("JSONBody").(map[string]interface{}))
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

	if assetVars.AssetID == childAssetLinkElement.AssetID {
		err = errors.New("You cannot attach an asset to itself")
		tracing.LogAndTraceErr(log, span, err, "Invalid attach operation")
		render.Render(w, r, responses.ErrConflict(err))
		span.Finish()
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
	log.Info().Msg("Child asset does exist and is retrievable")

	childSpan.Finish()

	// Update parent asset state with new subasset references
	log.Debug().Msg("Committing add child-reference to parent")
	childSpan = opentracing.StartSpan("Commit add child-reference to parent", opentracing.ChildOf(span.Context()))
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
		var attachedArray []helpers.AssetLinkElement
		requestAsset.AttachedChildren = attachedArray
	}

	for _, curChild := range requestAsset.AttachedChildren {
		if curChild.AssetID == childAssetLinkElement.AssetID {
			err = errors.New("this child is already part of the asset")
			tracing.LogAndTraceErr(log, childSpan, err, "Duplicate Subasset")
			render.Render(w, r, responses.ErrAlreadyAttached(err))
			childSpan.Finish()
			return
		}
	}
	requestAsset.AttachedChildren = append(requestAsset.AttachedChildren, childAssetLinkElement)

	_, err = requestAgent.Commit(ctx, agent.CommitArgs{
		ChannelID:  assetVars.ChannelID,
		AssetID:    assetVars.AssetID,
		CommitType: "ATTACH",
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

	// Commit modified child asset to channel
	log.Debug().Msg("Committing add parent-reference to child")
	childSpan = opentracing.StartSpan("Commit add parent-reference to child", opentracing.ChildOf(span.Context()))
	ctx = opentracing.ContextWithSpan(ctx, childSpan)

	childAsset.ParentAsset = &helpers.AssetLinkElement{
		AssetElement: helpers.AssetElement{
			RepoID:    assetVars.RepoID,
			ChannelID: assetVars.ChannelID,
			AssetID:   assetVars.AssetID,
		},
		Role:    childAssetLinkElement.Role,
		SubRole: childAssetLinkElement.SubRole,
	}

	_, err = childAssetAgent.Commit(ctx, agent.CommitArgs{
		ChannelID:  childAssetLinkElement.ChannelID,
		AssetID:    childAssetLinkElement.AssetID,
		CommitType: "ATTACH",
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

	render.Render(w, r, responses.SuccessfulAttachResponse())
}

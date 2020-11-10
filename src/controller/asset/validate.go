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
	"chainsource-gateway/pgp"
	"chainsource-gateway/responses"
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"
)

var logger = helpers.GetLogger("validateAsset")
var fingerprintKey = "manufactureFingerprint"
var signatureKey = "manufactureSignature"

// validateAsset is a controller function to validate the signature of an asset
func ValidateAsset(w http.ResponseWriter, r *http.Request) {
	logger.Info().Msgf("Validate Asset")
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Validate asset")
	defer span.Finish()

	assetVars := r.Context().Value("assetVars").(helpers.AssetRoutingVars)
	requestAgent := r.Context().Value("agent").(agent.Agent)
	signingService := r.Context().Value("signatureValidator").(pgp.SignatureValidator)
	log.Debug().Msgf("Getting %s/%s from agent at %s:%d", assetVars.ChannelID, assetVars.AssetID,
		requestAgent.GetHost(), requestAgent.GetPort())

	var result helpers.Asset
	resultStream, err := requestAgent.QueryStream(ctx, agent.QueryArgs{
		ChannelID: assetVars.ChannelID,
		AssetID:   assetVars.AssetID,
	})
	if err != nil {
		if err == helpers.ErrNotFound {
			render.Render(w, r, responses.ErrDoesNotExist(err))
		} else if err == helpers.ErrUnauthorized{
			render.Render(w, r, responses.ErrAgentUnauthorized(err))
		} else {
			render.Render(w, r, responses.ErrAgent(err))
		}
		return
	}
	err = json.NewDecoder(resultStream).Decode(&result)
	signature := result.ManufactureSignature
	if len(signature) == 0 {
		render.Render(w, r, responses.ErrNoSignature())
		return
	}
	meta := result.AssetMetadata.(map[string]interface{})
	if meta == nil || len(meta) < 1 {
		render.Render(w, r, responses.ErrNoSignature())
		return
	}
	fingerprint := ""
	if _, ok := meta[fingerprintKey]; ok {
		fingerprint = meta[fingerprintKey].(string)
	}
	if len(fingerprint) == 0 {
		render.Render(w, r, responses.ErrNoFingerprint())
		return
	}
	delete(meta, fingerprintKey)
	result.AssetMetadata = meta
	input := helpers.AssetNoChildParent(result)

	res, err := signingService.Validate(ctx, pgp.ValidateArgs{
		Fingerprint: fingerprint,
		Signature:   signature,
		Input:       input,
	})

	if err != nil {
		render.Render(w, r, responses.ErrInvalidSignature())
		return
	}

	render.JSON(w, r, res)
}

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
	"net/http"

	"github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"
)

// AuditAsset is a controller function to perform an audit of an assetID
// Audits are managed by the agent since it's implementation specific
// No validity checks are performed

func AuditAsset(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Audit Asset")
	defer span.Finish()

	assetVars := r.Context().Value("assetVars").(helpers.AssetRoutingVars)
	requestAgent := r.Context().Value("agent").(agent.Agent)

	log.Debug().Msgf("Getting audit trail of %s/%s from agent at %s:%d", assetVars.ChannelID, assetVars.AssetID,
		requestAgent.GetHost(), requestAgent.GetPort())

	auditTrail, err := requestAgent.QueryAuditTrail(ctx, agent.QueryArgs{
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
	render.JSON(w, r, auditTrail)
}

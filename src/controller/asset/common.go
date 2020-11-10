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

// Package asset contains all the controller functions for managing an asset
package asset

import (
	"chainsource-gateway/agent"
	"chainsource-gateway/helpers"
	"chainsource-gateway/tracing"
	"context"
	"encoding/json"
	"github.com/opentracing/opentracing-go"
)

var log = helpers.GetLogger("AssetController")

// Gets the child asset and the agent for the child asset, given a link element
func getChildAssetContextFromAssetElement(ctx context.Context, element helpers.AssetElement) (childAssetAgent agent.Agent, childAsset helpers.Asset, err error) {
	agentProvider := ctx.Value("agentProvider").(agent.Provider)
	childAssetAgentConfig, err := agentProvider.GetAgentConfigForRepo(element.RepoID)
	if err != nil {
		tracing.LogAndTraceErr(log, opentracing.SpanFromContext(ctx), err, "HttpAgent for child asset is not known")
		log.Err(err).Msg("HttpAgent for child asset is not known")
		return
	}
	childAssetAgent = agentProvider.NewAgent(&childAssetAgentConfig)

	childAssetStream, err := childAssetAgent.QueryStream(ctx, agent.QueryArgs{
		ChannelID: element.ChannelID,
		AssetID:   element.AssetID,
	})
	if err != nil {
		return
	}
	err = json.NewDecoder(childAssetStream).Decode(&childAsset)
	if err != nil {
		log.Err(err).Msg("Error Decoding Child Asset (!) Should not be possible ")
		return
	}
	return
}

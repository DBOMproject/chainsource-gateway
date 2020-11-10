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
	"chainsource-gateway/schema"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
)


// openTestJSON returns an opened raw test JSON file
func openTestJSON(path string) *os.File {
	jsonFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return jsonFile
}

// decodeAsset unmarshals an asset from a stream
func decodeAsset(closer io.ReadCloser) helpers.Asset {
	var result helpers.Asset
	json.NewDecoder(closer).Decode(&result)
	return result
}

// invalidJSON is used as an example of unparseable json
const invalidJSON = "in,VAl.ID"

// injectMockAssetContext injects general asset context into a request
func injectMockAssetContext(r *http.Request, repoID string, channelID string, assetID string, mockAgent agent.Agent, mockValidator schema.AssetSchema) (reqWithContext *http.Request) {
	vars := helpers.AssetRoutingVars{
		RepoID:    repoID,
		ChannelID: channelID,
		AssetID:   assetID,
	}
	ctx := context.WithValue(r.Context(), "JSONBody", map[string]interface{}{})
	ctx = context.WithValue(ctx, "assetVars", vars)
	ctx = context.WithValue(ctx, "agent", mockAgent)
	ctx = context.WithValue(ctx, "schemaValidator", mockValidator)
	reqWithContext = r.WithContext(ctx)
	return
}

// getAgentSuccessResponse returns a general success response expected from the agent
func getAgentSuccessResponse() map[string]interface{} {
	return map[string]interface{}{"success": "true"}
}

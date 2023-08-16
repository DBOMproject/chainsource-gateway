/*
 * Copyright 2023 Unisys Corporation
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
	"chainsource-gateway/controller/node"
	"chainsource-gateway/helpers"
	"chainsource-gateway/responses"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nats-io/nats.go"
)

var logger = helpers.GetLogger(helpers.AssetLogger)

// Validate is a helper function to validate the request if it is authorized
func Validate(w http.ResponseWriter, r *http.Request) bool {
	nodeIdFromUri := strings.Split(chi.URLParam(r, "node_uri"), ".")[0]
	channelFromUri := chi.URLParam(r, "channel_id")
	nodeMetaData := node.GetNodeDetailsFromId(helpers.GetNodeID())
	logger.Info().Msgf("Node _metadata from Asset: %v\n", nodeMetaData)

	for _, nodeConn := range nodeMetaData.NodeConnections {
		if nodeConn.NodeId == nodeIdFromUri && nodeConn.Status == "FEDERATION_SUCCESS" {
			for _, channelConn := range nodeConn.ChannelConnections {
				if channelConn.ChannelId == channelFromUri && channelConn.Status == "CONNECTED" {
					return true
				}
			}
		}
	}
	return false
}

// ListAssets is a controller function to list all the asset
func ListAssets(w http.ResponseWriter, r *http.Request) {
	var nodeUri = chi.URLParam(r, "node_uri")
	var nodeIdFromRequest = strings.Split(nodeUri, ".")[0]
	var channelIdFromRequest = chi.URLParam(r, "channel_id")
	logger.Info().Msgf("Received request to list assets for node %s and channel %s", nodeIdFromRequest, channelIdFromRequest)

	if nodeIdFromRequest == helpers.LocalNodeId || nodeIdFromRequest == helpers.GetNodeID() {
		nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
		if ncErr != nil {
			logger.Err(ncErr).Msgf("NATS connection error")
			helpers.HandleError(w, r, "NATS connection error")
			return
		}
		defer nc.Close()

		request := helpers.AssetRoutingVars{
			ChannelID: channelIdFromRequest,
		}

		requestMarshal, marshalErr := json.Marshal(request)
		if marshalErr != nil {
			logger.Err(marshalErr).Msgf("Request marshal error")
			helpers.HandleError(w, r, "Request marshal error")
			return
		}

		// Send the request
		msg, msgErr := nc.Request("asset.all", requestMarshal, helpers.TimeOut*time.Second)
		if msgErr != nil {
			logger.Err(msgErr).Msgf("NATS request error")
			helpers.HandleError(w, r, "NATS request error")
			return
		}

		logger.Info().Msgf("Received response: %s\n", msg.Data)

		// Use the response
		var response []helpers.AssetMeta
		unmarshalErr := json.Unmarshal([]byte(string(msg.Data)), &response)
		if unmarshalErr != nil {
			logger.Err(unmarshalErr).Msgf(helpers.UnmarshalErr)
			helpers.HandleError(w, r, "Unmarshal error")
			return
		}
		render.JSON(w, r, response)
	} else {
		logger.Info().Msgf("Request is not for this node, forwarding to %s", nodeUri)
		if Validate(w, r) {
			var host = nodeUri + ":7205"
			res, err := helpers.GetRequest("https://" + host + "/api/v2/federation/requests/nodes/" + nodeUri + "/channels/" + channelIdFromRequest + "/assets")
			if err != nil {
				logger.Err(err).Msgf(helpers.ResponseErr)
				helpers.HandleError(w, r, "Request error")
				return
			}
			data, err := io.ReadAll(res)
			if err != nil {
				http.Error(w, helpers.ReadingErr, http.StatusInternalServerError)
				helpers.HandleError(w, r, "Reading error")
				return
			}
			var response []helpers.AssetMeta
			unmarshalErr := json.Unmarshal([]byte(string(data)), &response)
			if unmarshalErr != nil {
				logger.Err(unmarshalErr).Msgf(helpers.UnmarshalErr)
				helpers.HandleError(w, r, "Unmarshal error")
				return
			}
			render.JSON(w, r, response)
		} else {
			render.Render(w, r, responses.ErrUnauthorizedQueryDestination(errors.New(helpers.Unauthorized)))
		}
	}
}

// ListOneAsset is a controller function to lists asset by ID
func ListOneAsset(w http.ResponseWriter, r *http.Request) {
	var nodeUri = chi.URLParam(r, "node_uri")
	var nodeIdFromRequest = strings.Split(nodeUri, ".")[0]
	var channelIdFromRequest = chi.URLParam(r, "channel_id")
	var assetIdFromRequest = chi.URLParam(r, "asset_id")

	logger.Info().Msgf("Received request to list asset %s for node %s and channel %s", assetIdFromRequest, nodeIdFromRequest, channelIdFromRequest)

	if nodeIdFromRequest == helpers.LocalNodeId || nodeIdFromRequest == helpers.GetNodeID() {
		nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
		if ncErr != nil {
			logger.Err(ncErr).Msgf(helpers.NatsConnectError)
			helpers.HandleError(w, r, "NATS connection error")
			return
		}
		defer nc.Close()

		request := helpers.AssetRoutingVars{
			ChannelID: channelIdFromRequest,
			AssetID:   assetIdFromRequest,
		}

		requestMarshal, marshalErr := json.Marshal(request)
		if marshalErr != nil {
			logger.Err(marshalErr).Msgf(helpers.MarshalErr)
			helpers.HandleError(w, r, "Marshal error")
			return
		}

		// Send the request
		msg, msgErr := nc.Request("asset.one", requestMarshal, helpers.TimeOut*time.Second)
		if msgErr != nil {
			logger.Err(msgErr).Msgf(helpers.MsgErr)
			helpers.HandleError(w, r, "NATS request error")
			return
		}

		logger.Info().Msgf(helpers.Success+"%s\n", msg.Data)

		// Use the response
		var response []helpers.AssetMeta
		unmarshalErr := json.Unmarshal([]byte(string(msg.Data)), &response)
		if unmarshalErr != nil {
			logger.Err(unmarshalErr).Msgf(helpers.UnmarshalErr)
			helpers.HandleError(w, r, "Unmarshal error")
			return
		}
		render.JSON(w, r, response)

	} else {
		if Validate(w, r) {
			var host = nodeUri + ":7205"
			res, err := helpers.GetRequest("https://" + host + "/api/v2/federation/requests/nodes/" + nodeUri + "/channels/" + channelIdFromRequest + "/assets/" + assetIdFromRequest)
			if err != nil {
				logger.Err(err).Msgf(helpers.ResponseErr)
				helpers.HandleError(w, r, "Request error")
				return
			}
			data, err := io.ReadAll(res)
			if err != nil {
				http.Error(w, helpers.ReadingErr, http.StatusInternalServerError)
				helpers.HandleError(w, r, "Reading error")
				return
			}
			var response []helpers.AssetMeta
			unmarshalErr := json.Unmarshal([]byte(string(data)), &response)
			if unmarshalErr != nil {
				logger.Err(unmarshalErr).Msgf(helpers.UnmarshalErr)
				helpers.HandleError(w, r, "Unmarshal error")
				return
			}
			render.JSON(w, r, response)
		} else {
			render.Render(w, r, responses.ErrUnauthorizedQueryDestination(errors.New(helpers.Unauthorized)))
		}
	}
}

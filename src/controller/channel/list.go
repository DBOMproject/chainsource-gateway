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

package channel

import (
	"chainsource-gateway/controller/federation"
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

var logger = helpers.GetLogger(helpers.ChannelLogger)

// Validate is a controller function to validate the request
func Validate(w http.ResponseWriter, r *http.Request) bool {
	nodeIdFromUri := strings.Split(chi.URLParam(r, "node_uri"), ".")[0]
	channelFromUri := chi.URLParam(r, "channel_id")
	nodeMetaData := node.GetNodeDetailsFromId(helpers.GetNodeID())
	logger.Info().Msgf("Node _metadata from Channel: %v\n", nodeMetaData)

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

// ListChannels is a controller function to list all the channels
func ListChannels(w http.ResponseWriter, r *http.Request) {
	var nodeUri = chi.URLParam(r, "node_uri")
	var nodeIdFromRequest = strings.Split(nodeUri, ".")[0]

	logger.Info().Msgf("Received request to list all channels from node: %v\n", nodeIdFromRequest)

	if nodeIdFromRequest == helpers.LocalNodeId || nodeIdFromRequest == helpers.GetNodeID() {
		nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
		if ncErr != nil {
			logger.Err(ncErr).Msgf(helpers.NatsConnectError)
			helpers.HandleError(w, r, helpers.NatsConnectError)
			return
		}
		defer nc.Close()

		requestMarshal, marshalErr := json.Marshal(nodeIdFromRequest)
		if marshalErr != nil {
			logger.Err(marshalErr).Msgf(helpers.MarshalErr)
			helpers.HandleError(w, r, helpers.MarshalErr)
			return
		}

		// Send the request
		msg, msgErr := nc.Request("channel.all", requestMarshal, helpers.TimeOut*time.Second)
		if msgErr != nil {
			logger.Err(msgErr).Msgf(helpers.MsgErr)
			helpers.HandleError(w, r, helpers.MsgErr)
			return
		}

		logger.Info().Msgf(helpers.Success+"%s\n", msg.Data)

		// Use the response
		var response []helpers.Channel
		unmarshalErr := json.Unmarshal([]byte(string(msg.Data)), &response)
		if unmarshalErr != nil {
			logger.Err(unmarshalErr).Msgf(helpers.UnmarshalErr)
			helpers.HandleError(w, r, helpers.UnmarshalErr)
			return
		}

		render.JSON(w, r, response)
	} else {
		render.Render(w, r, responses.ErrDoesNotExist(errors.New(helpers.InvalidRequest)))
	}
}

// ListOneChannel is a controller function to lists channel by ID
func ListOneChannel(w http.ResponseWriter, r *http.Request) {
	var nodeUri = chi.URLParam(r, "node_uri")
	var nodeIdFromRequest = strings.Split(nodeUri, ".")[0]
	var channelIdFromRequest = chi.URLParam(r, "channel_id")

	logger.Info().Msgf("Received request to list channel: %v from node: %v\n", channelIdFromRequest, nodeIdFromRequest)

	if nodeIdFromRequest == helpers.LocalNodeId || nodeIdFromRequest == helpers.GetNodeID() {
		nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
		if ncErr != nil {
			logger.Err(ncErr).Msgf(helpers.NatsConnectError)
			helpers.HandleError(w, r, helpers.NatsConnectError)
			return
		}
		defer nc.Close()

		request := helpers.ListOneChannelMeta{
			ChannelID: channelIdFromRequest,
		}

		requestMarshal, marshalErr := json.Marshal(request)
		if marshalErr != nil {
			logger.Err(marshalErr).Msgf(helpers.MarshalErr)
			helpers.HandleError(w, r, helpers.MarshalErr)
			return
		}

		// Send the request
		msg, msgErr := nc.Request("channel.one", requestMarshal, helpers.TimeOut*time.Second)
		if msgErr != nil {
			logger.Err(msgErr).Msgf(helpers.MsgErr)
			helpers.HandleError(w, r, helpers.MsgErr)
			return
		}

		logger.Info().Msgf(helpers.Success+"%s\n", msg.Data)

		// Use the response
		var response []helpers.Channel
		unmarshalErr := json.Unmarshal([]byte(string(msg.Data)), &response)
		if unmarshalErr != nil {
			logger.Err(unmarshalErr).Msgf(helpers.UnmarshalErr)
			helpers.HandleError(w, r, helpers.UnmarshalErr)
			return
		}

		render.JSON(w, r, response)
	} else {
		var host = nodeUri + ":7205"
		if Validate(w, r) {
			res, err := helpers.GetRequest("https://" + host + "/api/v2/federation/requests/nodes/" + nodeUri + "/channels/" + channelIdFromRequest)
			if err != nil {
				logger.Err(err).Msgf(helpers.ResponseErr)
				helpers.HandleError(w, r, helpers.ResponseErr)
				return
			}
			data, err := io.ReadAll(res)
			if err != nil {
				http.Error(w, helpers.ReadingErr, http.StatusInternalServerError)
				return
			}
			var response []helpers.Channel
			unmarshalErr := json.Unmarshal([]byte(string(data)), &response)
			if unmarshalErr != nil {
				logger.Err(unmarshalErr).Msgf(helpers.UnmarshalErr)
				helpers.HandleError(w, r, helpers.UnmarshalErr)
				return
			}
			render.JSON(w, r, response)
		} else {
			body := helpers.FederationRequestOperations{
				NodeURI:   helpers.GetNodeURI(),
				NodeID:    helpers.GetNodeID(),
				ChannelID: channelIdFromRequest,
				Type:      "INITIATE_REQUEST",
			}
			resToSend, marshalErr := json.Marshal(body)
			if marshalErr != nil {
				logger.Err(marshalErr).Msg(helpers.MarshalErr)
			}
			res, err := helpers.PostJSONRequest("https://"+host+"/api/v2/federation/requests", []byte(resToSend))
			if err != nil {
				logger.Err(err).Msgf(helpers.ResponseErr)
				helpers.HandleError(w, r, helpers.ResponseErr)
				return
			}

			if res.StatusCode == http.StatusOK {
				// Update _metadata locally
				nodeMetaBody := helpers.FederationRequestOperations{
					NodeID:    nodeIdFromRequest,
					ChannelID: channelIdFromRequest,
					Type:      "FEDERATION_SUCCESS",
				}

				federation.SelfUpdateNodeDetails(nodeMetaBody)

				data, err := io.ReadAll(res.Body)
				if err != nil {
					http.Error(w, "Error reading data", http.StatusInternalServerError)
					return
				}
				parsedData, err := helpers.ParseJSONData(data)
				if err != nil {
					http.Error(w, "Error parsing JSON", http.StatusInternalServerError)
					return
				}
				render.JSON(w, r, parsedData)
			} else {
				render.Render(w, r, responses.ErrDoesNotExist(errors.New(helpers.InvalidRequest)))
			}
		}
	}
}

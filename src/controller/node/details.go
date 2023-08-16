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

package node

import (
	"chainsource-gateway/helpers"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nats-io/nats.go"
)

var logger = helpers.GetLogger(helpers.NodeLogger)

// GetNodeDetails is a controller function to get node details
func GetNodeDetails(w http.ResponseWriter, r *http.Request) {
	var nodeUri = chi.URLParam(r, "node_uri")
	var nodeIdFromRequest = strings.Split(nodeUri, ".")[0]

	logger.Info().Msgf("Received request to get node details for node id: %s", nodeIdFromRequest)

	if nodeIdFromRequest == helpers.LocalNodeId || nodeIdFromRequest == helpers.GetNodeID() {
		nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
		if ncErr != nil {
			logger.Err(ncErr).Msg(helpers.NatsConnectError)
			return
		}
		defer nc.Close()

		requestMarshal, marshalErr := json.Marshal(nodeIdFromRequest)
		if marshalErr != nil {
			logger.Err(marshalErr).Msg(helpers.MarshalErr)
			return
		}

		// Send the request
		msg, msgErr := nc.Request("node.details", requestMarshal, helpers.TimeOut*time.Second)
		if msgErr != nil {
			logger.Err(msgErr).Msgf(helpers.MsgErr)
			return
		}

		logger.Info().Msgf(helpers.Success + "\n")

		// Use the response
		var response []helpers.Node
		unmarshalErr := json.Unmarshal([]byte(string(msg.Data)), &response)
		if unmarshalErr != nil {
			logger.Err(unmarshalErr).Msgf(helpers.UnmarshalErr)
			return
		}

		render.JSON(w, r, response)
	}
}

// GetNodeDetailsFromId is a function to get node details from node id
func GetNodeDetailsFromId(nodeId string) helpers.Node {
	nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
	if ncErr != nil {
		logger.Err(ncErr).Msgf(helpers.NatsConnectError)
		return helpers.Node{}
	}
	defer nc.Close()

	requestMarshal, marshalErr := json.Marshal(nodeId)
	if marshalErr != nil {
		logger.Err(marshalErr).Msgf(helpers.MarshalErr)
		return helpers.Node{}
	}

	// Send the request
	msg, msgErr := nc.Request("node.details", requestMarshal, helpers.TimeOut*time.Second)
	if msgErr != nil {
		logger.Err(msgErr).Msgf(helpers.MsgErr)
		return helpers.Node{}
	}

	var response []helpers.Node
	unmarshalErr := json.Unmarshal([]byte(string(msg.Data)), &response)
	if unmarshalErr != nil {
		logger.Err(unmarshalErr).Msgf(helpers.UnmarshalErr)
		return helpers.Node{}
	}

	return response[0]
}

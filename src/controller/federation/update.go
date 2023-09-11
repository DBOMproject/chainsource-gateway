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

package federation

import (
	"chainsource-gateway/helpers"
	"chainsource-gateway/responses"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/render"
	"github.com/nats-io/nats.go"
)

// UpdateNodeDetails is a controller function that updates the node details
func UpdateNodeDetails(w http.ResponseWriter, r *http.Request) {
	logger.Info().Msgf("Received request to update node details")

	nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
	if ncErr != nil {
		logger.Err(ncErr).Msgf(helpers.NatsConnectError)
		helpers.HandleError(w, r, helpers.NatsConnectError)
		return
	}
	defer nc.Close()

	var nodeUpdateRequest helpers.FederationRequestOperations
	json.NewDecoder(r.Body).Decode(&nodeUpdateRequest)

	request := helpers.FederationRequestOperations{
		NodeID:    nodeUpdateRequest.NodeID,
		ChannelID: nodeUpdateRequest.ChannelID,
		Type:      nodeUpdateRequest.Type,
	}

	requestMarshal, marshalErr := json.Marshal(request)
	if marshalErr != nil {
		logger.Err(marshalErr).Msgf(helpers.MarshalErr)
		helpers.HandleError(w, r, helpers.MarshalErr)
		return
	}

	// Send the request
	msg, msgErr := nc.Request("node.update", requestMarshal, helpers.TimeOut*time.Second)
	if msgErr != nil {
		logger.Err(msgErr).Msgf(helpers.MsgErr)
		return
	}

	logger.Info().Msgf(helpers.Response + "\n")

	// Use the response
	var response helpers.NodeResultResponse
	unmarshalErr := json.Unmarshal([]byte(string(msg.Data)), &response)
	if unmarshalErr != nil {
		logger.Err(unmarshalErr).Msgf(helpers.UnmarshalErr)
		return
	}

	if response.Success {
		render.JSON(w, r, response)
	} else {
		render.Render(w, r, responses.ErrCustom(errors.New(response.Status)))
	}
}

// SelfUpdateNodeDetails is a function that updates the node details locally and sends the update to NATS
func SelfUpdateNodeDetails(body helpers.FederationRequestOperations) {
	logger.Info().Msgf("Received request to self update node details")

	nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
	if ncErr != nil {
		logger.Err(ncErr).Msgf(helpers.NatsConnectError)
		return
	}
	defer nc.Close()

	requestMarshal, marshalErr := json.Marshal(body)
	if marshalErr != nil {
		logger.Err(marshalErr).Msgf(helpers.MarshalErr)
		return
	}

	// Send the request
	_, msgErr := nc.Request("node.update", requestMarshal, helpers.TimeOut*time.Second)
	if msgErr != nil {
		logger.Err(msgErr).Msgf(helpers.MsgErr)
		return
	}
}

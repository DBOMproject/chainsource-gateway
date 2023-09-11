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

// CreateRequest is a controller function to create federation request
func CreateRequest(w http.ResponseWriter, r *http.Request) {
	logger.Info().Msgf("Received request create request")

	// Connect to NATS
	nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
	if ncErr != nil {
		logger.Err(ncErr).Msgf(helpers.NatsConnectError)
		helpers.HandleError(w, r, helpers.NatsConnectError)
		return
	}
	defer nc.Close()

	// Get body contents
	var federationRequest helpers.FederationRequestOperations
	json.NewDecoder(r.Body).Decode(&federationRequest)

	request := helpers.FederationRequestOperations{
		NodeURI:   federationRequest.NodeURI,
		NodeID:    federationRequest.NodeID,
		ChannelID: federationRequest.ChannelID,
		Type:      federationRequest.Type,
	}
	requestMarshal, marshalErr := json.Marshal(request)
	if marshalErr != nil {
		logger.Err(marshalErr).Msgf(helpers.MarshalErr)
		helpers.HandleError(w, r, helpers.MarshalErr)
		return
	}
	// Send the request
	msg, msgErr := nc.Request("federation.create", requestMarshal, helpers.TimeOut*time.Second)
	if msgErr != nil {
		logger.Err(msgErr).Msgf(helpers.MsgErr)
		helpers.HandleError(w, r, helpers.MsgErr)
		return
	}

	logger.Info().Msgf(helpers.Response+" %s\n", msg.Data)

	// Use the response
	var response helpers.FederationResultResponse
	unmarshalErr := json.Unmarshal([]byte(string(msg.Data)), &response)
	if unmarshalErr != nil {
		logger.Err(unmarshalErr).Msgf(helpers.UnmarshalErr)
		helpers.HandleError(w, r, helpers.UnmarshalErr)
		return
	}

	if response.Success {
		render.Render(w, r, responses.SuccessfulCreateResponse(response.Status))
	} else {
		render.Render(w, r, responses.ErrCustom(errors.New(response.Status)))
	}
}

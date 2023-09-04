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
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nats-io/nats.go"
)

// AcceptRequest is a controller function to accept the federation request
func AcceptRequest(w http.ResponseWriter, r *http.Request) {
	logger.Info().Msgf("Received request accept request")

	// Connect to NATS
	nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
	if ncErr != nil {
		logger.Err(ncErr).Msgf(helpers.NatsConnectError)
		helpers.HandleError(w, r, helpers.NatsConnectError)
		return
	}
	defer nc.Close()

	// Get body contents
	var acceptRequest helpers.FederationRequestOperations
	json.NewDecoder(r.Body).Decode(&acceptRequest)

	request := helpers.FederationRequestBody{
		RequestID: chi.URLParam(r, "request_id"),
		Type:      acceptRequest.Type,
	}

	requestMarshal, marshalErr := json.Marshal(request)
	if marshalErr != nil {
		logger.Err(marshalErr).Msgf(helpers.MarshalErr)
		helpers.HandleError(w, r, helpers.MarshalErr)
		return
	}
	// Send the request
	msgAccept, msgAcceptErr := nc.Request("federation.accept", requestMarshal, helpers.TimeOut*time.Second)
	if msgAcceptErr != nil {
		logger.Err(msgAcceptErr).Msgf(helpers.MsgErr)
		helpers.HandleError(w, r, helpers.MsgErr)
		return
	}

	var msgAcceptResponse helpers.FederationResultResponse
	unmarshalErr := json.Unmarshal([]byte(string(msgAccept.Data)), &msgAcceptResponse)
	if unmarshalErr != nil {
		logger.Err(unmarshalErr).Msgf(helpers.UnmarshalErr)
		helpers.HandleError(w, r, helpers.UnmarshalErr)
		return
	}

	if msgAcceptResponse.Success {
		msg, msgErr := nc.Request("federation.one", requestMarshal, helpers.TimeOut*time.Second)
		if msgErr != nil {
			logger.Err(msgErr).Msgf(helpers.MsgErr)
			helpers.HandleError(w, r, helpers.MsgErr)
			return
		}

		logger.Info().Msgf(helpers.Success+"%s\n", msg.Data)

		// Use the response
		var response helpers.FederationResultResponse
		unmarshalResErr := json.Unmarshal([]byte(string(msg.Data)), &response)
		if unmarshalResErr != nil {
			logger.Err(unmarshalResErr).Msgf(helpers.UnmarshalErr)
			helpers.HandleError(w, r, helpers.UnmarshalErr)
			return
		}
		acceptRequestSelf := helpers.FederationRequestOperations{
			Type:      "ACCEPT",
			NodeID:    helpers.GetNodeID(),
			ChannelID: response.Result[0].ChannelID,
		}
		marshalRequest, marshalErr := json.Marshal(acceptRequestSelf)
		if marshalErr != nil {
			logger.Err(marshalErr).Msgf(helpers.MarshalErr)
			helpers.HandleError(w, r, helpers.MarshalErr)
			return
		}
		var host = response.Result[0].NodeURI + ":7205"
		updateRes, updateErr := helpers.PostJSONRequest("https://"+host+"/api/v2/federation/requests/nodes/update", []byte(marshalRequest))
		if updateErr != nil {
			http.Error(w, "Error in updating data", http.StatusInternalServerError)
			helpers.HandleError(w, r, "Error in updating data")
			return
		}
		data, err := io.ReadAll(updateRes.Body)
		if err != nil {
			http.Error(w, "Error reading data", http.StatusInternalServerError)
			helpers.HandleError(w, r, "Error reading data")
			return
		}
		parsedData, err := helpers.ParseJSONData(data)
		logger.Info().Msgf("Parsed Data %v", parsedData)

		if err != nil {
			http.Error(w, "Error parsing JSON", http.StatusInternalServerError)
			helpers.HandleError(w, r, "Error parsing JSON")
			return
		}
		render.Render(w, r, responses.SuccessfulFederationAcceptResponse())
	} else {
		render.Render(w, r, responses.ErrCustom(errors.New(msgAcceptResponse.Status)))
	}
}

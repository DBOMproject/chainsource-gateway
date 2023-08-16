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

	"github.com/go-chi/render"
	"github.com/nats-io/nats.go"
)

// RevokeRequest is a controller function to revoke access
func RevokeRequest(w http.ResponseWriter, r *http.Request) {
	logger.Info().Msgf("Received request to revoke access")

	// Connect to NATS
	nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
	if ncErr != nil {
		logger.Err(ncErr).Msgf(helpers.NatsConnectError)
		helpers.HandleError(w, r, helpers.NatsConnectError)
		return
	}
	defer nc.Close()

	// Get body contents
	var revokeRequest helpers.FederationRequestOperations
	json.NewDecoder(r.Body).Decode(&revokeRequest)

	js, jsErr := nc.JetStream(nats.PublishAsyncMaxPending(helpers.PublishAsyncMaxPendingConstant))
	if jsErr != nil {
		logger.Err(jsErr).Msg(helpers.NatsJetStreamError)
		helpers.HandleError(w, r, helpers.NatsJetStreamError)
		return
	}

	js.AddStream(&nats.StreamConfig{
		Name:     "federation",
		Subjects: []string{"revoke"},
	})

	requestMarshal, marshalErr := json.Marshal(revokeRequest)
	if marshalErr != nil {
		logger.Err(marshalErr).Msgf(helpers.MarshalErr)
		helpers.HandleError(w, r, helpers.MarshalErr)
		return
	}

	js.PublishAsync("federation.revoke", requestMarshal)

	select {
	case <-js.PublishAsyncComplete():
		{
			var host = revokeRequest.NodeURI + ":7205"
			revokeRequestSelf := helpers.FederationRequestOperations{
				Type:      revokeRequest.Type,
				NodeID:    helpers.GetNodeID(),
				ChannelID: revokeRequest.ChannelID,
			}
			revokeRequestMarshal, marshalErr := json.Marshal(revokeRequestSelf)
			if marshalErr != nil {
				logger.Err(marshalErr).Msgf(helpers.MarshalErr)
				helpers.HandleError(w, r, helpers.MarshalErr)
				return
			}
			updateRes, updateErr := helpers.PostJSONRequest("https://"+host+"/api/v2/federation/requests/nodes/update", []byte(revokeRequestMarshal))
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
			render.Render(w, r, responses.SuccessfulFederationRevokeResponse())
		}
	case <-time.After(helpers.TimeOut * time.Second):
		render.Render(w, r, responses.ErrInvalidRequest(errors.New(helpers.TimeOutErr)))
	}
}

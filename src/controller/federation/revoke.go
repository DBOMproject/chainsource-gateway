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

	requestMarshal, marshalErr := json.Marshal(revokeRequest)
	if marshalErr != nil {
		logger.Err(marshalErr).Msgf(helpers.MarshalErr)
		helpers.HandleError(w, r, helpers.MarshalErr)
		return
	}

	// Send the request
	msgRevoke, msgRevokeErr := nc.Request("federation.revoke", requestMarshal, helpers.TimeOut*time.Second)
	if msgRevokeErr != nil {
		logger.Err(msgRevokeErr).Msgf(helpers.MsgErr)
		helpers.HandleError(w, r, helpers.MsgErr)
		return
	}
	logger.Info().Msgf(helpers.Response+" %s\n", msgRevoke.Data)

	var msgRevokeResponse helpers.FederationResultResponse
	unmarshalErr := json.Unmarshal([]byte(string(msgRevoke.Data)), &msgRevokeResponse)
	if unmarshalErr != nil {
		logger.Err(unmarshalErr).Msgf(helpers.UnmarshalErr)
		helpers.HandleError(w, r, helpers.UnmarshalErr)
		return
	}

	if msgRevokeResponse.Success {
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
			http.Error(w, helpers.UpdateErr, http.StatusInternalServerError)
			helpers.HandleError(w, r, helpers.UpdateErr)
			return
		}
		data, err := io.ReadAll(updateRes.Body)
		if err != nil {
			http.Error(w, helpers.ReadingErr, http.StatusInternalServerError)
			helpers.HandleError(w, r, helpers.ReadingErr)
			return
		}
		parsedData, err := helpers.ParseJSONData(data)
		logger.Info().Msgf("Parsed Data %v", parsedData)

		if err != nil {
			http.Error(w, helpers.ParseErr, http.StatusInternalServerError)
			helpers.HandleError(w, r, helpers.ParseErr)
			return
		}
		if msgRevokeResponse.Success {
			render.Render(w, r, responses.SuccessfulOkResponse(msgRevokeResponse.Status))
		} else {
			render.Render(w, r, responses.ErrCustom(errors.New(msgRevokeResponse.Status)))
		}
	}
}

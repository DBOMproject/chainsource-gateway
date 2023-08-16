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
	"chainsource-gateway/helpers"
	"chainsource-gateway/responses"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nats-io/nats.go"
)

// CreateChannel is a controller function to create channel
func CreateChannel(w http.ResponseWriter, r *http.Request) {
	var nodeUri = chi.URLParam(r, "node_uri")
	var nodeIdFromRequest = strings.Split(nodeUri, ".")[0]

	logger.Info().Msgf("Received request to create channel on node %s", nodeIdFromRequest)

	if nodeIdFromRequest == helpers.LocalNodeId || nodeIdFromRequest == helpers.GetNodeID() {
		// Connect to NATS
		nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
		if ncErr != nil {
			logger.Err(ncErr).Msgf(helpers.NatsConnectError)
			helpers.HandleError(w, r, helpers.NatsConnectError)
			return
		}
		defer nc.Close()

		// Get body contents
		var channel helpers.Channel
		json.NewDecoder(r.Body).Decode(&channel)

		js, jsErr := nc.JetStream(nats.PublishAsyncMaxPending(helpers.PublishAsyncMaxPendingConstant))
		if jsErr != nil {
			logger.Err(jsErr).Msg(helpers.NatsJetStreamError)
			helpers.HandleError(w, r, helpers.NatsJetStreamError)
			return
		}

		js.AddStream(&nats.StreamConfig{
			Name:     "channel",
			Subjects: []string{"create"},
		})

		request := helpers.CreateChannelMeta{
			ChannelMeta: channel,
		}

		requestMarshal, marshalErr := json.Marshal(request)
		if marshalErr != nil {
			logger.Err(marshalErr).Msgf(helpers.MarshalErr)
			helpers.HandleError(w, r, helpers.MarshalErr)
			return
		}

		js.PublishAsync("channel.create", requestMarshal)

		select {
		case <-js.PublishAsyncComplete():
			render.Render(w, r, responses.SuccessfulChannelsCreationResponse())
		case <-time.After(helpers.TimeOut * time.Second):
			helpers.HandleError(w, r, helpers.TimeOutErr)
			return
		}
	} else {
		render.Render(w, r, responses.ErrDoesNotExist(errors.New(helpers.InvalidRequest)))
	}
}

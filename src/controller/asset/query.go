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
	"chainsource-gateway/helpers"
	"chainsource-gateway/responses"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/nats-io/nats.go"
)

// QueryAsset is a controller function to rich query and query assets
func QueryAsset(w http.ResponseWriter, r *http.Request) {
	var nodeUri = chi.URLParam(r, "node_uri")
	var nodeIdFromRequest = strings.Split(nodeUri, ".")[0]
	var channelIdFromRequest = chi.URLParam(r, "channel_id")

	logger.Info().Msgf("Received request to query asset from node %s", nodeIdFromRequest)

	if nodeIdFromRequest == helpers.LocalNodeId || nodeIdFromRequest == helpers.GetNodeID() {
		// Connect to NATS
		nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
		if ncErr != nil {
			logger.Err(ncErr).Msg(helpers.NatsConnectError)
			helpers.HandleError(w, r, helpers.NatsConnectError)
			return
		}
		defer nc.Close()

		// Get body contents
		var query helpers.QueryMeta
		if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
			logger.Err(err).Msg(helpers.DecodeErr)
			helpers.HandleError(w, r, helpers.DecodeErr)
			return
		}

		js, jsErr := nc.JetStream(nats.PublishAsyncMaxPending(helpers.PublishAsyncMaxPendingConstant))
		if jsErr != nil {
			logger.Err(jsErr).Msg(helpers.NatsJetStreamError)
			helpers.HandleError(w, r, helpers.NatsJetStreamError)
			return
		}

		js.AddStream(&nats.StreamConfig{
			Name:     "asset",
			Subjects: []string{"query"},
		})

		request := helpers.AssetQueryRoutingVars{
			ChannelID: channelIdFromRequest,
			Query:     query,
		}

		requestMarshal, marshalErr := json.Marshal(request)
		if marshalErr != nil {
			logger.Err(marshalErr).Msg(helpers.MarshalErr)
			helpers.HandleError(w, r, helpers.MarshalErr)
			return
		}

		// Send the request
		msg, msgErr := nc.Request("asset.query", requestMarshal, helpers.TimeOut*time.Second)
		if msgErr != nil {
			logger.Err(msgErr).Msgf(helpers.MsgErr)
			helpers.HandleError(w, r, helpers.MsgErr)
			return
		}

		logger.Info().Msgf(helpers.Success+"%s\n", msg.Data)

		var response []helpers.AssetMeta
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

// RichQueryAsset is a controller function to rich query assets
func RichQueryAsset(w http.ResponseWriter, r *http.Request) {
	var nodeUri = chi.URLParam(r, "node_uri")
	var nodeIdFromRequest = strings.Split(nodeUri, ".")[0]
	var channelIdFromRequest = chi.URLParam(r, "channel_id")

	logger.Info().Msgf("Received request to rich query asset from node %s", nodeIdFromRequest)

	if nodeIdFromRequest == helpers.LocalNodeId || nodeIdFromRequest == helpers.GetNodeID() {
		// Connect to NATS
		nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
		if ncErr != nil {
			logger.Err(ncErr).Msg(helpers.NatsConnectError)
			helpers.HandleError(w, r, helpers.NatsConnectError)
			return
		}
		defer nc.Close()

		// Get body contents
		var query helpers.QueryMeta
		if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
			logger.Err(err).Msg(helpers.DecodeErr)
			helpers.HandleError(w, r, helpers.DecodeErr)
			return
		}

		js, jsErr := nc.JetStream(nats.PublishAsyncMaxPending(helpers.PublishAsyncMaxPendingConstant))
		if jsErr != nil {
			logger.Err(jsErr).Msg(helpers.NatsJetStreamError)
			helpers.HandleError(w, r, helpers.NatsJetStreamError)
			return
		}

		js.AddStream(&nats.StreamConfig{
			Name:     "asset",
			Subjects: []string{"query"},
		})

		request := helpers.AssetRichQueryRoutingVars{
			ChannelID: channelIdFromRequest,
			Query:     r.URL.Query().Get("query"),
			Fields:    r.URL.Query().Get("fields"),
			Limit:     r.URL.Query().Get("limit"),
			Skip:      r.URL.Query().Get("skip"),
		}

		requestMarshal, marshalErr := json.Marshal(request)
		if marshalErr != nil {
			logger.Err(marshalErr).Msg(helpers.MarshalErr)
			helpers.HandleError(w, r, helpers.MarshalErr)
			return
		}

		// Send the request
		msg, msgErr := nc.Request("asset.richquery", requestMarshal, helpers.TimeOut*time.Second)
		if msgErr != nil {
			logger.Err(msgErr).Msgf(helpers.MsgErr)
			helpers.HandleError(w, r, helpers.MsgErr)
			return
		}

		logger.Info().Msgf(helpers.Success+"%s\n", msg.Data)

		var response []helpers.AssetMeta
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

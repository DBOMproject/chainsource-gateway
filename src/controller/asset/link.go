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

// LinkAsset is a controller function to add link to an asset
func LinkAsset(w http.ResponseWriter, r *http.Request) {
	var nodeUri = chi.URLParam(r, "node_uri")
	var nodeIdFromRequest = strings.Split(nodeUri, ".")[0]
	var channelIdFromRequest = chi.URLParam(r, "channel_id")
	var assetIdFromRequest = chi.URLParam(r, "asset_id")

	logger.Info().Msgf("Received request to link asset %s on channel %s", assetIdFromRequest, channelIdFromRequest)

	if nodeIdFromRequest == helpers.LocalNodeId || nodeIdFromRequest == helpers.GetNodeID() {
		// Connect to NATS
		nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
		if ncErr != nil {
			logger.Err(ncErr).Msg(helpers.NatsConnectError)
			helpers.HandleError(w, r, helpers.NatsConnectError)
			return
		}
		defer nc.Close()

		var assetLink helpers.LinksDefinition
		if err := json.NewDecoder(r.Body).Decode(&assetLink); err != nil {
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
			Subjects: []string{"link"},
		})

		request := helpers.LinkAssetMeta{
			ChannelID: channelIdFromRequest,
			AssetID:   assetIdFromRequest,
			Payload:   assetLink,
		}

		requestMarshal, marshalErr := json.Marshal(request)
		if marshalErr != nil {
			logger.Err(marshalErr).Msg(helpers.MarshalErr)
			helpers.HandleError(w, r, helpers.MarshalErr)
			return
		}

		// Send NATS request
		js.PublishAsync("asset.link", requestMarshal)

		select {
		case <-js.PublishAsyncComplete():
			render.Render(w, r, responses.SuccessfulLinkResponse())
		case <-time.After(helpers.TimeOut * time.Second):
			helpers.HandleError(w, r, helpers.TimeOutErr)
		}
	} else {
		render.Render(w, r, responses.ErrDoesNotExist(errors.New(helpers.InvalidRequest)))
	}
}

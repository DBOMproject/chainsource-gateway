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

// UnlinkAsset is a controller function to remove a link to another asset
func UnlinkAsset(w http.ResponseWriter, r *http.Request) {
	var nodeUri = chi.URLParam(r, "node_uri")
	var nodeIdFromRequest = strings.Split(nodeUri, ".")[0]
	var channelIdFromRequest = chi.URLParam(r, "channel_id")
	var assetIdFromRequest = chi.URLParam(r, "asset_id")
	var linkIdFromRequest = chi.URLParam(r, "link_id")

	logger.Info().Msgf("Received request to unlink asset %s from asset %s on channel %s from node %s", assetIdFromRequest, linkIdFromRequest, channelIdFromRequest, nodeIdFromRequest)

	// If the request is for this node, audit the asset locally
	if nodeIdFromRequest == helpers.LocalNodeId || nodeIdFromRequest == helpers.GetNodeID() {
		// Connect to NATS
		nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
		if ncErr != nil {
			logger.Err(ncErr).Msgf(helpers.NatsConnectError)
			helpers.HandleError(w, r, helpers.NatsConnectError)
			return
		}
		defer nc.Close()

		js, jsErr := nc.JetStream(nats.PublishAsyncMaxPending(helpers.PublishAsyncMaxPendingConstant))
		if jsErr != nil {
			logger.Err(jsErr).Msg(helpers.NatsJetStreamError)
			helpers.HandleError(w, r, helpers.NatsJetStreamError)
			return
		}

		js.AddStream(&nats.StreamConfig{
			Name:     "asset",
			Subjects: []string{"unlink"},
		})

		request := helpers.UnlinkAssetMeta{
			ChannelID: channelIdFromRequest,
			AssetID:   assetIdFromRequest,
			LinkID:    linkIdFromRequest,
		}

		requestMarshal, marshalErr := json.Marshal(request)
		if marshalErr != nil {
			logger.Err(marshalErr).Msg(helpers.MarshalErr)
			helpers.HandleError(w, r, helpers.MarshalErr)
			return
		}

		// Send NATS request
		js.PublishAsync("asset.unlink", requestMarshal)

		select {
		case <-js.PublishAsyncComplete():
			render.Render(w, r, responses.SuccessfulUnlinkResponse())
		case <-time.After(helpers.TimeOut * time.Second):
			helpers.HandleError(w, r, helpers.TimeOutErr)
			return
		}
	} else {
		render.Render(w, r, responses.ErrDoesNotExist(errors.New(helpers.InvalidRequest)))
	}
}

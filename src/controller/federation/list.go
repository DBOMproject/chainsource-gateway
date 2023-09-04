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
 
	 "github.com/go-chi/chi"
	 "github.com/go-chi/render"
	 "github.com/nats-io/nats.go"
 )
 
 // TODO: Validate the request
 func Validate(w http.ResponseWriter, r *http.Request) bool {
	 return true
 }
 
 var logger = helpers.GetLogger(helpers.FederationLogger)
 
 // GetAllRequest is a controller function to list all the federation requests
 func GetAllRequest(w http.ResponseWriter, r *http.Request) {
	 logger.Info().Msgf("Received request to list all federation requests")
 
	 nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
	 if ncErr != nil {
		 logger.Err(ncErr).Msgf(helpers.NatsConnectError)
		 helpers.HandleError(w, r, helpers.NatsConnectError)
		 return
	 }
	 defer nc.Close()
 
	 // Send the request
	 msg, msgErr := nc.Request("federation.all", nil, helpers.TimeOut*time.Second)
	 if msgErr != nil {
		 logger.Err(msgErr).Msgf(helpers.MsgErr)
		 helpers.HandleError(w, r, helpers.MsgErr)
		 return
	 }
 
	 logger.Info().Msgf(helpers.Success+"%s\n", msg.Data)
 
	 // Use the response
	 var response helpers.FederationResultResponse
	 unmarshalErr := json.Unmarshal([]byte(string(msg.Data)), &response)
	 if unmarshalErr != nil {
		 logger.Err(unmarshalErr).Msgf(helpers.UnmarshalErr)
		 helpers.HandleError(w, r, helpers.UnmarshalErr)
		 return
	 }
	 render.JSON(w, r, response)
 }
 
 // GetOneRequest is a controller function to list one federation request
 func GetOneRequest(w http.ResponseWriter, r *http.Request) {
	 logger.Info().Msgf("Received request to list one federation requests")
 
	 nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
	 if ncErr != nil {
		 logger.Err(ncErr).Msgf(helpers.NatsConnectError)
		 helpers.HandleError(w, r, helpers.NatsConnectError)
		 return
	 }
	 defer nc.Close()
 
	 // Get request Id
	 request := helpers.FederationRoutingVars{
		 RequestID: chi.URLParam(r, "request_id"),
	 }
 
	 requestMarshal, marshalErr := json.Marshal(request)
	 if marshalErr != nil {
		 logger.Err(marshalErr).Msgf(helpers.MarshalErr)
		 helpers.HandleError(w, r, helpers.MarshalErr)
		 return
	 }
 
	 // Send the request
	 msg, msgErr := nc.Request("federation.one", requestMarshal, helpers.TimeOut*time.Second)
	 if msgErr != nil {
		 logger.Err(msgErr).Msgf(helpers.MsgErr)
		 helpers.HandleError(w, r, helpers.MsgErr)
		 return
	 }
 
	 logger.Info().Msgf(helpers.Success+"%s\n", msg.Data)
 
	 // Use the response
	 var response helpers.FederationResultResponse
	 unmarshalErr := json.Unmarshal([]byte(string(msg.Data)), &response)
	 if unmarshalErr != nil {
		 logger.Err(unmarshalErr).Msgf(helpers.UnmarshalErr)
		 helpers.HandleError(w, r, helpers.UnmarshalErr)
		 return
	 }
	 render.JSON(w, r, response)
 }
 
 // FedListOneChannel is a controller function to lists channel by ID
 func FedListOneChannel(w http.ResponseWriter, r *http.Request) {
	 logger.Info().Msgf("Received request to list one channel")
 
	 if Validate(w, r) {
		 nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
		 if ncErr != nil {
			 logger.Err(ncErr).Msgf(helpers.NatsConnectError)
			 helpers.HandleError(w, r, helpers.NatsConnectError)
			 return
		 }
		 defer nc.Close()
 
		 request := helpers.ListOneChannelMeta{
			 ChannelID: chi.URLParam(r, "channel_id"),
		 }
 
		 requestMarshal, marshalErr := json.Marshal(request)
		 if marshalErr != nil {
			 logger.Err(marshalErr).Msgf(helpers.MarshalErr)
			 helpers.HandleError(w, r, helpers.MarshalErr)
			 return
		 }
 
		 // Send the request
		 msg, msgErr := nc.Request("channel.one", requestMarshal, helpers.TimeOut*time.Second)
		 if msgErr != nil {
			 logger.Err(msgErr).Msgf(helpers.MsgErr)
			 helpers.HandleError(w, r, helpers.MsgErr)
			 return
		 }
 
		 logger.Info().Msgf(helpers.Success+"%s\n", msg.Data)
 
		 // Use the response
		 var response []helpers.Channel
		 unmarshalErr := json.Unmarshal([]byte(string(msg.Data)), &response)
		 if unmarshalErr != nil {
			 logger.Err(unmarshalErr).Msgf(helpers.UnmarshalErr)
			 helpers.HandleError(w, r, helpers.UnmarshalErr)
			 return
		 }
		 render.JSON(w, r, response)
	 } else {
		 render.Render(w, r, responses.ErrUnauthorizedQueryDestination(errors.New("Unauthorized")))
	 }
 }
 
 // FedListAssets is a controller function to list all the asset
 func FedListAssets(w http.ResponseWriter, r *http.Request) {
	 logger.Info().Msgf("Received request to list all assets")
 
	 if Validate(w, r) {
		 nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
		 if ncErr != nil {
			 logger.Err(ncErr).Msgf(helpers.NatsConnectError)
			 helpers.HandleError(w, r, helpers.NatsConnectError)
			 return
		 }
		 defer nc.Close()
 
		 request := helpers.AssetRoutingVars{
			 ChannelID: chi.URLParam(r, "channel_id"),
		 }
 
		 requestMarshal, marshalErr := json.Marshal(request)
		 if marshalErr != nil {
			 logger.Err(marshalErr).Msgf(helpers.MarshalErr)
			 helpers.HandleError(w, r, helpers.MarshalErr)
			 return
		 }
 
		 // Send the request
		 msg, msgErr := nc.Request("asset.all", requestMarshal, helpers.TimeOut*time.Second)
		 if msgErr != nil {
			 logger.Err(msgErr).Msgf(helpers.MsgErr)
			 helpers.HandleError(w, r, helpers.MsgErr)
			 return
		 }
 
		 logger.Info().Msgf(helpers.Success+"%s\n", msg.Data)
 
		 // Use the response
		 var response []helpers.AssetMeta
		 unmarshalErr := json.Unmarshal([]byte(string(msg.Data)), &response)
		 if unmarshalErr != nil {
			 logger.Err(unmarshalErr).Msgf(helpers.UnmarshalErr)
			 helpers.HandleError(w, r, helpers.UnmarshalErr)
			 return
		 }
		 render.JSON(w, r, response)
	 } else {
		 render.Render(w, r, responses.ErrUnauthorizedQueryDestination(errors.New("Unauthorized")))
	 }
 }
 
 // FedListOneAsset is a controller function to lists asset by ID
 func FedListOneAsset(w http.ResponseWriter, r *http.Request) {
	 logger.Info().Msgf("Received request to list one asset")
 
	 if Validate(w, r) {
		 nc, ncErr := nats.Connect(os.Getenv("NATS_URI"))
		 if ncErr != nil {
			 logger.Err(ncErr).Msgf(helpers.NatsConnectError)
			 helpers.HandleError(w, r, helpers.NatsConnectError)
			 return
		 }
		 defer nc.Close()
 
		 request := helpers.AssetRoutingVars{
			 ChannelID: chi.URLParam(r, "channel_id"),
			 AssetID:   chi.URLParam(r, "asset_id"),
		 }
 
		 requestMarshal, marshalErr := json.Marshal(request)
		 if marshalErr != nil {
			 logger.Err(marshalErr).Msgf(helpers.MarshalErr)
			 helpers.HandleError(w, r, helpers.MarshalErr)
			 return
		 }
 
		 // Send the request
		 msg, msgErr := nc.Request("asset.one", requestMarshal, helpers.TimeOut*time.Second)
		 if msgErr != nil {
			 logger.Err(msgErr).Msgf(helpers.MsgErr)
			 helpers.HandleError(w, r, helpers.MsgErr)
			 return
		 }
 
		 logger.Info().Msgf(helpers.Success+"%s\n", msg.Data)
 
		 // Use the response
		 var response []helpers.AssetMeta
		 unmarshalErr := json.Unmarshal([]byte(string(msg.Data)), &response)
		 if unmarshalErr != nil {
			 logger.Err(unmarshalErr).Msgf(helpers.UnmarshalErr)
			 helpers.HandleError(w, r, helpers.UnmarshalErr)
			 return
		 }
		 render.JSON(w, r, response)
	 } else {
		 render.Render(w, r, responses.ErrUnauthorizedQueryDestination(errors.New("Unauthorized")))
	 }
 }
 
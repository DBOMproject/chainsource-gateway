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

package helpers

const (
	// Errors
	NatsConnectError   = "error connecting to NATS"
	NatsJetStreamError = "error connecting to JetStream"
	MarshalErr         = "error marshalling message"
	UnmarshalErr       = "error unmarshalling message"
	MsgErr             = "error sending message"
	TimeOutErr         = "timeout"
	NotImplemented     = "not implemented"
	ResponseErr        = "error in receiving response"
	ReadingErr         = "error reading data"
	DecodeErr          = "error decoding data"
	ParseErr           = "error parsing JSON"
	Unauthorized       = "requested action not authorized"
	InvalidRequest     = "invalid request"
	UpdateErr          = "error in updating data"

	// Results
	Response = "received response:"

	// Loggers
	AssetLogger      = "[ASSET]"
	ChannelLogger    = "[CHANNEL]"
	FederationLogger = "[FEDERATION]"
	NodeLogger       = "[NODE]"
	HttpClient       = "[HttpClient]"
	Main             = "[MAIN]"

	// Constant Numbers
	TimeOut                        = 5
	PublishAsyncMaxPendingConstant = 256

	// Local Constants
	LocalNodeId = "_local"
)

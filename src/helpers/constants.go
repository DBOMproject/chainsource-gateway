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
	NatsConnectError   = "Error connecting to NATS"
	NatsJetStreamError = "Error connecting to JetStream"
	MarshalErr         = "Error marshalling message"
	UnmarshalErr       = "Error unmarshalling message"
	MsgErr             = "Error sending message"
	TimeOutErr         = "timeout"
	NotImplemented     = "Not Implemented"
	ResponseErr        = "Error in receiving response"
	ReadingErr         = "Error reading data"
	DecodeErr          = "Error decoding data"
	Unauthorized       = "Requested action not authorized"
	InvalidRequest     = "invalid request, please check the request body and try again"

	// Results
	Success = "Success - Received response:"

	// Loggers
	AssetLogger      = "Asset"
	ChannelLogger    = "Channel"
	FederationLogger = "Federation"
	NodeLogger       = "Node"

	// Constant Numbers
	TimeOut                        = 5
	PublishAsyncMaxPendingConstant = 256

	// Local Constants
	LocalNodeId = "_local"
)

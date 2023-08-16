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

// Package responses is a module that contains all the responses rendered by the gateway as JSON marshal-able structs
package responses

import (
	"net/http"

	"github.com/go-chi/render"
)

// Assets

// SuccessResponse is a type for a success response
type SuccessResponse struct {
	HTTPStatusCode int    `json:"-"`
	IsSuccessful   bool   `json:"success"`
	StatusText     string `json:"status"`
	Result         string `json:"result,omitempty"`
}

// Render renders the success response
func (e SuccessResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func successfulResponse(code int, statusText string) render.Renderer {
	return &SuccessResponse{
		IsSuccessful:   true,
		HTTPStatusCode: code,
		StatusText:     statusText,
	}
}

/** Node **/
// SuccessfulNodeInitResponse returns success when request is sent
func SuccessfulNodeInitResponse() render.Renderer {
	return successfulResponse(201, "Successfully sent request to initialize node")
}

/** Asset **/

// SuccessfulAssetCreationResponse returns success when request is sent
func SuccessfulAssetCreationResponse() render.Renderer {
	return successfulResponse(201, "Successfully sent request to create asset")
}

// SuccessfulAssetUpdateResponse returns success when request is sent
func SuccessfulAssetUpdateResponse() render.Renderer {
	return successfulResponse(200, "Successfully sent request to update asset")
}

// SuccessfulLinkResponse returns success for when request is sent
func SuccessfulLinkResponse() render.Renderer {
	return successfulResponse(200, "Successfully sent request to link asset")
}

// SuccessfulUnlinkResponse returns success when request is sent
func SuccessfulUnlinkResponse() render.Renderer {
	return successfulResponse(200, "Successfully sent request to unlink asset")
}

/** Channels **/

// SuccessfulChannelsCreationResponse returns success when request is sent
func SuccessfulChannelsCreationResponse() render.Renderer {
	return successfulResponse(201, "Successfully sent request to create channel")
}

// SuccessfulChannelUpdateResponse returns success when request is sent
func SuccessfulChannelUpdateResponse() render.Renderer {
	return successfulResponse(200, "Successfully sent request to update notary details to a channel")
}

// SuccessfulChannelUpdateResponse returns success when request is sent
func SuccessfulChannelDeleteResponse() render.Renderer {
	return successfulResponse(200, "Successfully sent request to delete notary details from a channel")
}

/** Federation **/

// SuccessfulFederationSentResponse
func SuccessfulFederationSentResponse() render.Renderer {
	return successfulResponse(200, "Successfully sent federation request")
}

func SuccessfulFederationAcceptResponse() render.Renderer {
	return successfulResponse(200, "Successfully accepted federation request")
}
func SuccessfulFederationRejectResponse() render.Renderer {
	return successfulResponse(200, "Successfully rejected federation request")
}

func SuccessfulFederationRevokeResponse() render.Renderer {
	return successfulResponse(200, "Successfully sent revoke request")
}

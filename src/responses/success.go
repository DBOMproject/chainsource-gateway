/*
 * Copyright 2020 Unisys Corporation
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

//SuccessResponse is a type for a success response
type SuccessResponse struct {
	HTTPStatusCode int `json:"-"`

	IsSuccessful  bool   `json:"success"`
	StatusText    string `json:"status"`
	AgentResponse string `json:"agentResponse,omitempty"`
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

//SuccessfulCreationResponse returns success for when an asset is created on an agent
func SuccessfulCreationResponse() render.Renderer {
	return successfulResponse(201, "Successfully created on the agent")
}

//SuccessfulUpdateResponse returns success for when an asset is updated on an agent
func SuccessfulUpdateResponse() render.Renderer {
	return successfulResponse(200, "Successfully updated on the agent")
}

//SuccessfulAttachResponse returns success for when a subasset is attached on an agent
func SuccessfulAttachResponse() render.Renderer {
	return successfulResponse(200, "Successfully attached on agent")
}

//SuccessfulTransferResponse returns success when a subasset is attached on an agent
func SuccessfulTransferResponse() render.Renderer {
	return successfulResponse(200, "Successfully transferred asset")
}

//SuccessfulDetachResponse returns success when a subasset is detached on an agent
func SuccessfulDetachResponse() render.Renderer {
	return successfulResponse(200, "Successfully detached on agent")
}

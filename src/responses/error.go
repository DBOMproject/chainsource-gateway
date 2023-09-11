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

package responses

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrResponse is a type for a error response
type ErrResponse struct {
	IsSuccessful   bool  `json:"success"`
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// Render renders the error response
func (e ErrResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// NATS errors
// NatsURLError returns error for when an invalid request is received
func NatsURLError(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "NATS connection error",
		ErrorText:      err.Error(),
	}
}

// Custom errors
func ErrCustom(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "invalid request",
		ErrorText:      err.Error(),
	}
}

// ErrInvalidRequest returns error for when an invalid request is received
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "invalid request",
		ErrorText:      err.Error(),
	}
}

// ErrAlreadyExists returns error for when an asset already exists
func ErrAlreadyExists(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusConflict,
		StatusText:     "already exists",
		ErrorText:      err.Error(),
	}
}

// ErrUnauthorizedQueryDestination returns error for when querying the destination asset fails
func ErrUnauthorizedQueryDestination(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "unauthorized to access destination channel",
		ErrorText:      err.Error(),
	}
}

// ErrUnimplemented returns error for when a method is not implemented
func ErrUnimplemented(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusNotImplemented,
		StatusText:     "unimplemented method",
		ErrorText:      err.Error(),
	}
}

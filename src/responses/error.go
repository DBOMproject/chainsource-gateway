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

package responses

import (
	"net/http"

	"github.com/go-chi/render"
)

//ErrResponse is a type for a error response
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

// User-Errors (4xx)

//ErrInvalidRequest returns error for when an invalid request is received
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request",
		ErrorText:      err.Error(),
	}
}

//ErrAlreadyExists returns error for when an asset already exists
func ErrAlreadyExists(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusConflict,
		StatusText:     "Already Exists",
		ErrorText:      err.Error(),
	}
}

//ErrDoesNotExist returns error for when an asset does not exist
func ErrDoesNotExist(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     "Asset or Channel does not exist",
		ErrorText:      err.Error(),
	}
}

//ErrAgentUnauthorized returns error for when an agent is unauthorized to access a channel
func ErrAgentUnauthorized(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Agent Failure",
		ErrorText:      err.Error(),
	}
}


//ErrNoAgent returns error for when an agent does not exist
func ErrNoAgent(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     "No agent",
		ErrorText:      err.Error(),
	}
}

//ErrAlreadyAttached returns error for when a sub asset is already attached
func ErrAlreadyAttached(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusForbidden,
		StatusText:     "Subasset already attached",
		ErrorText:      err.Error(),
	}
}

//ErrNotAttached returns error for when a sub asset is not attached
func ErrNotAttached(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusForbidden,
		StatusText:     "Subasset not attached",
		ErrorText:      err.Error(),
	}
}

//ErrConflict returns error for when conflicting information is found
func ErrConflict(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusConflict,
		StatusText:     "Conflicting Information",
		ErrorText:      err.Error(),
	}
}

//ErrNoSignature returns the json response for when no signature is found during validation
func ErrNoSignature() render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     "No signature found",
	}
}

//ErrInvalidSignature returns the json response for when an invalid signature is found during validation
func ErrInvalidSignature() render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		HTTPStatusCode: 400,
		StatusText:     "Invalid Signature",
	}
}

//ErrNoFingerprint returns the json response for when no fingerprint is found during validation
func ErrNoFingerprint() render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     "No signature fingerprint found",
	}
}

//ErrReadOnly returns the json response for when an asset is read only
func ErrReadOnly() render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		HTTPStatusCode: http.StatusConflict,
		StatusText:     "Asset is read only",
	}
}

// Gateway/Agent-Errors (5xx)

//ErrAgent returns error for when there is a failure on the agent
func ErrAgent(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusBadGateway,
		StatusText:     "Agent Failure",
		ErrorText:      err.Error(),
	}
}

//ErrFailedQueryChild returns error for when querying the child asset fails
func ErrFailedQueryChild(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusBadGateway,
		StatusText:     "Failed to get child",
		ErrorText:      err.Error(),
	}
}

//ErrUnauthorizedQueryChild returns error for when querying the child asset fails
func ErrUnauthorizedQueryChild(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Unauthorized to get child",
		ErrorText:      err.Error(),
	}
}

//ErrFailedQueryDestination returns error for when querying the destination asset fails
func ErrFailedQueryDestination(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusBadGateway,
		StatusText:     "Unable to query to the destination channel",
		ErrorText:      err.Error(),
	}
}

//ErrUnauthorizedQueryDestination returns error for when querying the destination asset fails
func ErrUnauthorizedQueryDestination(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Unauthorized to access destination channel",
		ErrorText:      err.Error(),
	}
}

//ErrFailedModifyChild returns error for when when an attach/detach update on a child asset occurs
func ErrFailedModifyChild(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusBadGateway,
		StatusText:     "Failed to update parent entry in child on agent",
		ErrorText:      err.Error(),
	}
}

//ErrFailedModifyParent returns error for when when an attach/detach update on a parent asset occurs
func ErrUnauthorizedModifyChild(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Unauthorized to update parent entry in child on agent",
		ErrorText:      err.Error(),
	}
}

//ErrFailedModifyParent returns error for when when an attach/detach update on a parent asset occurs
func ErrFailedModifyParent(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusBadGateway,
		StatusText:     "Failed to update subasset list in parent on agent",
		ErrorText:      err.Error(),
	}
}

//ErrFailedModifyParent returns error for when when an attach/detach update on a parent asset occurs
func ErrUnauthorizedModifyParent(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Unauthorized to update subasset list in parent on agent",
		ErrorText:      err.Error(),
	}
}

//ErrInternalServer returns error for when an internal server error occurs
func ErrInternalServer(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Validator Failure",
		ErrorText:      err.Error(),
	}
}

//ErrUnimplemented returns error for when a method is not implemented
func ErrUnimplemented(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusNotImplemented,
		StatusText:     "Unimplemented Method",
		ErrorText:      err.Error(),
	}
}

//ErrFailedExport returns the json response for when an export fails
func ErrFailedExport(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusBadGateway,
		StatusText:     "Failed to export asset",
		ErrorText:      err.Error(),
	}
}

//ErrFailedModifyOrigin returns the json response for when the origin asset could not be updated
func ErrFailedModifyOrigin(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusBadGateway,
		StatusText:     "Failed to update origin asset on agent",
		ErrorText:      err.Error(),
	}
}

//ErrUnauthorizedModifyOrigin returns the json response for when the origin asset could not be updated due to no authorization
func ErrUnauthorizedModifyOrigin(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Unauthorized to update origin asset on agent",
		ErrorText:      err.Error(),
	}
}


//ErrFailedModifyDestination returns the json response for when the destination asset could not be updated
func ErrFailedModifyDestination(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusBadGateway,
		StatusText:     "Failed to update destination on agent",
		ErrorText:      err.Error(),
	}
}

//ErrUnauthorizedModifyDestination returns the json response for when the destination asset could not be updated
func ErrUnauthorizedModifyDestination(err error) render.Renderer {
	return &ErrResponse{
		IsSuccessful:   false,
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     "Unauthorized to update destination on agent",
		ErrorText:      err.Error(),
	}
}


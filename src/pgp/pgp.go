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

// Package pgp contains functions and types for validating a pgp signature
package pgp

import (
	"chainsource-gateway/helpers"
	"chainsource-gateway/tracing"
	"context"
	"encoding/json"

	"github.com/opentracing/opentracing-go"
)

// SignatureValidator is an interface that provides the Validate API
type SignatureValidator interface {
	Validate(ctx context.Context, args ValidateArgs) (result map[string]interface{}, err error)
}

// SigningServiceValidator is an implementation of the SignatureValidator interface, utilizing the signing service
type SigningServiceValidator struct {
}

// NewSigningServiceValidator returns a new signing service validator
func NewSigningServiceValidator() SigningServiceValidator {
	return SigningServiceValidator{}
}

// ValidateArgs is a type representing the arguments sent to the pgp service to validate a signature
type ValidateArgs struct {
	Fingerprint string
	Signature   string
	Input       helpers.AssetNoChildParent
}

// Validate a pgp signature
func (SigningServiceValidator) Validate(ctx context.Context, args ValidateArgs) (result map[string]interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "PGP Validate")
	defer span.Finish()
	url := GetPGPServiceAddress()

	span.SetTag("pgp-service-url", url)
	validatePath := "/pgp/validate"

	body := helpers.ValidateBody{
		Fingerprint: args.Fingerprint,
		Signature:   args.Signature,
		Input:       args.Input,
	}
	bytesRepresentation, err := json.Marshal(body)
	if err != nil {
		tracing.LogAndTraceErr(log, span, err, "Validate signature failed")
	}
	result, err = helpers.PostJSONRequest(url, validatePath, nil, bytesRepresentation)

	if err != nil {
		tracing.LogAndTraceErr(log, span, err, "Validate signature failed")
	}
	return
}

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

package pgp

import (
	"chainsource-gateway/helpers"
	"context"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"math"
	"net/http"
	"os"
	"testing"
)

const mockSigningServiceEndpoint = "http://mock-signer"

// setupMockSigningService creates a gock endpoint that emulates the pgp signing service
func setupMockSigningService(statusCode int){
	gock.New(mockSigningServiceEndpoint).
		Post("/pgp/validate").
		Reply(statusCode).
		JSON(map[string]string{"foo": "bar"})
}

// TestNewSigningServiceValidator tests the Signing Service Validator getter
func TestNewSigningServiceValidator(t *testing.T) {
	assert.Implements(t, (*SignatureValidator)(nil),  NewSigningServiceValidator(), "Implements signature validator interface")
}

// TestSigningServiceValidator_Validate tests the validate implementation of the Signing Service Validator
func TestSigningServiceValidator_Validate(t *testing.T) {
	os.Setenv(pgpServiceAddressVar, mockSigningServiceEndpoint)
	defer os.Unsetenv(pgpServiceAddressVar)
	t.Run("When_Remote_Success", func(t *testing.T) {
		setupMockSigningService(http.StatusOK)
		defer gock.Off()

		signingServiceValidator := NewSigningServiceValidator()
		_, err := signingServiceValidator.Validate(context.Background(), ValidateArgs{
			Fingerprint: "",
			Signature:   "",
			Input:       helpers.AssetNoChildParent{},
		})

		assert.NoError(t, err, "No error must be returned")
	})
	t.Run("When_Remote_Failure", func(t *testing.T) {
		setupMockSigningService(http.StatusInternalServerError)
		defer gock.Off()

		signingServiceValidator := NewSigningServiceValidator()
		_, err := signingServiceValidator.Validate(context.Background(), ValidateArgs{
			Fingerprint: "",
			Signature:   "",
			Input:       helpers.AssetNoChildParent{},
		})

		assert.Error(t, err, "Error must be returned")
	})
	t.Run("When_Bad_Payload", func(t *testing.T) {

		signingServiceValidator := NewSigningServiceValidator()
		_, err := signingServiceValidator.Validate(context.Background(), ValidateArgs{
			Fingerprint: "",
			Signature:   "",
			Input:       helpers.AssetNoChildParent{
				StandardVersion: math.Inf(1),
			},
		})

		assert.Error(t, err, "Error must be returned")
	})
}

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
	"github.com/magiconair/properties/assert"
	"os"
	"testing"
)

const testAddr = "http://addr"

// TestGetPGPServiceAddress tests the pgp service address getter
func TestGetPGPServiceAddress(t *testing.T) {
	t.Run("When_Environment_Set", func(t *testing.T) {
		os.Setenv(pgpServiceAddressVar, testAddr)
		defer os.Unsetenv(pgpServiceAddressVar)
		assert.Equal(t, testAddr, GetPGPServiceAddress(), "Appropriate address returned")
	})
	t.Run("When_Environment_Not_Set", func(t *testing.T) {
		assert.Equal(t, defaultPgpServiceAddress, GetPGPServiceAddress(), "Appropriate address returned")
	})
}



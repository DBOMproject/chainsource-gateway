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

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetServiceAddress tests the method that gets the go chi address from the env vars
func TestGetServiceAddress(t *testing.T) {
	t.Run("When_Set_Environment", func(t *testing.T) {
		os.Setenv("PORT", "3050")
		defer os.Unsetenv("PORT")

		address := GetServiceAddress()
		assert.Equal(t, ":3050", address, "Should return the set address")
	})

	t.Run("When_No_Environment", func(t *testing.T) {
		address := GetServiceAddress()
		assert.Equal(t, ":3050", address, "Should return the default address")
	})
}
func TestGetFederationAddress(t *testing.T) {
	t.Run("When_Set_Environment", func(t *testing.T) {
		os.Setenv("FED_PORT", "7205")
		defer os.Unsetenv("FED_PORT")

		address := GetFederationAddress()
		assert.Equal(t, ":7205", address, "Should return the set address")
	})

	t.Run("When_No_Environment", func(t *testing.T) {
		address := GetFederationAddress()
		assert.Equal(t, ":7205", address, "Should return the default address")
	})
}

func TestGetNodeID(t *testing.T) {
	t.Run("When_Set_Environment", func(t *testing.T) {
		expectedNodeID := "my-node-id"
		os.Setenv("NODE_ID", expectedNodeID)
		defer os.Unsetenv("NODE_ID")

		nodeID := GetNodeID()
		assert.Equal(t, expectedNodeID, nodeID, "Should return the set Node ID")
	})

	t.Run("When_No_Environment", func(t *testing.T) {
		nodeID := GetNodeID()
		assert.Equal(t, "", nodeID, "Should return an empty string when not set")
	})
}

func TestGetNodeURI(t *testing.T) {
	t.Run("When_Set_Environment", func(t *testing.T) {
		expectedNodeURI := "my-node-uri"
		os.Setenv("NODE_URI", expectedNodeURI)
		defer os.Unsetenv("NODE_URI")

		nodeURI := GetNodeURI()
		assert.Equal(t, expectedNodeURI, nodeURI, "Should return the set Node URI")
	})

	t.Run("When_No_Environment", func(t *testing.T) {
		nodeURI := GetNodeURI()
		assert.Equal(t, "", nodeURI, "Should return an empty string when not set")
	})
}

// Test_existsInEnv tests the function that checks if a key exists in the environment
func TestExistsInEnv(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		os.Setenv("foo", "bar")
		defer os.Unsetenv("foo")

		exists := ExistsInEnv("foo")
		assert.True(t, exists, "Should return true for existing environment variable")
	})

	t.Run("DoesNotExist", func(t *testing.T) {
		exists := ExistsInEnv("bar")
		assert.False(t, exists, "Should return false for nonexistent environment variable")
	})
}

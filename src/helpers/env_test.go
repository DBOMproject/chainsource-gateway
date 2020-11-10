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

package helpers

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestGetServiceAddress tests the method that gets the go chi address from the env vars
func TestGetServiceAddress(t *testing.T) {
	t.Run("When_Set_Environment", func(t *testing.T) {
		os.Setenv("PORT", "3000")
		defer os.Unsetenv("PORT")

		assert.NotPanics(t, func() {
			GetServiceAddress()
		}, "Does not panic while getting")
	})
	t.Run("When_No_Environment", func(t *testing.T) {
		assert.NotPanics(t, func() {
			GetServiceAddress()
		}, "Does not panic while getting")
	})
}

// Test_existsInEnv tests the function that checks if a key exists in the environment
func TestExistsInEnv(t *testing.T) {
	os.Setenv("foo", "bar")
	defer os.Unsetenv("foo")
	assert.Equal(t, true, ExistsInEnv("foo"), "Returns true for existing environment variable")
	assert.Equal(t, false, ExistsInEnv("bar"), "Returns false for nonexistent environment variable")
}

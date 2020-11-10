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

package schema

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestNewAssetSchemaImpl tests if the Asset Schema Validator getter
func TestNewAssetSchemaImpl(t *testing.T) {
	assert.NotPanics(t, func() {
		NewAssetSchemaImpl()
	}, "Runs without panicking")
}

// TestAssetSchemaImpl_ValidateAsset tests if the Asset Schema Validator ValidateAsset works
func TestAssetSchemaImpl_ValidateAsset(t *testing.T) {
	assert.NotPanics(t, func() {
		os.Chdir("../..")
		defer os.Chdir("./src/schema")
		payload := make(map[string]interface{})
		NewAssetSchemaImpl().ValidateAsset(context.Background(), payload)
	}, "Runs without panicking")
}

// TestAssetSchemaImpl_ValidateAttachSubasset tests if the Asset Schema Validator ValidateAttachSubasset works
func TestAssetSchemaImpl_ValidateAttachSubasset(t *testing.T) {
	assert.NotPanics(t, func() {
		os.Chdir("../..")
		defer os.Chdir("./src/schema")
		payload := make(map[string]interface{})
		NewAssetSchemaImpl().ValidateAttachSubasset(context.Background(), payload)
	}, "Runs without panicking")
}

// TestAssetSchemaImpl_ValidateDetachSubasset tests if the Asset Schema Validator ValidateDetachSubasset works
func TestAssetSchemaImpl_ValidateDetachSubasset(t *testing.T) {
	assert.NotPanics(t, func() {
		os.Chdir("../..")
		defer os.Chdir("./src/schema")
		payload := make(map[string]interface{})
		NewAssetSchemaImpl().ValidateDetachSubasset(context.Background(), payload)
	}, "Runs without panicking")
}

// TestAssetSchemaImpl_ValidateTransferAsset tests if the Asset Schema Validator ValidateTransferAsset works
func TestAssetSchemaImpl_ValidateTransferAsset(t *testing.T) {
	assert.NotPanics(t, func() {
		os.Chdir("../..")
		defer os.Chdir("./src/schema")
		payload := make(map[string]interface{})
		NewAssetSchemaImpl().ValidateTransferAsset(context.Background(), payload)
	}, "Runs without panicking")
}


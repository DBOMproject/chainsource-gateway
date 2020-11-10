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
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const validPayloadPath = "../testdata/schema_tests/core/valid.json"
const invalidPayloadPath = "../testdata/schema_tests/core/invalid.json"
const nonexistentSchemaPath = "../testdata/schema_tests/core/missing.json"
const validSchemaPath = "../testdata/schema_tests/core/good_schema.json"
const invalidSchemaPath = "../testdata/schema_tests/core/bad_schema.json"

// loadJSONFileAsMap loads a JSON file from a given path as a map[string]interface{}
func loadJSONFileAsMap(path string) (asMap map[string]interface{}) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	json.NewDecoder(file).Decode(&asMap)
	return
}

// Test_coreValidator tests if the qri.io Schema Validator implementation works
func Test_coreValidator(t *testing.T) {
	t.Run("When_Valid_Schema_Valid_Payload", func(t *testing.T) {
		schemaErrors, isValid, err := coreValidator(context.Background(),
			loadJSONFileAsMap(validPayloadPath),
			validSchemaPath)
		assert.NoError(t, err, "Validator returned no errors")
		assert.Equal(t, true, isValid, "Payload marked as valid")
		assert.Len(t, schemaErrors, 0, "No schema errors were returned")
	})
	t.Run("When_Valid_Schema_Invalid_Payload", func(t *testing.T) {
		schemaErrors, isValid, err := coreValidator(context.Background(),
			loadJSONFileAsMap(invalidPayloadPath),
			validSchemaPath)
		assert.NoError(t, err, "Validator returned no errors")
		assert.Equal(t, false, isValid, "Payload marked as invalid")
		assert.NotEqual(t, "", schemaErrors, "Schema errors were returned")
	})
	t.Run("When_Invalid_Schema", func(t *testing.T) {
		assert.Panics(t, func() {
			coreValidator(context.Background(),
				loadJSONFileAsMap(validPayloadPath),
				invalidSchemaPath)
		}, "Validator panics")
	})
	t.Run("When_Nonexistent_Schema", func(t *testing.T) {
		assert.Panics(t, func() {
			coreValidator(context.Background(),
				loadJSONFileAsMap(validPayloadPath),
				nonexistentSchemaPath)
		}, "Validator panics")
	})


}

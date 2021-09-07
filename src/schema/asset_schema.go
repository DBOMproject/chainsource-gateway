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
	"chainsource-gateway/helpers"
	"context"
)

type AssetSchema interface {
	ValidateAsset(ctx context.Context, json map[string]interface{}) (schemaErrors string, isValid bool, err error)
	ValidateAttachSubasset(ctx context.Context, json map[string]interface{}) (schemaErrors string, isValid bool, err error)
	ValidateDetachSubasset(ctx context.Context, json map[string]interface{}) (schemaErrors string, isValid bool, err error)
	ValidateQueryAsset(ctx context.Context, json map[string]interface{}) (schemaErrors string, isValid bool, err error)
	ValidateTransferAsset(ctx context.Context, json map[string]interface{}) (schemaErrors string, isValid bool, err error)
}

// AssetSchemaImpl implements the AssetSchema interface using qri-io/jsonschema
type AssetSchemaImpl struct {
	assetSchemaLocation          string
	attachSubassetSchemaLocation string
	detachSubassetSchemaLocation string
	querySchemaLocation          string
	transferAssetSchemaLocation  string
}

func NewAssetSchemaImpl() AssetSchemaImpl {
	return AssetSchemaImpl{
		assetSchemaLocation:          "./src/rest_schema/asset.json",
		attachSubassetSchemaLocation: "./src/rest_schema/attach-subasset.json",
		detachSubassetSchemaLocation: "./src/rest_schema/detach-subasset.json",
		querySchemaLocation:          "./src/rest_schema/query-asset.json",
		transferAssetSchemaLocation:  "./src/rest_schema/transfer-asset.json",
	}
}

var schemaLog = helpers.GetLogger("SchemaController")

// ValidateAsset is a wrapper function to validate Asset Object
func (a AssetSchemaImpl) ValidateAsset(ctx context.Context, json map[string]interface{}) (schemaErrors string, isValid bool, err error) {
	schemaErrors, isValid, err = coreValidator(ctx, json, a.assetSchemaLocation)
	return
}

// validateAttachSubasset is a wrapper function to validate schema for the Attach API
func (a AssetSchemaImpl) ValidateAttachSubasset(ctx context.Context, json map[string]interface{}) (schemaErrors string, isValid bool, err error) {
	schemaErrors, isValid, err = coreValidator(ctx, json, a.attachSubassetSchemaLocation)
	return
}

// validateDetachSubasset is a wrapper function to validate schema for the Detach API
func (a AssetSchemaImpl) ValidateDetachSubasset(ctx context.Context, json map[string]interface{}) (schemaErrors string, isValid bool, err error) {
	schemaErrors, isValid, err = coreValidator(ctx, json, a.detachSubassetSchemaLocation)
	return
}

// ValidateQueryAsset is a wrapper function to validate schema for the query API
func (a AssetSchemaImpl) ValidateQueryAsset(ctx context.Context, json map[string]interface{}) (schemaErrors string, isValid bool, err error) {
	schemaErrors, isValid, err = coreValidator(ctx, json, a.querySchemaLocation)
	return
}

// validateTransferAsset is a wrapper function to validate schema for the Transfer API
func (a AssetSchemaImpl) ValidateTransferAsset(ctx context.Context, json map[string]interface{}) (schemaErrors string, isValid bool, err error) {
	schemaErrors, isValid, err = coreValidator(ctx, json, a.transferAssetSchemaLocation)
	return
}

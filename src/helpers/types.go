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

// AssetRoutingVars is a type to hold the asset context derived from the parametrized url
type AssetRoutingVars struct {
	RepoID    string
	ChannelID string
	AssetID   string
}

// AssetRoutingVars is a type to hold the asset context derived from the parametrized url
type AssetQueryVars struct {
	Query  string
	Filter string
	Skip   string
	Limit  string
}

// ExportRoutingVars is a type to hold the asset context derived from the parametrized url
type ExportRoutingVars struct {
	FileName       string
	InlineResponse string
}

// Asset is a type representing an asset recognized by the gateway
// If the JSON Schema is changed, it must also be updated here otherwise the unmarshalled JSON will not contain the changes
type Asset struct {
	StandardVersion       float64                `json:"standardVersion"`
	DocumentName          string                 `json:"documentName"`
	DocumentCreator       string                 `json:"documentCreator"`
	DocumentCreatedDate   string                 `json:"documentCreatedDate"`
	AssetType             string                 `json:"assetType"`
	AssetSubType          string                 `json:"assetSubType"`
	AssetManufacturer     string                 `json:"assetManufacturer"`
	AssetModelNumber      string                 `json:"assetModelNumber"`
	AssetDescription      string                 `json:"assetDescription"`
	AssetMetadata         interface{}            `json:"assetMetadata"`
	AttachedChildren      []AssetLinkElement     `json:"attachedChildren,omitempty"`
	ParentAsset           *AssetLinkElement      `json:"parentAsset,omitempty"`
	CustodyTransferEvents []CustodyTransferEvent `json:"custodyTransferEvents ,omitempty"`
	ManufactureSignature  string                 `json:"manufactureSignature"`
	Children              map[string]*Asset      `json:"children,omitempty"`
	Parent                map[string]*Asset      `json:"parent,omitempty"`
	ReadOnly              bool                   `json:"readOnly,omitempty"`
}

// AssetElement is a type representing a link element (parent or child)
type AssetElement struct {
	RepoID    string `json:"repoID"`
	ChannelID string `json:"channelID"`
	AssetID   string `json:"assetID"`
	Asset     *Asset `json:"Asset,omitempty"`
}

// AssetLinkElement is a type representing a link element (parent or child)
type AssetLinkElement struct {
	AssetElement
	Role    string `json:"role"`
	SubRole string `json:"subRole"`
}

// AssetTransferElement is a type representing a link element (parent or child)
type AssetTransferElement struct {
	AssetElement
	TransferDescription string `json:"transferDescription"`
}

// CustodyTransferEvent is a type representing a custody transfer event
type CustodyTransferEvent struct {
	Timestamp            string `json:"timestamp,omitempty"`
	TransferDescription  string `json:"transferDescription,omitempty"`
	SourceRepoID         string `json:"sourceRepoID"`
	SourceChannelID      string `json:"sourceChannelID"`
	SourceAssetID        string `json:"sourceAssetID"`
	DestinationRepoID    string `json:"destinationRepoID"`
	DestinationChannelID string `json:"destinationChannelID"`
	DestinationAssetID   string `json:"destinationAssetID"`
}

// Fingerprint is a type representing a manufacture fingerprint
type Fingerprint struct {
	ManufactureFingerprint string `json:"manufactureFingerprint"`
}

// ValidateBody is a type representing the request body expected by the pgp service for validating a signature
type ValidateBody struct {
	Fingerprint string             `json:"fingerprint"`
	Signature   string             `json:"signature"`
	Input       AssetNoChildParent `json:"input"`
}

// AssetNoChildParent is a type representing an asset recognized by the gateway
// If the JSON Schema is changed, it must also be updated here otherwise the unmarshalled JSON will not contain the changes
// This is the same as Asset, except it excludes the signature field, and the parent and children links from the json output for validation of a signature
type AssetNoChildParent struct {
	StandardVersion       float64                `json:"standardVersion"`
	DocumentName          string                 `json:"documentName"`
	DocumentCreator       string                 `json:"documentCreator"`
	DocumentCreatedDate   string                 `json:"documentCreatedDate"`
	AssetType             string                 `json:"assetType"`
	AssetSubType          string                 `json:"assetSubType"`
	AssetManufacturer     string                 `json:"assetManufacturer"`
	AssetModelNumber      string                 `json:"assetModelNumber"`
	AssetDescription      string                 `json:"assetDescription"`
	AssetMetadata         interface{}            `json:"assetMetadata"`
	AttachedChildren      []AssetLinkElement     `json:"-"`
	ParentAsset           *AssetLinkElement      `json:"-"`
	CustodyTransferEvents []CustodyTransferEvent `json:"-"`
	ManufactureSignature  string                 `json:"-"`
	Children              map[string]*Asset      `json:"-"`
	Parent                map[string]*Asset      `json:"-"`
	ReadOnly              bool                   `json:"-"`
}

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

import "github.com/nats-io/nats.go"

type AssetRoutingVars struct {
	ChannelID string `json:"channelId"`
	AssetID   string `json:"assetId"`
}

type AssetQueryRoutingVars struct {
	ChannelID string      `json:"channelId"`
	Query     interface{} `json:"query"`
}

type AssetRichQueryRoutingVars struct {
	ChannelID string `json:"channelId"`
	Query     string `json:"query"`
	Fields    string `json:"fields"`
	Limit     string `json:"limit"`
	Skip      string `json:"skip"`
}

type ChannelRoutingVars struct {
	ChannelID string `json:"channelId"`
	NotaryID  string `json:"notaryId"`
}

type FederationRoutingVars struct {
	RequestID string `json:"requestId"`
}

type NotarizationsDefinition struct {
	NotaryId   string      `json:"notaryId"`
	NotaryMeta interface{} `json:"notaryMeta"`
}

type LinksDefinition struct {
	AssetURI string `json:"assetUri"`
	Type     string `json:"type"`
	Comment  string `json:"comment"`
	Id       string `json:"id"`
}

type SignMetaDefinition struct {
	Authority string `json:"authority"`
	KeyId     string `json:"keyId"`
	Sign      string `json:"sign"`
}
type SignatureDefinition struct {
	HashType string             `json:"hashType"`
	SignType string             `json:"signType"`
	SignMeta SignMetaDefinition `json:"signMeta"`
}

type Asset struct {
	StandardVersion float64                   `json:"standardVersion"`
	SchemaUrl       string                    `json:"schemaUrl"`
	CreatedAt       string                    `json:"createdAt"`
	ModifiedAt      string                    `json:"modifiedAt"`
	Notarizations   []NotarizationsDefinition `json:"notarizations"`
	Links           []LinksDefinition         `json:"links"`
	Signatures      []SignatureDefinition     `json:"signatures"`
	Body            interface{}               `json:"body"`
	Id              string                    `json:"id,omitempty"`
}

type NotaryMeta struct {
	Id     string      `json:"id"`
	Type   string      `json:"type"`
	Config interface{} `json:"config"`
}

type Channel struct {
	ChannelId   string       `json:"channelId"`
	Description string       `json:"description"`
	Type        string       `json:"type"`
	Notaries    []NotaryMeta `json:"notaries"`
	CreatedAt   string       `json:"createdAt"`
	ModifiedAt  string       `json:"modifiedAt"`
	Id          string       `json:"id,omitempty"`
}

type NodeMeta struct {
	NodeID   string `json:"nodeId"`
	KnownURI string `json:"knownUri"`
	Status   string `json:"status"`
	Access   string `json:"access"`
}

type ChannelMeta struct {
	NodeID    string `json:"nodeId"`
	ChannelID string `json:"channelId"`
	Status    string `json:"status"`
	Access    string `json:"access"`
}

type ChannelConnection struct {
	ChannelId string `json:"channelId"`
	Status    string `json:"status"`
	Access    string `json:"access"`
}

type NodeConnection struct {
	NodeId             string              `json:"nodeId"`
	Status             string              `json:"status"`
	ChannelConnections []ChannelConnection `json:"channelConnections"`
}

type Node struct {
	Id              string           `json:"id"`
	NodeID          string           `json:"nodeId"`
	PublicKeys      []string         `json:"publicKeys"`
	NodeConnections []NodeConnection `json:"nodeConnections"`
	CreatedAt       string           `json:"createdAt"`
	ModifiedAt      string           `json:"modifiedAt"`
}

type CreateChannelMeta struct {
	ChannelMeta interface{} `json:"channelMeta"`
}

type ListOneChannelMeta struct {
	ChannelID string `json:"channelId"`
}

type FederationRequestOperations struct {
	Type      string `json:"type"`
	NodeURI   string `json:"nodeUri,omitempty"`
	NodeID    string `json:"nodeId,omitempty"`
	ChannelID string `json:"channelId,omitempty"`
}

type FederationRequestBody struct {
	RequestID string `json:"requestId"`
	Type      string `json:"type"`
	NodeID    string `json:"nodeId,omitempty"`
	ChannelID string `json:"channelId,omitempty"`
}

type UpdateChannelRequestMeta struct {
	NotaryId string `json:"notaryId"`
}

type UpdateChannelMeta struct {
	ChannelID string      `json:"channelId"`
	NotaryID  interface{} `json:"notaryId"`
}

type AssetMeta struct {
	ChannelID string `json:"channelId"`
	AssetID   string `json:"assetId"`
	Payload   Asset  `json:"payload"`
}

type LinkAssetMeta struct {
	ChannelID string          `json:"channelId"`
	AssetID   string          `json:"assetId"`
	Payload   LinksDefinition `json:"payload"`
}

type UnlinkAssetMeta struct {
	ChannelID string `json:"channelId"`
	AssetID   string `json:"assetId"`
	LinkID    string `json:"linkId"`
}

type PrismaQuery struct {
	AssetID struct {
		Equals string `json:"equals"`
	} `json:"assetId"`
	ChannelID struct {
		Equals string `json:"equals,omitempty"`
	} `json:"channelId,omitempty"`
}

type QueryMeta struct {
	Take    *int         `json:"take,omitempty"`
	Skip    *int         `json:"skip,omitempty"`
	Cursor  *interface{} `json:"cursor,omitempty"`
	Where   PrismaQuery  `json:"where,omitempty"`
	OrderBy interface{}  `json:"orderBy,omitempty"`
}

type FederationMeta struct {
	RequestID  string `json:"requestId"`
	NodeURI    string `json:"nodeUri,omitempty"`
	NodeID     string `json:"nodeId,omitempty"`
	ChannelID  string `json:"channelId,omitempty"`
	Status     string `json:"status"`
	CreatedAt  string `json:"createdAt"`
	ModifiedAt string `json:"modifiedAt"`
}

type HistoryMeta struct {
	Id        string      `json:"id"`
	ChannelID string      `json:"channelId"`
	AssetID   string      `json:"assetId"`
	Action    string      `json:"action"`
	Payload   interface{} `json:"payload"`
	Timestamp string      `json:"timestamp"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Server struct {
	nc *nats.Conn
}

type NodeResultResponse struct {
	Success bool   `json:"success"`
	Status  string `json:"status,omitempty"`
	Result  []Node `json:"result,omitempty"`
}

type FederationResultResponse struct {
	Success bool             `json:"success"`
	Status  string           `json:"status,omitempty"`
	Result  []FederationMeta `json:"result,omitempty"`
}

type AssetResultResponse struct {
	Success bool        `json:"success"`
	Status  string      `json:"status,omitempty"`
	Result  []AssetMeta `json:"result,omitempty"`
}

type ChannelResultResponse struct {
	Success bool      `json:"success"`
	Status  string    `json:"status,omitempty"`
	Result  []Channel `json:"result,omitempty"`
}

type HistoryResultResponse struct {
	Success bool          `json:"success"`
	Status  string        `json:"status,omitempty"`
	Result  []HistoryMeta `json:"result,omitempty"`
}

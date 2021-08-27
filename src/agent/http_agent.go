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

// Package agent contains functions and types to encapsulate the gateway's interactions with the agent
package agent

import (
	"chainsource-gateway/helpers"
	"chainsource-gateway/tracing"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strconv"

	"github.com/spf13/viper"

	"github.com/opentracing/opentracing-go"
)

// HttpAgent is a type representing an agent that works over a RESTFUL HTTP interface
type HttpAgent struct {
	Config   *Config
	AgentURL string
}

type HttpAgentProvider struct{}

// CommitArgs is a type representing the arguments sent to the agent interfaces "commit" function
type CommitArgs struct {
	ChannelID  string
	AssetID    string
	CommitType string
	Payload    helpers.Asset
}

// QueryArgs is a type representing the arguments sent to the agent interfaces "query" function
type QueryArgs struct {
	ChannelID string
	AssetID   string
}

// RichQueryArgs is a type representing the arguments sent to the agent interfaces "rich query" function
type RichQueryArgs struct {
	Query  interface{} `json:"query"`
	Filter []string    `json:"filter"`
	Skip   int         `json:"skip"`
	Limit  int         `json:"limit"`
}

// CommitBody is a type representing the request body expected by the agent on it's commit interface
type CommitBody struct {
	RecordID        string        `json:"recordID"`
	RecordIDPayload helpers.Asset `json:"recordIDPayload"`
}

func NewHTTPAgentProvider() (p Provider) {
	p = &HttpAgentProvider{}
	return
}

// NewAgent creates a new agent, given an agent config element
func (provider *HttpAgentProvider) NewAgent(config *Config) (a Agent) {
	a = HttpAgent{
		Config:   config,
		AgentURL: "http://" + config.Host + ":" + strconv.Itoa(int(config.Port)),
	}
	return
}

func (provider *HttpAgentProvider) GetAgentConfigForRepo(repoID string) (Config, error) {
	agents := make(map[string]Config)
	err := viper.UnmarshalKey("agents", &agents)
	if err != nil {
		return Config{}, err
	}
	val, exists := agents[repoID]
	if !exists {
		err = errors.New("No agent configured for sent repoID")
		log.Err(err).Str("repoID", repoID)
		return Config{}, err
	}
	return val, err
}

// Commit performs a commit on the agent
func (a HttpAgent) Commit(ctx context.Context, args CommitArgs) (result map[string]interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Commit on agent")
	defer span.Finish()
	url := a.AgentURL

	extraHeaders := make(map[string]string)
	extraHeaders["Commit-Type"] = args.CommitType
	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.TextMapCarrier(extraHeaders))

	span.SetTag("agent-url", url)
	commitPath := "/channels/" + args.ChannelID + "/records/"

	body := CommitBody{
		RecordID:        args.AssetID,
		RecordIDPayload: args.Payload,
	}
	//helpers.DebugPrintMap("COMMIT", body)
	bytesRepresentation, err := json.Marshal(body)

	result, err = helpers.PostJSONRequest(url, commitPath, extraHeaders, bytesRepresentation)
	if err != nil {
		tracing.LogAndTraceErr(log, span, err, "Commit failed on agent")
	}

	return
}

// QueryAssets performs a rich query on the agent
func (a HttpAgent) QueryAssets(ctx context.Context, args QueryArgs, body RichQueryArgs) (result map[string]interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Query Assets")
	defer span.Finish()
	url := a.AgentURL

	extraHeaders := make(map[string]string)
	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.TextMapCarrier(extraHeaders))

	span.SetTag("agent-url", url)
	queryPath := "/channels/" + args.ChannelID + "/records/_query"
	bytesRepresentation, err := json.Marshal(body)

	result, err = helpers.PostJSONRequest(url, queryPath, extraHeaders, bytesRepresentation)
	if err != nil {
		tracing.LogAndTraceErr(log, span, err, "Query failed on agent")
	}

	return
}

// ListChannels performs a listChannels on the agent. Returns it as a io stream
func (a HttpAgent) ListChannels(ctx context.Context, args QueryArgs) (resultStream io.ReadCloser, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "List Channels")
	defer span.Finish()

	extraHeaders := make(map[string]string)
	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.TextMapCarrier(extraHeaders))

	url := a.AgentURL
	queryPath := "/channels"
	resultStream, err = helpers.GetRequest(url, queryPath, extraHeaders)
	if err != nil {
		tracing.LogAndTraceErr(log, span, err, "Failed to list channels from agent")
	}
	return
}

// ListAssets performs a listAssets on the agent. Returns it as a io stream
func (a HttpAgent) ListAssets(ctx context.Context, args QueryArgs) (resultStream io.ReadCloser, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "List Assets")
	defer span.Finish()

	extraHeaders := make(map[string]string)
	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.TextMapCarrier(extraHeaders))

	url := a.AgentURL
	queryPath := "/channels/" + args.ChannelID + "/records"
	resultStream, err = helpers.GetRequest(url, queryPath, extraHeaders)
	if err != nil {
		tracing.LogAndTraceErr(log, span, err, "Failed to list assets from agent")
	}
	return
}

// QueryStream performs a query on the agent. Returns it as a io stream
func (a HttpAgent) QueryStream(ctx context.Context, args QueryArgs) (resultStream io.ReadCloser, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Query on HttpAgent")
	defer span.Finish()

	extraHeaders := make(map[string]string)
	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.TextMapCarrier(extraHeaders))

	url := a.AgentURL
	queryPath := "/channels/" + args.ChannelID + "/records/" + args.AssetID
	resultStream, err = helpers.GetRequest(url, queryPath, extraHeaders)
	if err != nil {
		tracing.LogAndTraceErr(log, span, err, "Failed to query From agent")
	}
	return
}

// QueryAuditTrail performs an audit on the agent. Returns it as a map
func (a HttpAgent) QueryAuditTrail(ctx context.Context, args QueryArgs) (result map[string]interface{}, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Audit on HttpAgent")
	defer span.Finish()

	extraHeaders := make(map[string]string)
	opentracing.GlobalTracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.TextMapCarrier(extraHeaders))

	url := a.AgentURL
	queryPath := "/channels/" + args.ChannelID + "/records/" + args.AssetID + "/audit"
	resultStream, err := helpers.GetRequest(url, queryPath, extraHeaders)
	if err != nil {
		tracing.LogAndTraceErr(log, span, err, "Could not retrieve audit trail from agent")
		return nil, err
	}
	err = json.NewDecoder(resultStream).Decode(&result)
	if err != nil {
		tracing.LogAndTraceErr(log, span, err, "Audit trail is not valid JSON")
		return
	}
	return

}

func (a HttpAgent) GetHost() string {
	return a.Config.Host
}

func (a HttpAgent) GetPort() int {
	return int(a.Config.Port)
}

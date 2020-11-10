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

package agent

import (
	"context"
	"io"
)

// Agent is an interface that is expected to be implemented by an agent
type Agent interface {
	Commit(ctx context.Context, args CommitArgs) (result map[string]interface{}, err error)
	QueryStream(ctx context.Context, args QueryArgs) (resultStream io.ReadCloser, err error)
	QueryAuditTrail(ctx context.Context, args QueryArgs) (result map[string]interface{}, err error)
	GetHost() string
	GetPort() int
}

type Provider interface {
	NewAgent(config *Config) (agent Agent)
	GetAgentConfigForRepo(repoID string) (Config, error)
}

// Config is a type representing the agent-config.yaml file
type Config struct {
	Version int64
	Host    string
	Port    int64
	RepoID  string
	Enabled bool
}




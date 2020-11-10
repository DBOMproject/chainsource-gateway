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

package mocks

import (
	"chainsource-gateway/agent"
	"github.com/golang/mock/gomock"
)

type newAgentRequestMatcher struct {
	repoID string
}

func (a *newAgentRequestMatcher) Matches(args interface{}) bool {
	rArgs := args.(*agent.Config)
	valid := a.repoID == rArgs.RepoID
	return valid
}

func (a *newAgentRequestMatcher) String() string {
	return "is an agent request for " + a.repoID
}

func AgentRequestFor(repoID string) gomock.Matcher {
	return &newAgentRequestMatcher{repoID: repoID}
}

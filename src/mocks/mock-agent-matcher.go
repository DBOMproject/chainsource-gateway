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

type queryArgsMatcher struct {
	channelID string
	assetID   string
}

func (q *queryArgsMatcher) Matches(args interface{}) bool {
	qArgs := args.(agent.QueryArgs)
	valid := (q.channelID == qArgs.ChannelID) && (q.assetID == qArgs.AssetID)
	return valid
}

func (q *queryArgsMatcher) String() string {
	return "is a query for " + q.channelID + "/" + q.assetID
}

func AgentQueryFor(channelID string, assetID string) gomock.Matcher {
	return &queryArgsMatcher{channelID: channelID, assetID: assetID}
}


type commitArgsMatcher struct {
	channelID string
	assetID   string
	commitType string
}

func (c *commitArgsMatcher) Matches(args interface{}) bool {
	cArgs := args.(agent.CommitArgs)
	valid := (c.channelID == cArgs.ChannelID) && (c.assetID == cArgs.AssetID) && (c.commitType == cArgs.CommitType)
	return valid
}

func (c *commitArgsMatcher) String() string {
	return "is a commit to " + c.channelID + "/" + c.assetID
}

func AgentCommitTo(channelID string, assetID string, commitType string) gomock.Matcher {
	return &commitArgsMatcher{channelID: channelID, assetID: assetID, commitType:commitType}
}

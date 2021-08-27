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
	"chainsource-gateway/helpers"
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

// setupOKMockRemoteHttpAgent sets up a mock HTTP agent that always returns OK
func setupOKMockRemoteHttpAgent(channelID string, recordID string) {
	gock.New("http://mock-agent").
		Get("/channels").
		Reply(200).
		JSON([]string{"test", "test2"})
	gock.New("http://mock-agent").
		Get("/channels/" + channelID + "/records").
		Reply(200).
		JSON(map[string]bool{"success": true})
	gock.New("http://mock-agent").
		Post("/channels/" + channelID + "/records").
		Reply(200).
		JSON(map[string]bool{"success": true})
	gock.New("http://mock-agent").
		Get("/channels/" + channelID + "/records/" + recordID).
		Reply(200).
		JSON(map[string]string{"foo": "bar"})
	gock.New("http://mock-agent").
		Get("/channels/" + channelID + "/records/" + recordID + "/audit").
		Reply(200).
		JSON(map[string]string{"foo": "test"})
	gock.New("http://mock-agent").
		Get("/channels/" + channelID + "/records/_query").
		Reply(200).
		JSON(map[string]string{"foo": "bar"})
}

// setupFailMockRemoteHttpAgent sets up a mock HTTP agent that always returns InternalServerError
func setupFailMockRemoteHttpAgent(channelID string, recordID string) {
	gock.New("http://mock-agent").
		Get("/channels").
		Reply(500)
	gock.New("http://mock-agent").
		Get("/channels/" + channelID + "/records").
		Reply(500)
	gock.New("http://mock-agent").
		Post("/channels/" + channelID + "/records").
		Reply(500)
	gock.New("http://mock-agent").
		Get("/channels/" + channelID + "/records/" + recordID).
		Reply(500)
	gock.New("http://mock-agent").
		Get("/channels/" + channelID + "/records/" + recordID + "/audit").
		Reply(500)
	gock.New("http://mock-agent").
		Post("/channels/" + channelID + "/records/_query").
		Reply(500)
}

// setupInvalidJSONMockRemoteHttpAgent sets up a mock HTTP agent that always returns Bad JSON
func setupInvalidJSONMockRemoteHttpAgent(channelID string, recordID string) {
	gock.New("http://mock-agent").
		Get("/channels/" + channelID + "/records/" + recordID).
		Reply(200).
		BodyString("invalid.JSON")
}

// getMockHTTPAgent gets a fake HTTP agent at mock-agent
func getMockHTTPAgent() HttpAgent {
	return HttpAgent{
		Config: &Config{
			Version: 1,
			Host:    "mock-agent",
			Port:    80,
			RepoID:  "T1",
			Enabled: true,
		},
		AgentURL: "http://mock-agent",
	}
}

// TestHttpAgentProvider_GetAgentConfigForRepo tests the agent provider's get config method
func TestHttpAgentProvider_GetAgentConfigForRepo(t *testing.T) {

	t.Run("With_Agent_That_Exists", func(t *testing.T) {
		provider := NewHTTPAgentProvider()
		agents := make(map[string]Config)
		agents["T1"] = Config{
			Version: 1,
			Host:    "mock-agent",
			Port:    80,
			RepoID:  "T1",
			Enabled: true,
		}
		viper.Set("agents", agents)
		defer viper.Set("agents", nil)
		config, err := provider.GetAgentConfigForRepo("T1")
		assert.NoError(t, err, "No error occurred while retrieving config")
		assert.Equal(t, agents["T1"], config, "The correct agent config was returned")
	})

	t.Run("With_Agent_That_Does_Not_Exist", func(t *testing.T) {
		provider := NewHTTPAgentProvider()
		_, err := provider.GetAgentConfigForRepo("T2")
		assert.Error(t, err, "An error was returned")
	})

}

// TestHttpAgentProvider_NewAgent tests the agent provider's new agent method
func TestHttpAgentProvider_NewAgent(t *testing.T) {
	t.Run("With_Valid_Config", func(t *testing.T) {
		provider := NewHTTPAgentProvider()
		config := Config{
			Version: 1,
			Host:    "mock-agent",
			Port:    80,
			RepoID:  "T1",
			Enabled: true,
		}
		agent := provider.NewAgent(&config)
		assert.Implements(t, (*Agent)(nil), agent, "Implements the agent interface")
	})
}

// TestHttpAgentProvider_Commit tests the agent's commit method
func TestHttpAgent_Commit(t *testing.T) {
	agent := getMockHTTPAgent()
	t.Run("With_OK_Agent", func(t *testing.T) {
		setupOKMockRemoteHttpAgent("C1", "A1")
		defer gock.Off()
		_, err := agent.Commit(context.Background(), CommitArgs{
			ChannelID:  "C1",
			AssetID:    "A1",
			CommitType: "COMMIT_TYPE_EXAMPLE",
			Payload:    helpers.Asset{},
		})
		assert.NoError(t, err, "Completes commit successfully")
	})
	t.Run("With_Fail_Agent", func(t *testing.T) {
		setupFailMockRemoteHttpAgent("C1", "A1")
		defer gock.Off()
		_, err := agent.Commit(context.Background(), CommitArgs{
			ChannelID:  "C1",
			AssetID:    "A1",
			CommitType: "COMMIT_TYPE_EXAMPLE",
			Payload:    helpers.Asset{},
		})
		assert.Error(t, err, "Commit fails")
	})
}

// TestHttpAgent_QueryStream tests the agent's query method
func TestHttpAgent_QueryStream(t *testing.T) {
	agent := getMockHTTPAgent()
	t.Run("With_OK_Agent", func(t *testing.T) {
		setupOKMockRemoteHttpAgent("C1", "A1")
		defer gock.Off()
		_, err := agent.QueryStream(context.Background(), QueryArgs{
			ChannelID: "C1",
			AssetID:   "A1",
		})
		assert.NoError(t, err, "Completes query successfully")
	})
	t.Run("With_Fail_Agent", func(t *testing.T) {
		setupFailMockRemoteHttpAgent("C1", "A1")
		defer gock.Off()
		_, err := agent.QueryStream(context.Background(), QueryArgs{
			ChannelID: "C1",
			AssetID:   "A1",
		})
		assert.Error(t, err, "Query fails")
	})
}

// TestHttpAgent_QueryAssets tests the agent's query method
func TestHttpAgent_QueryAssets(t *testing.T) {
	agent := getMockHTTPAgent()
	t.Run("With_OK_Agent", func(t *testing.T) {
		setupOKMockRemoteHttpAgent("C1", "A1")
		defer gock.Off()
		_, err := agent.QueryAssets(context.Background(), QueryArgs{
			ChannelID: "C1",
			AssetID:   "A1",
		}, RichQueryArgs{Query: nil, Filter: nil, Limit: 10, Skip: 0})
		assert.NoError(t, err, "Completes query successfully")
	})
	t.Run("With_Fail_Agent", func(t *testing.T) {
		setupFailMockRemoteHttpAgent("C1", "A1")
		defer gock.Off()
		_, err := agent.QueryAssets(context.Background(), QueryArgs{
			ChannelID: "C1",
			AssetID:   "A1",
		}, RichQueryArgs{Query: nil, Filter: nil, Limit: 10, Skip: 0})
		assert.Error(t, err, "Query fails")
	})
}

// TestHttpAgent_ListAssets tests the agent's ListAssets method
func TestHttpAgent_ListAssets(t *testing.T) {
	agent := getMockHTTPAgent()
	t.Run("With_OK_Agent", func(t *testing.T) {
		setupOKMockRemoteHttpAgent("C1", "A1")
		defer gock.Off()
		_, err := agent.ListAssets(context.Background(), QueryArgs{
			ChannelID: "C1",
		})
		assert.NoError(t, err, "Completes list assets successfully")
	})
	t.Run("With_Fail_Agent", func(t *testing.T) {
		setupFailMockRemoteHttpAgent("C1", "A1")
		defer gock.Off()
		_, err := agent.ListAssets(context.Background(), QueryArgs{
			ChannelID: "C1",
		})
		assert.Error(t, err, "Query fails")
	})
}

// TestHttpAgent_ListChannels tests the agent's ListChannels method
func TestHttpAgent_ListChannels(t *testing.T) {
	agent := getMockHTTPAgent()
	t.Run("With_OK_Agent", func(t *testing.T) {
		setupOKMockRemoteHttpAgent("", "")
		defer gock.Off()
		_, err := agent.ListChannels(context.Background(), QueryArgs{})
		assert.NoError(t, err, "Completes list channels successfully")
	})
	t.Run("With_Fail_Agent", func(t *testing.T) {
		setupFailMockRemoteHttpAgent("C1", "A1")
		defer gock.Off()
		_, err := agent.ListChannels(context.Background(), QueryArgs{})
		assert.Error(t, err, "Query fails")
	})
}

// TestHttpAgent_QueryAuditTrail tests the agent's query audit trail method
func TestHttpAgent_QueryAuditTrail(t *testing.T) {
	agent := getMockHTTPAgent()
	t.Run("With_OK_Agent", func(t *testing.T) {
		setupOKMockRemoteHttpAgent("C1", "A1")
		defer gock.Off()
		_, err := agent.QueryAuditTrail(context.Background(), QueryArgs{
			ChannelID: "C1",
			AssetID:   "A1",
		})
		assert.NoError(t, err, "Completes query audit trail successfully")
	})
	t.Run("With_Fail_Agent", func(t *testing.T) {
		setupFailMockRemoteHttpAgent("C1", "A1")
		defer gock.Off()
		_, err := agent.QueryAuditTrail(context.Background(), QueryArgs{
			ChannelID: "C1",
			AssetID:   "A1",
		})
		assert.Error(t, err, "Query audit trail fails")
	})
	t.Run("With_Invalid_JSON_From_Agent", func(t *testing.T) {
		setupInvalidJSONMockRemoteHttpAgent("C1", "A1")
		defer gock.Off()
		_, err := agent.QueryAuditTrail(context.Background(), QueryArgs{
			ChannelID: "C1",
			AssetID:   "A1",
		})
		assert.Error(t, err, "Query audit trail fails due to unmarshaling")
	})
}

// TestHttpAgent_GetHost tests the agent's get host method
func TestHttpAgent_GetHost(t *testing.T) {
	agent := getMockHTTPAgent()
	assert.Equal(t, "mock-agent", agent.GetHost(), "Correct host is returned")
}

// TestHttpAgent_GetPort tests the agent's get port method
func TestHttpAgent_GetPort(t *testing.T) {
	agent := getMockHTTPAgent()
	assert.Equal(t, 80, agent.GetPort(), "Correct port is returned")
}

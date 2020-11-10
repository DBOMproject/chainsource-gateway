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
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

// Test_logAgentConfig contains tests for the agent config logger
func Test_logAgentConfig(t *testing.T) {
	agents := make(map[string]Config)
	agents["T1"] = Config{
		Version: 1,
		Host:    "mock-agent",
		Port:    80,
		RepoID:  "T1",
		Enabled: true,
	}
	assert.NotPanics(t, func() { logAgentConfig(agents) }, "Agent config logger does not panic")
}

// TestGetAgentConfig contains tests for the agent config getter
func TestGetAgentConfig(t *testing.T) {
	t.Run("When_HappyPath", func(t *testing.T) {
		os.Chdir("../..")
		defer os.Chdir("./src/agent")
		defer viper.Reset()
		assert.NotPanics(t, func() {
			GetAgentConfig()
		}, "Does not panic on trying to retrieve config")
		agents := make(map[string]Config)
		assert.NoError(t, viper.UnmarshalKey("agents", &agents))

		// Hot Reload Test
		currentTime := time.Now().Local()
		os.Chtimes("./config/agent-config.yaml", currentTime, currentTime)
	})

	t.Run("When_NotFound", func(t *testing.T) {
		assert.Panics(t, func() {
			GetAgentConfig()
		}, "Panics appropriately")
	})
}

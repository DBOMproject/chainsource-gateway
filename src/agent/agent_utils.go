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
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var log = helpers.GetLogger("AgentService")

// GetAgentConfig uses viper to get and store the agent config from disk.
// Automatically reloads configuration on any change
func GetAgentConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("agent-config")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("agent-config.yaml could not be loaded: %s \n", err))
	}

	agents := make(map[string]Config)
	_ = viper.UnmarshalKey("agents", &agents)
	logAgentConfig(agents)

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Info().Msgf("Hot Reload: AgentConfig has changed")
		agents = make(map[string]Config)
		_ = viper.UnmarshalKey("agents", &agents)
		logAgentConfig(agents)
	})
}

// Logs the agents loaded from the agent-config fie
func logAgentConfig(agentMap map[string]Config) {
	keys := make([]string, 0, len(agentMap))
	for key := range agentMap {
		keys = append(keys, key)
	}
	log.Info().Msgf("Agents Loaded: %s", strings.Join(keys, ","))
}

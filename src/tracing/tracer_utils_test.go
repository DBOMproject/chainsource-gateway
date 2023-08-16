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

package tracing

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_getSamplerConfigFromEnv tests if the SamplerConfig can be successfully gotten from the env
func Test_getSamplerConfigFromEnv(t *testing.T) {
	t.Run("When_Env_Vars_Set", func(t *testing.T) {
		os.Setenv(jaegerSamplerTypeVar, "const")
		os.Setenv(jaegerSamplerParamVar, "0.2")
		defer os.Unsetenv(jaegerSamplerParamVar)
		defer os.Unsetenv(jaegerSamplerTypeVar)

		_, err := getSamplerConfigFromEnv()
		assert.NoError(t, err, "Retrieves without an error")
	})
	t.Run("When_Env_Vars_Set_Unparseable", func(t *testing.T) {
		os.Setenv(jaegerSamplerTypeVar, "const")
		os.Setenv(jaegerSamplerParamVar, "UNKNOWN")
		defer os.Unsetenv(jaegerSamplerParamVar)
		defer os.Unsetenv(jaegerSamplerTypeVar)

		_, err := getSamplerConfigFromEnv()
		assert.Error(t, err, "Error Raised")
	})
	t.Run("When_Env_Vars_Not_Set", func(t *testing.T) {
		_, err := getSamplerConfigFromEnv()
		assert.NoError(t, err, "Retrieves without an error")
	})
}

// Test_cfgFromEnv tests if the Jaeger Config can be successfully gotten from the env
func Test_cfgFromEnv(t *testing.T) {
	t.Run("When_Empty_Env", func(t *testing.T) {
		cfg, err := cfgFromEnv()
		assert.NoError(t, err, "Retrieves without an error")
		assert.Equal(t, true, cfg.Disabled, "Disabled by default")
	})
	t.Run("When_Jaeger_On_As_Sidecar", func(t *testing.T) {
		os.Setenv(serviceNameVar, "Jaeger Utils Test")
		os.Setenv(sidecarEnabledVar, "true")
		os.Setenv(jaegerEnabledVar, "true")
		defer os.Unsetenv(sidecarEnabledVar)
		defer os.Unsetenv(serviceNameVar)
		defer os.Unsetenv(jaegerEnabledVar)

		_, err := cfgFromEnv()
		assert.NoError(t, err, "Retrieves without an error")
	})
	t.Run("When_Jaeger_On_As_Sidecar_Unparseable_Bool", func(t *testing.T) {
		os.Setenv(serviceNameVar, "Jaeger Utils Test")
		os.Setenv(sidecarEnabledVar, "NoParse")
		os.Setenv(jaegerEnabledVar, "true")
		defer os.Unsetenv(sidecarEnabledVar)
		defer os.Unsetenv(serviceNameVar)
		defer os.Unsetenv(jaegerEnabledVar)

		_, err := cfgFromEnv()
		assert.Error(t, err, "Retrieves with an error")
	})
	t.Run("When_Jaeger_On_Without_Sidecar", func(t *testing.T) {
		os.Setenv(jaegerEnabledVar, "true")
		os.Setenv(jaegerHostVar, "localhost")
		defer os.Unsetenv(jaegerEnabledVar)
		defer os.Unsetenv(jaegerHostVar)

		_, err := cfgFromEnv()
		assert.NoError(t, err, "Retrieves without an error")
	})
	t.Run("When_Jaeger_On_Inconsistent_Env", func(t *testing.T) {
		os.Setenv(jaegerEnabledVar, "true")
		defer os.Unsetenv(jaegerEnabledVar)

		_, err := cfgFromEnv()
		assert.Error(t, err, "Error raised")
	})
	t.Run("When_Jaeger_Off_Inconsistent_Env", func(t *testing.T) {
		os.Setenv(jaegerEnabledVar, "false")
		defer os.Unsetenv(jaegerEnabledVar)

		cfg, err := cfgFromEnv()
		assert.NoError(t, err, "Retrieves without an error")
		assert.Equal(t, true, cfg.Disabled, "Jaeger is disabled")
	})
}

// TestSetupGlobalTracer tests if the global tracer can be set up
func TestSetupGlobalTracer(t *testing.T) {
	t.Run("When_Empty_Env", func(t *testing.T) {
		_, err := SetupGlobalTracer()
		assert.NoError(t, err, "Sets up without an error")
	})
	t.Run("When_Bad_Env", func(t *testing.T) {
		os.Setenv(serviceNameVar, "Jaeger Utils Test")
		os.Setenv(sidecarEnabledVar, "NoParse")
		os.Setenv(jaegerEnabledVar, "true")
		os.Setenv(jaegerSamplerTypeVar, "const")
		os.Setenv(jaegerSamplerParamVar, "UNKNOWN")
		defer os.Unsetenv(sidecarEnabledVar)
		defer os.Unsetenv(serviceNameVar)
		defer os.Unsetenv(jaegerEnabledVar)
		defer os.Unsetenv(jaegerSamplerParamVar)
		defer os.Unsetenv(jaegerSamplerTypeVar)
		_, err := SetupGlobalTracer()
		assert.Error(t, err, "Error raised")
	})

	t.Run("When_Jaeger_Host_Unresolvable", func(t *testing.T) {
		os.Setenv(jaegerEnabledVar, "true")
		os.Setenv(jaegerHostVar, "un-resolvable")
		defer os.Unsetenv(jaegerEnabledVar)
		defer os.Unsetenv(jaegerHostVar)

		_, err := SetupGlobalTracer()
		assert.Error(t, err, "Error raised")
	})
}

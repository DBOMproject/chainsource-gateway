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
	"chainsource-gateway/helpers"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

var log = helpers.GetLogger("TracingUtil")

const defaultUDPSpanServerPort = 6831
const serviceNameVar = "JAEGER_SERVICE_NAME"
const sidecarEnabledVar = "JAEGER_AGENT_SIDECAR_ENABLED"
const jaegerEnabledVar = "JAEGER_ENABLED"
const jaegerSamplerTypeVar = "JAEGER_SAMPLER_TYPE"
const jaegerSamplerParamVar = "JAEGER_SAMPLER_PARAM"
const jaegerHostVar = "JAEGER_HOST"

// SetupGlobalTracer sets up a global tracer
func SetupGlobalTracer() (closer io.Closer, err error) {
	cfg, err := cfgFromEnv()
	if err != nil {
		log.Err(err).Msg("Could not parse Jaeger env vars")
		return
	}
	logger := NewZeroLogJaegerLogger(helpers.GetLogger("Jaeger"))

	tracer, closer, err := cfg.NewTracer(
		config.Logger(logger),
		config.Metrics(metrics.NullFactory))
	if err != nil {
		if _, ok := err.(*net.DNSError); ok {
			log.Err(err).Msg("DNS Lookup for the Jaeger agent host failed")
		}
		return
	}
	opentracing.SetGlobalTracer(tracer)
	return
}

func cfgFromEnv() (cfg *config.Configuration, err error) {
	cfg = &config.Configuration{}
	var sidecarEnabled bool
	if helpers.ExistsInEnv(serviceNameVar) {
		cfg.ServiceName = os.Getenv(serviceNameVar)
	} else {
		cfg.ServiceName = "Chainsource Gateway"
	}

	cfg.Sampler, err = getSamplerConfigFromEnv()
	if err != nil {
		return
	}

	cfg.Reporter = &config.ReporterConfig{}
	if helpers.ExistsInEnv(sidecarEnabledVar) {
		sidecarEnabled, err = strconv.ParseBool(os.Getenv(sidecarEnabledVar))
		if err != nil {
			return
		}
		if sidecarEnabled {
			cfg.Reporter.LocalAgentHostPort = "localhost:6831"
		}
	} else {
		cfg.Reporter.LocalAgentHostPort = fmt.Sprintf("%s:%d", os.Getenv(jaegerHostVar), defaultUDPSpanServerPort)
	}

	if helpers.ExistsInEnv(jaegerEnabledVar) {
		var enabled bool
		enabled, err = strconv.ParseBool(os.Getenv(jaegerEnabledVar))
		cfg.Disabled = !enabled
		if enabled {
			log.Info().Msg("Jaeger tracing enabled")
			if !sidecarEnabled && (!helpers.ExistsInEnv(jaegerHostVar) || os.Getenv(jaegerHostVar) == "") {
				log.Error().Msg("JAEGER_ENABLED is set to true but no agent address or sidecar configured to send traces to. " +
					"Either set JAEGER_ENABLED to false or set one of JAEGER_HOST or JAEGER_AGENT_SIDECAR_ENABLED to their " +
					"appropriate values.")
				err = errors.New("inconsistent Jaeger environment variables")
				return
			}
			log.Debug().Msgf("Jaeger Config:\n%s", helpers.PrettyInterfaceFormat(cfg))
		}
	} else {
		cfg.Disabled = true
	}

	return
}

func getSamplerConfigFromEnv() (samplerConfig *config.SamplerConfig, err error) {
	samplerConfig = &config.SamplerConfig{}

	if helpers.ExistsInEnv(jaegerSamplerTypeVar) {
		samplerConfig.Type = os.Getenv(jaegerSamplerTypeVar)
	} else {
		samplerConfig.Type = "const"
	}

	if helpers.ExistsInEnv(jaegerSamplerParamVar) {
		samplerConfig.Param, err = strconv.ParseFloat(os.Getenv(jaegerSamplerParamVar), 64)
	} else {
		samplerConfig.Param = 1
	}

	if err != nil {
		return
	}
	return
}

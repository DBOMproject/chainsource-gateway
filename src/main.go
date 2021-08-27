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

package main

import (
	"chainsource-gateway/agent"
	"chainsource-gateway/helpers"
	"chainsource-gateway/routes"
	"chainsource-gateway/tracing"
	"fmt"
	"net/http"
	"runtime"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var logger = helpers.GetLogger("Main")

func main() {
	address := helpers.GetServiceAddress()

	// Initialise log level from env
	helpers.SetLogLevelFromEnv()

	// Setup Jaeger distributed tracing from env, defer graceful shutdown
	closer, err := tracing.SetupGlobalTracer()
	if err != nil {
		logger.Err(err).Msg("Unable to initialize Jaeger tracer. Falling back to the NoopTracer")
	} else {
		defer closer.Close()
	}

	// Initialise AgentConfig and setup hot reload
	agent.GetAgentConfig()

	// Setup chi HTTP Server
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	helpers.SetupLoggingMiddleware(r)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Chainsource Gateway. runtime:" + runtime.Version()))
	})
	r.Mount("/api/v1", routes.AssetRouter())

	logger.Info().Msgf("Setting up on %s", address)
	err = http.ListenAndServe(address, r)
	if err != nil {
		panic(fmt.Errorf("failed to start server: %s", err))
	}
}

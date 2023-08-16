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

package main

import (
	"chainsource-gateway/helpers"
	"chainsource-gateway/routes"
	"chainsource-gateway/tracing"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"runtime"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var logger = helpers.GetLogger("Main")

func main() {
	var wg sync.WaitGroup
	address := helpers.GetServiceAddress()
	federationAddress := helpers.GetFederationAddress()
	// Initialize log level from env
	helpers.SetLogLevelFromEnv()

	// Setup Jaeger distributed tracing from env, defer graceful shutdown
	closer, err := tracing.SetupGlobalTracer()
	if err != nil {
		logger.Err(err).Msg("Unable to initialize Jaeger tracer. Falling back to the NoopTracer")
	} else {
		defer closer.Close()
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		// Setup chi HTTP Server
		r := chi.NewRouter()
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Recoverer)
		helpers.SetupLoggingMiddleware(r)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("chainsource-gateway user and client APIs runtime:" + runtime.Version()))
		})
		r.Mount("/api/v2", routes.GatewayRouter())

		logger.Info().Msgf("setting up user and client APIs on %s", address)
		err := http.ListenAndServe(address, r)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		// Setup chi HTTPS Federation Server
		r := chi.NewRouter()
		r.Use(middleware.RequestID)
		r.Use(middleware.RealIP)
		r.Use(middleware.Recoverer)
		helpers.SetupLoggingMiddleware(r)
		certFile, keyFile := helpers.GetNodeCertificate()
		caCertLocation := helpers.GetCACertificate()
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			panic(err)
		}

		caCert, err := os.ReadFile(caCertLocation)
		if err != nil {
			panic(err)
		}
		caCertPool := x509.NewCertPool()
		if err != nil {
			panic(err)
		}
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientCAs:    caCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		}

		server := &http.Server{
			Addr:      federationAddress,
			TLSConfig: tlsConfig,
			Handler:   r,
		}
		logger.Info().Msgf("Setting up Federation APIs on %s", federationAddress)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("chainsource-gateway federation APIs runtime:" + runtime.Version()))
		})
		r.Mount("/api/v2", routes.FederationRouter())

		err = server.ListenAndServeTLS(certFile, keyFile)
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}

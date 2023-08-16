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

package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/go-chi/chi"

	//"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

const DefaultLogLevel = zerolog.InfoLevel
const LogLevelVar = "LOG_LEVEL"

// GetLogger gets a zerolog logger with the "from" parameter set to the string sent to it as a parameter
func GetLogger(component string) zerolog.Logger {
	return zerolog.New(os.Stdout).
		With().
		Str("from", component).
		Timestamp().
		Logger().
		Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

// GetLoggerMiddleware returns a go chi middleware that logs HTTP requests
func GetLoggerMiddleware(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := logger.With().Logger()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			// Don't log liveliness check
			if fmt.Sprint(r.URL) == "/" {
				next.ServeHTTP(ww, r)
				return
			}
			defer func() {
				// Recover and record stack traces in case of a panic
				if rec := recover(); rec != nil {
					fmt.Print(string(debug.Stack()))
					log.Error().
						Timestamp().
						Interface("recover_info", rec).
						Bytes("debug_stack", debug.Stack()).
						Msg("log system error")
					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
				// log end request
				log.Info().
					Timestamp().
					Msgf("%d [%s] %s from %s ", ww.Status(), r.Method, r.URL, r.RemoteAddr)
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

// PrettyInterfaceFormat is the format the passed interface as pretty-printed json and returns it as a string
func PrettyInterfaceFormat(x interface{}) string {
	b, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		return "{{INTERFACE PARSE ERROR}}"
	}
	return string(b)
}

// Setup logging for the REST API
func SetupLoggingMiddleware(r *chi.Mux) {
	zeroLogHTTP := GetLogger("HTTP")
	r.Use(GetLoggerMiddleware(&zeroLogHTTP))
	zeroLogHTTP.Info().Msg("HTTP Logging Middleware Ready")
}

// Get the LOG_LEVEL from the environment as a level understood by zerolog, and then set the global log level to that value
func SetLogLevelFromEnv() {
	var setLogLevel = DefaultLogLevel
	logLevel := os.Getenv(LogLevelVar)
	if logLevel != "" {
		parsedLogLevel, err := zerolog.ParseLevel(logLevel)
		if err != nil {
			log.Warn().Msg("Could not parse LOG_LEVEL from environment")
		} else {
			setLogLevel = parsedLogLevel
		}
	} else {
		log.Warn().Msg("No LOG_LEVEL present in environment")
	}
	zerolog.SetGlobalLevel(setLogLevel)
	log.Info().Msgf("Set global log level to %s", setLogLevel)
}

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

// Package tracing contains utility types and functions to facilitate Jaeger Tracing
package tracing

import (
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog"
)

// ZerologJaegerLogger is a wrapper for zerolog to work with the Jaeger Logger interface
type ZerologJaegerLogger struct {
	logInstance zerolog.Logger
}

// NewZeroLogJaegerLogger creates a new instance of a logger the implements the Jaeger Logger interface
func NewZeroLogJaegerLogger(logInstance zerolog.Logger) *ZerologJaegerLogger {
	jaegerLogger := new(ZerologJaegerLogger)
	jaegerLogger.logInstance = logInstance
	return jaegerLogger
}

// Error adds an error log to the span
func (jaegerLogger ZerologJaegerLogger) Error(msg string) {
	jaegerLogger.logInstance.Error().Msg(msg)
}

// Infof adds an info log to the span
func (jaegerLogger ZerologJaegerLogger) Infof(msg string, args ...interface{}) {
	jaegerLogger.logInstance.Info().Msgf(strings.TrimSpace(msg), args...)
}

// LogAndTraceErr is a convenience function that sets an error on the trace and logs it
func LogAndTraceErr(logger zerolog.Logger, span opentracing.Span, err error, msg string) {
	logger.Err(err).Msg(msg)
	if span != nil {
		span.SetTag("error", true)
		span.SetTag("error.description", err.Error())
	}
}

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
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

// TestNewZeroLogJaegerLogger tests if the ZeroLogJaegerLogger getter works
func TestNewZeroLogJaegerLogger(t *testing.T) {
	assert.NotPanics(t, func() {
		NewZeroLogJaegerLogger(helpers.GetLogger("Test Jaeger"))
	}, "Returns a new Zerolog logger implementation without panicking")
}

// TestNewZeroLogJaegerLogger tests if the ZeroLogJaegerLogger error function works
func TestZerologJaegerLogger_Error(t *testing.T) {
	jLogger := NewZeroLogJaegerLogger(helpers.GetLogger("Test Jaeger"))
	assert.NotPanics(t, func() {
		jLogger.Error("Test Error")
	}, "Logs without panic")
}

// TestNewZeroLogJaegerLogger tests if the ZeroLogJaegerLogger infof function works
func TestZerologJaegerLogger_Infof(t *testing.T) {
	jLogger := NewZeroLogJaegerLogger(helpers.GetLogger("Test Jaeger"))
	assert.NotPanics(t, func() {
		jLogger.Infof("Test Error")
	}, "Logs without panic")
}

// TestNewZeroLogJaegerLogger tests if the log and trace helper function works
func TestLogAndTraceErr(t *testing.T) {
	span := opentracing.StartSpan("NewSpan")
	assert.NotPanics(t, func() {
		LogAndTraceErr(zerolog.Logger{}, span, errors.New("ERR"), "ERR")
	}, "Logs without panic")
}

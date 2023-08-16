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
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

// GetOKHandler returns a HTTP handler function that always returns OK
func GetOKHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}
	return fn
}

// GetOKHandler returns a HTTP handler function that always panics
func GetPanicHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		panic(errors.New("expected panic"))
	}
	return fn
}

// TestPrettyInterfaceFormat tests the interface to string converter
func TestPrettyInterfaceFormat(t *testing.T) {

	t.Run("Happy Case", func(t *testing.T) {
		mapIface := make(map[string]string)
		mapIface["foo"] = "bar"
		assert.NotPanics(t, func() { PrettyInterfaceFormat(mapIface) }, "Interface log formatter does not panic")
	})

	t.Run("Happy Case", func(t *testing.T) {
		assert.Contains(t, PrettyInterfaceFormat(math.Inf(1)), "ERROR", "Interface log formatter returns error")
	})
}

// TestGetLogger tests the logger getter
func TestGetLogger(t *testing.T) {
	assert.NotPanics(t, func() { GetLogger("test") }, "Logger getter does not panic")
}

// TestLoggerMiddleware tests the logger middleware
func TestLoggerMiddleware(t *testing.T) {
	logger := GetLogger("TestHTTP")
	assert.NotPanics(t, func() { GetLoggerMiddleware(&logger) }, "Logger middleware getter does not panic")
	loggerMiddleware := GetLoggerMiddleware(&logger)

	t.Run("For_OK_Request", func(t *testing.T) {
		mockRequest := httptest.NewRequest("POST", "http://mock-target/", strings.NewReader(""))
		responseRecorder := httptest.NewRecorder()
		loggerMiddleware(GetOKHandler()).ServeHTTP(responseRecorder, mockRequest)
		assert.Equal(t, http.StatusOK, responseRecorder.Code, "A 200 OK is returned")
	})
	t.Run("For_Liveliness_Case", func(t *testing.T) {
		mockRequest := httptest.NewRequest("POST", "/", strings.NewReader(""))
		responseRecorder := httptest.NewRecorder()
		loggerMiddleware(GetOKHandler()).ServeHTTP(responseRecorder, mockRequest)
		assert.Equal(t, http.StatusOK, responseRecorder.Code, "A 200 OK is returned")
	})
	t.Run("For_Panic", func(t *testing.T) {
		mockRequest := httptest.NewRequest("POST", "http://mock-target/", strings.NewReader(""))
		responseRecorder := httptest.NewRecorder()
		t.Log("\n==== A stack trace being displayed is expected behaviour for this test ====\n")
		assert.NotPanics(t, func() {
			loggerMiddleware(GetPanicHandler()).ServeHTTP(responseRecorder, mockRequest)
		}, "Recovers from panic after dumping error to logs")
		t.Log("\n==== END panic test ====\n")
	})
}

// TestSetupLoggingMiddleware tests the logger middleware setup function
func TestSetupLoggingMiddleware(t *testing.T) {
	assert.NotPanics(t, func() {
		SetupLoggingMiddleware(chi.NewRouter())
	}, "Can be setup without panicking")
}

// TestSetLogLevelFromEnv tests the logger level setter function
func TestSetLogLevelFromEnv(t *testing.T) {
	t.Run("When_OK_Environment", func(t *testing.T) {
		os.Setenv(LogLevelVar, "info")
		defer os.Unsetenv(LogLevelVar)
		defer zerolog.SetGlobalLevel(zerolog.Disabled)
		assert.NotPanics(t, func() {
			SetLogLevelFromEnv()
		}, "Can be set without panicking")

	})
	t.Run("When_Bad_Environment", func(t *testing.T) {
		os.Setenv(LogLevelVar, "badlevel")
		defer os.Unsetenv(LogLevelVar)
		defer zerolog.SetGlobalLevel(zerolog.Disabled)
		assert.NotPanics(t, func() {
			SetLogLevelFromEnv()
		}, "Can be set without panicking")
	})
	t.Run("When_No_Environment", func(t *testing.T) {
		defer zerolog.SetGlobalLevel(zerolog.Disabled)
		assert.NotPanics(t, func() {
			SetLogLevelFromEnv()
		}, "Can be set without panicking")
	})
}

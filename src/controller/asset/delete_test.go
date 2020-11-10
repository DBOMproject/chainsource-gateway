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

package asset

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestDelete ensures that Delete is unimplemented
func TestDelete(t *testing.T) {
	mockRequest := httptest.NewRequest("DELETE", "/", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteAsset)
	handler.ServeHTTP(responseRecorder, mockRequest)

	assert.Equal(t, http.StatusNotImplemented, responseRecorder.Result().StatusCode, "Response Should be 501 UNIMPLEMENTED")
}

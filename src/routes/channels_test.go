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

package routes

import (
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

// TestChannelRouter tests if the channel router initializes successfully
func TestChannelRouter(t *testing.T) {
	assert.NotPanics(t, func() {
		ChannelRouter()
	}, "Router initializes without panic")
}

// Test_channelSubRouting tests if the channel sub router mounts successfully
func Test_channelSubRouting(t *testing.T) {
	assert.NotPanics(t, func() {
		channelSubRouting(chi.NewRouter())
	}, "Router mounts without panic")
}

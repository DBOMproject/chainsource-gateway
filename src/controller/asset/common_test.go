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
	"chainsource-gateway/helpers"
	"chainsource-gateway/mocks"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestAssetContextGetterWithProviderFailureConditions contains tests for failures associated with the AgentProvider
func TestAssetContextGetterWithProviderFailureConditions(t *testing.T) {
	t.Run("No_Usable_Agent_For_RepoID", func (t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := mocks.InjectAgentProviderIntoContext(context.Background(), ctrl, nil)

		_, _, err := getChildAssetContextFromAssetElement(ctx,  helpers.AssetElement{
			RepoID: "NIL_REPO",
		})
		assert.Error(t, err, "Agent unavailable case must be handled gracefully" )

	})
}

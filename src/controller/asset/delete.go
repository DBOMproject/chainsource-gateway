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
	"chainsource-gateway/responses"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

// DeleteAsset is a controller function to delete an asset on a channel on the repository with a given assetID
// All agents would not be able to implement delete since some of them interact with a immutable DLT
func DeleteAsset(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, responses.ErrUnimplemented(errors.New("use PUT to update the asset instead")))
}

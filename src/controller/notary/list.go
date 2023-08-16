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

package notary

import (
	"chainsource-gateway/helpers"
	"chainsource-gateway/responses"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

// ListNotaries is a controller function to list all the locally configured notaries
func ListNotaries(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, responses.ErrUnimplemented(errors.New(helpers.NotImplemented)))
}

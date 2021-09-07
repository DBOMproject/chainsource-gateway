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

// Package routes contains all the routes for the Gateway API
package routes

import (
	"chainsource-gateway/agent"
	"chainsource-gateway/controller/channel"
	"chainsource-gateway/helpers"
	"chainsource-gateway/responses"
	"chainsource-gateway/tracing"
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"
)

var channelLog = helpers.GetLogger("Router")

// channelSubRouting defines the sub routes for the channel APIs
func channelSubRouting(r chi.Router) {
	r.Use(injectSpanMiddleware)
	r.Use(agentProvider)
	r.Use(signingServiceProvider)
	r.Use(channelContext)
	r.Use(unmarshalBody)

	// List assets
	r.Get("/asset", channel.ListAssets)

}

// channelContext parses the context for the channel APIs
func channelContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, ctx := opentracing.StartSpanFromContext(r.Context(), "Embedding Asset Context")

		// Append context with routing details

		vars := helpers.AssetRoutingVars{
			RepoID:    chi.URLParam(r, "repoID"),
			ChannelID: chi.URLParam(r, "channelID"),
		}

		err := errors.New("Invalid URL")

		if vars.ChannelID == "" || vars.RepoID == "" {
			render.Render(w, r, responses.ErrInvalidRequest(err))
			tracing.LogAndTraceErr(channelLog, span, err, "Empty URL parameters found")
			return
		}

		agentProvider := r.Context().Value("agentProvider").(agent.Provider)
		channelLog.Debug().Msgf("ctx: %+v ", vars)
		ctx = context.WithValue(r.Context(), "assetVars", vars)

		// Getting the agent-config for the repo we're trying to access
		agentConfig, err := agentProvider.GetAgentConfigForRepo(vars.RepoID)
		if err != nil {
			tracing.LogAndTraceErr(channelLog, span, err, "No agent configured for request")
			_ = render.Render(w, r, responses.ErrNoAgent(err))
			return
		}

		requestAgent := agentProvider.NewAgent(&agentConfig)
		ctx = context.WithValue(ctx, "agent", requestAgent)
		span.Finish()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

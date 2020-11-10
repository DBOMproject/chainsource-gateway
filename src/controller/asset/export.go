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
	"chainsource-gateway/agent"
	"chainsource-gateway/helpers"
	"chainsource-gateway/responses"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/render"
	"github.com/opentracing/opentracing-go"
)

type exportAsset struct {
	id    string
	asset *helpers.Asset
}

// ExportAsset exports a DBoM asset as JSON
func ExportAsset(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), "Export Asset")
	defer span.Finish()
	assetVars := r.Context().Value("assetVars").(helpers.AssetRoutingVars)
	exportVars := r.Context().Value("exportVars").(helpers.ExportRoutingVars)
	requestAgent := r.Context().Value("agent").(agent.Agent)
	log.Debug().Msgf("Getting %s/%s from agent at %s:%d", assetVars.ChannelID, assetVars.AssetID,
		requestAgent.GetHost(), requestAgent.GetPort())

	var result helpers.Asset
	resultStream, err := requestAgent.QueryStream(ctx, agent.QueryArgs{
		ChannelID: assetVars.ChannelID,
		AssetID:   assetVars.AssetID,
	})
	if err != nil {
		if err == helpers.ErrNotFound {
			render.Render(w, r, responses.ErrDoesNotExist(err))
		} else if err == helpers.ErrUnauthorized{
			render.Render(w, r, responses.ErrAgentUnauthorized(err))
		} else {
			render.Render(w, r, responses.ErrAgent(err))
		}
		return
	}
	err = json.NewDecoder(resultStream).Decode(&result)
	if err != nil {
		render.Render(w, r, responses.ErrInternalServer(err))
		return
	}

	var wg sync.WaitGroup
	var channel = make(chan exportAsset, 1)
	var parentChannel = make(chan exportAsset, 1)
	var fatalErrors = make(chan error, 1)
	var parentFatalErrors = make(chan error, 1)
	var seenAssets = make([]string, 0)
	wg.Add(1)
	go ExportChildren(ctx, assetVars.AssetID, result, seenAssets, &wg, &channel, &fatalErrors)
	wg.Add(1)
	go ExportParent(ctx, assetVars.AssetID, result, seenAssets, &wg, &parentChannel, &parentFatalErrors)
	wg.Wait()
	close(channel)
	close(parentChannel)
	close(fatalErrors)
	close(parentFatalErrors)
	for err := range fatalErrors {
		render.Render(w, r, responses.ErrFailedExport(err))
		return
	}
	for err := range parentFatalErrors {
		render.Render(w, r, responses.ErrFailedExport(err))
		return
	}

	var children = make(map[string]*helpers.Asset)
	var parents = make(map[string]*helpers.Asset)
	for child := range channel {
		children[child.id] = child.asset
	}
	for parent := range parentChannel {
		parents[parent.id] = parent.asset
	}
	if len(parents) == 1 {
		parent := parents[assetVars.AssetID].Parent
		children[assetVars.AssetID].Parent = parent
	}

	result.AttachedChildren = nil

	if exportVars.InlineResponse == "true" {
		render.JSON(w, r, children)
	} else if exportVars.FileName != "" {
		w.Header().Set("Content-Disposition", "attachment; filename="+exportVars.FileName)
		render.JSON(w, r, children)
	} else {
		w.Header().Set("Content-Disposition", "attachment; filename="+assetVars.AssetID+".json")
		render.JSON(w, r, children)
	}
}

//ExportChildren recursively populates the children of the asset
func ExportChildren(ctx context.Context, assetID string, asset helpers.Asset, seenAssets []string, parentWG *sync.WaitGroup, parentChan *chan exportAsset, fatalErrors *chan error) {
	agentProvider := ctx.Value("agentProvider").(agent.Provider)
	var children = make(map[string]*helpers.Asset)
	var exAsset exportAsset
	seenAssets = append(seenAssets, assetID)
	if len(asset.AttachedChildren) > 0 {
		var wg sync.WaitGroup
		var channel = make(chan exportAsset, len(asset.AttachedChildren))
		for i, child := range asset.AttachedChildren {
			log.Info().Msg(strconv.Itoa(i))
			log.Info().Msg(child.AssetID)
			log.Debug().Msg("Check for loop")
			if Find(seenAssets, child.AssetID) {
				log.Error().Msg("Parent Child Loop Detected for asset " + assetID)
				err := errors.New("Parent Child Loop Detected for asset ")
				(*fatalErrors) <- err
				(*parentWG).Done()
				return
			}
			var result2 helpers.Asset
			agentConfig, err := agentProvider.GetAgentConfigForRepo(child.RepoID)
			childAgent := agentProvider.NewAgent(&agentConfig)
			resultStream2, err := childAgent.QueryStream(ctx, agent.QueryArgs{
				ChannelID: child.ChannelID,
				AssetID:   child.AssetID,
			})
			if err != nil {
				(*fatalErrors) <- err
				(*parentWG).Done()
				return
			}
			err = json.NewDecoder(resultStream2).Decode(&result2)
			if err != nil {
				(*parentChan) <- exAsset
				(*fatalErrors) <- err
				(*parentWG).Done()
				return
			}
			wg.Add(1)
			go ExportChildren(ctx, child.AssetID, result2, seenAssets, &wg, &channel, fatalErrors)
		}
		wg.Wait()
		close(channel)
		for child := range channel {
			children[child.id] = child.asset
		}
	}
	asset.AttachedChildren = nil
	asset.ParentAsset = nil
	asset.Children = children
	exAsset.id = assetID
	exAsset.asset = &asset
	(*parentChan) <- exAsset
	(*parentWG).Done()
}

//ExportParent recursively populates the parent of the asset
func ExportParent(ctx context.Context, assetID string, asset helpers.Asset, seenAssets []string, parentWG *sync.WaitGroup, parentChan *chan exportAsset, fatalErrors *chan error) {
	var parents = make(map[string]*helpers.Asset)
	agentProvider := ctx.Value("agentProvider").(agent.Provider)
	var exAsset exportAsset
	seenAssets = append(seenAssets, assetID)
	if asset.ParentAsset != nil && (*asset.ParentAsset).AssetID != "" {
		log.Debug().Msg("Check for loop")
		if Find(seenAssets, (*asset.ParentAsset).AssetID) {
			log.Error().Msg("Parent Child Loop Detected for asset " + assetID)
			err := errors.New("Parent Child Loop Detected for asset ")
			(*fatalErrors) <- err
			(*parentWG).Done()
			return
		}
		var wg sync.WaitGroup
		var channel = make(chan exportAsset, 1)
		log.Info().Msg((*asset.ParentAsset).AssetID)
		var result2 helpers.Asset
		agentConfig, err := agentProvider.GetAgentConfigForRepo((*asset.ParentAsset).RepoID)
		parentAgent := agentProvider.NewAgent(&agentConfig)
		resultStream2, err := parentAgent.QueryStream(ctx, agent.QueryArgs{
			ChannelID: (*asset.ParentAsset).ChannelID,
			AssetID:   (*asset.ParentAsset).AssetID,
		})
		if err != nil {
			(*fatalErrors) <- err
			(*parentWG).Done()
			return
		}
		err = json.NewDecoder(resultStream2).Decode(&result2)
		if err != nil {
			*parentChan <- exAsset
			(*fatalErrors) <- err
			(*parentWG).Done()
			return
		}
		wg.Add(1)
		go ExportParent(ctx, (*asset.ParentAsset).AssetID, result2, seenAssets, &wg, &channel, fatalErrors)
		wg.Wait()
		close(channel)
		for parent := range channel {
			parents[parent.id] = parent.asset
		}
	}
	asset.ParentAsset = nil
	asset.Parent = parents
	asset.AttachedChildren = nil
	asset.ParentAsset = nil
	exAsset.id = assetID
	exAsset.asset = &asset
	*parentChan <- exAsset
	(*parentWG).Done()
}

func Find(slice []string, val string) bool {
	for i, item := range slice {
		if item == val {
			log.Debug().Msg(strconv.Itoa(i))
			return true
		}
	}
	return false
}

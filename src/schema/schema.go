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

// Package schema contains the utilities and interfaces to validate schemas
// for the controllers
package schema

import (
	"chainsource-gateway/tracing"
	"context"
	json2 "encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/qri-io/jsonschema"
)

// A JSONSchema Validator to validate a json object in memory against a JSONSchema file
func coreValidator(ctx context.Context, json map[string]interface{}, schemaFileName string) (schemaErrors string, isValid bool, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Validate Schema")
	defer span.Finish()
	file, err := os.Open(schemaFileName)
	if err != nil {
		tracing.LogAndTraceErr(schemaLog, span, err, "Schema Open Error")
		panic("open schema: " + err.Error())
	}
	schemaRoot := &jsonschema.RootSchema{}
	if err := json2.NewDecoder(file).Decode(schemaRoot); err != nil {
		tracing.LogAndTraceErr(schemaLog, span, err, "Schema Unmarshal Error")
		panic("unmarshal schema: " + err.Error())
	}
	var errs []jsonschema.ValError
	schemaRoot.Validate("/", json, &errs)

	if len(errs) > 0 {
		isValid = false
		schemaErrors = fmt.Sprint(errs)
		tracing.LogAndTraceErr(schemaLog, span, errors.New(schemaErrors), "Schema Incorrect")
	} else {
		isValid = true
	}
	return
}

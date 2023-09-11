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

// Package helpers contains common functions used by the rest of the Gateway
package helpers

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/render"
)

var log = GetLogger(HttpClient)

// Timeout for a HTTP Request to the agent.
// Currently unenforced
const timeout = 3

// PostJSONRequest executes a POST request to a url with a Text/JSON body and optionally additional headers
func PostJSONRequest(url string, body []byte) (resp *http.Response, err error) {
	certFile, keyFile := GetNodeCertificate()
	caCertFile := GetCACertificate()
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}

	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		panic(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}

	client := &http.Client{
		Timeout: time.Minute * timeout,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Err(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err = client.Do(req)

	if err != nil {
		log.Err(err)
		return
	}

	return
}

// GetRequest executes a GET request to a URL, with optional headers
func GetRequest(url string) (result io.ReadCloser, err error) {
	certFile, keyFile := GetNodeCertificate()
	caCertFile := GetCACertificate()
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Err(err)
		return
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}

	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		panic(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}

	client := &http.Client{
		Timeout: time.Minute * 3,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err)
		return
	}
	result = resp.Body

	return
}

// HandleError responds with an error JSON response.
func HandleError(w http.ResponseWriter, r *http.Request, errMessage string) {
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, ErrorResponse{Error: errMessage})
}

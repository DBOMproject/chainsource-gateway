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
	"os"
)

const defaultPort = "3050"
const federationPort = "7205"
const certPath = "/app/certs/"

func GetServiceAddress() (address string) {
	if port := os.Getenv("PORT"); port != "" {
		address = ":" + port
	} else {
		address = ":" + defaultPort
	}
	return
}

func GetFederationAddress() (address string) {
	if port := os.Getenv("FED_PORT"); port != "" {
		address = ":" + port
	} else {
		address = ":" + federationPort
	}
	return
}

func GetNodeID() string {
	nodeId := os.Getenv("NODE_ID")
	return nodeId
}

func GetNodeURI() string {
	nodeUri := os.Getenv("NODE_URI")
	return nodeUri
}

func ExistsInEnv(key string) (exists bool) {
	_, exists = os.LookupEnv(key)
	return
}

func GetCACertificate() (CAPath string) {
	CAPath = certPath + "ca.crt"
	return
}

func GetNodeCertificate() (certFile string, keyFile string) {
	certFile = certPath + os.Getenv("NODE_URI") + ".crt"
	keyFile = certPath + os.Getenv("NODE_URI") + ".key"
	return
}

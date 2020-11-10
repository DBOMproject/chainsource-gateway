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

package pgp

import (
	"chainsource-gateway/helpers"
	"os"
)

var log = helpers.GetLogger("PGPUtil")

const pgpServiceAddressVar = "PGP_SERVICE_ADDRESS"
const defaultPgpServiceAddress = "http://localhost:5000"

// GetPGPServiceAddress gets the address of the pgp service
func GetPGPServiceAddress() (pgpServiceAddress string) {
	log.Info().Msg("Get PGP Service Address")
	if helpers.ExistsInEnv(pgpServiceAddressVar) {
		pgpServiceAddress = os.Getenv(pgpServiceAddressVar)
	} else {
		pgpServiceAddress = defaultPgpServiceAddress
	}
	log.Info().Msg(pgpServiceAddress)
	return
}


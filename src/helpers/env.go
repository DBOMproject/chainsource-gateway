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

package helpers

import "os"

const defaultPort = "3005"

func GetServiceAddress() (address string) {
	if port := os.Getenv("PORT"); port != "" {
		address = ":" + port
	} else {
		address = ":" + defaultPort
	}
	return
}

func ExistsInEnv(key string) (exists bool) {
	_, exists = os.LookupEnv(key)
	return
}

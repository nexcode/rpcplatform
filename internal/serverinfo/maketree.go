/*
 * Copyright 2025 RPCPlatform Authors
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
 */

package serverinfo

import (
	"strings"

	"github.com/nexcode/rpcplatform/attributes"
)

func MakeTree(m map[string]string) map[string]*ServerInfo {
	serverInfoTree := map[string]*ServerInfo{}

	for key, value := range m {
		path := strings.SplitN(key, "/", 2)

		if serverInfoTree[path[0]] == nil {
			serverInfoTree[path[0]] = &ServerInfo{
				Attributes: attributes.New(),
			}
		}

		if len(path) == 1 {
			serverInfoTree[path[0]].Address = value
		} else {
			serverInfoTree[path[0]].Attributes.Load(path[1], value)
		}
	}

	return serverInfoTree
}

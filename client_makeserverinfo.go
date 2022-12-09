/*
 * Copyright 2022 RPCPlatform Authors
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

package rpcplatform

import (
	"strings"
)

type addrAndAttrs struct {
	address    string
	attributes map[string]string
}

func (c *Client) makeServerInfo(serverInfo map[string]string) map[string]*addrAndAttrs {
	tree := map[string]*addrAndAttrs{}

	for key, value := range serverInfo {
		path := strings.SplitN(key, "/", 2)

		if tree[path[0]] == nil {
			tree[path[0]] = &addrAndAttrs{
				attributes: map[string]string{},
			}
		}

		if len(path) == 1 {
			tree[path[0]].address = value
		} else {
			tree[path[0]].attributes[path[1]] = value
		}
	}

	return tree
}

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
	"github.com/nexcode/rpcplatform/internal/serverinfo"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

func (c *Client) updateState(init bool, serverInfoTree map[string]*serverinfo.ServerInfo) {
	state := resolver.State{
		Addresses: make([]resolver.Address, 0, len(serverInfoTree)),
	}

	for _, value := range serverInfoTree {
		state.Addresses = append(state.Addresses, resolver.Address{
			Addr:       value.Address,
			Attributes: attributes.New(struct{}{}, value.Attributes),
		})
	}

	if init {
		c.resolver.InitialState(state)
	} else {
		c.resolver.UpdateState(state)
	}
}

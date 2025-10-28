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

package options

import (
	"github.com/nexcode/rpcplatform/internal/config"
	"google.golang.org/grpc"
)

// Client provides options used when creating new client.
var Client = client{}

type client struct{}

// MaxActiveServers sets the maximum number of active servers the client will connect to.
// If there are more servers than this value, no requests will be made to them.
func (client) MaxActiveServers(count int) func(*config.Client) {
	return func(c *config.Client) {
		c.MaxActiveServers = count
	}
}

// GRPCOptions provide []grpc.DialOption to the client.
func (client) GRPCOptions(options ...grpc.DialOption) func(*config.Client) {
	return func(c *config.Client) {
		c.GRPCOptions = append(c.GRPCOptions, options...)
	}
}

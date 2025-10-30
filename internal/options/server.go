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
	"github.com/nexcode/rpcplatform/internal/attributes"
	"github.com/nexcode/rpcplatform/internal/config"
	"google.golang.org/grpc"
)

type Server struct{}

// PublicAddr is used when the server is not accessible to clients at the address it is located at.
func (Server) PublicAddr(publicAddr string) func(*config.Server) {
	return func(c *config.Server) {
		c.PublicAddr = publicAddr
	}
}

// Attributes are server settings that are accessible via the API.
func (Server) Attributes(attributes *attributes.Attributes) func(*config.Server) {
	return func(c *config.Server) {
		c.Attributes = attributes
	}
}

// GRPCOptions provide []grpc.ServerOption to the server.
func (Server) GRPCOptions(options ...grpc.ServerOption) func(*config.Server) {
	return func(c *config.Server) {
		c.GRPCOptions = append(c.GRPCOptions, options...)
	}
}

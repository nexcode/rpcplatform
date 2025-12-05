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

// PublicAddr sets the public address for the server when it is not accessible to clients at its listening address.
func (Server) PublicAddr(publicAddr string) func(*config.Server) {
	return func(c *config.Server) {
		c.PublicAddr = publicAddr
	}
}

// Attributes sets server attributes that are applied by the server and accessible via the Lookup method.
func (Server) Attributes(attributes *attributes.Attributes) func(*config.Server) {
	return func(c *config.Server) {
		c.Attributes = attributes
	}
}

// GRPCOptions adds gRPC server options to the server.
func (Server) GRPCOptions(options ...grpc.ServerOption) func(*config.Server) {
	return func(c *config.Server) {
		c.GRPCOptions = append(c.GRPCOptions, options...)
	}
}

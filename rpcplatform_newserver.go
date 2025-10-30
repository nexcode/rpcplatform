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
	"net"

	"github.com/nexcode/rpcplatform/internal/config"
	"github.com/nexcode/rpcplatform/internal/gears"
	"google.golang.org/grpc"
)

// NewServer creates a new server. You need to provide the server name and listening address.
func (p *RPCPlatform) NewServer(name, addr string, options ...ServerOption) (*Server, error) {
	config := config.NewServer()

	for _, option := range p.config.ServerOptions {
		option(config)
	}

	for _, option := range options {
		option(config)
	}

	if config.Attributes == nil {
		config.Attributes = NewAttributes()
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	if config.PublicAddr != "" {
		addr = config.PublicAddr
	} else {
		addr = listener.Addr().String()
	}

	id := gears.UID()

	if p.config.OpenTelemetry != nil {
		statsHandler, err := p.openTelemetry(id, listener.Addr(), addr)
		if err != nil {
			return nil, err
		}

		config.GRPCOptions = append(config.GRPCOptions, grpc.StatsHandler(statsHandler))
	}

	return &Server{
		id:         id,
		name:       p.etcdPrefix + gears.FixPath(name),
		etcd:       p.etcdClient,
		server:     grpc.NewServer(config.GRPCOptions...),
		listener:   listener,
		attributes: config.Attributes,
		publicAddr: addr,
	}, nil
}

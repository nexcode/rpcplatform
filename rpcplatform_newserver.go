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

	"github.com/nexcode/rpcplatform/internal/gears"
	"google.golang.org/grpc"
)

// NewServer creates a new server. You need to provide the server name, listening address and attributes.
// Optional param `publicAddr` used if server can't reachable by listening address.
// If no additional settings are needed, attributes can be nil.
func (p *RPCPlatform) NewServer(name, addr string, attributes *ServerAttributes, publicAddr ...string) (*Server, error) {
	if attributes == nil {
		attributes = Attributes().Server()
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	if len(publicAddr) != 0 {
		addr = publicAddr[0]
	} else {
		addr = listener.Addr().String()
	}

	id := gears.UID()

	if err = p.grpcinject(id, listener.Addr(), addr); err != nil {
		return nil, err
	}

	return &Server{
		id:         id,
		name:       p.config.EtcdPrefix + gears.FixPath(name),
		etcd:       p.config.EtcdClient,
		server:     grpc.NewServer(p.config.GRPCOptions.Server...),
		listener:   listener,
		attributes: attributes,
		publicAddr: addr,
	}, nil
}

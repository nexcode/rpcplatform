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
	"github.com/nexcode/rpcplatform/internal/balancer"
	"github.com/nexcode/rpcplatform/internal/gears"
	"github.com/nexcode/rpcplatform/internal/resolver"
	"google.golang.org/grpc"
)

// NewClient creates a new client. You need to provide the target server name.
// If no additional settings are needed, attributes can be nil.
func (p *RPCPlatform) NewClient(target string, attributes *ClientAttributes) (*Client, error) {
	if attributes == nil {
		attributes = Attributes().Client()
	}

	balancerName := gears.UID()
	balancer.Register(balancerName, attributes.maxActiveServers)

	resolver := resolver.NewResolver()

	options := append(p.config.GRPCOptions.Client,
		grpc.WithResolvers(resolver),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig":[{"`+balancerName+`":{}}]}`),
	)

	if err := p.grpcinject(gears.UID(), nil, ""); err != nil {
		return nil, err
	}

	c := &Client{
		target:   p.config.EtcdPrefix + gears.FixPath(target) + "/",
		etcd:     p.config.EtcdClient,
		resolver: resolver,
	}

	serverInfo, revision, err := c.stateInit()
	if err != nil {
		return nil, err
	}

	c.client, err = grpc.NewClient(target, options...)
	if err != nil {
		return nil, err
	}

	if err = c.stateWatcher(serverInfo, revision); err != nil {
		return nil, err
	}

	return c, nil
}

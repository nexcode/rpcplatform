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
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nexcode/rpcplatform/internal/balancer"
	"github.com/nexcode/rpcplatform/internal/config"
	"github.com/nexcode/rpcplatform/internal/gears"
	"github.com/nexcode/rpcplatform/internal/resolver"
	"google.golang.org/grpc"
)

// NewClient creates a new client connecting to the specified server name.
func (p *RPCPlatform) NewClient(target string, options ...ClientOption) (*Client, error) {
	if target == "" || strings.Contains(target, "/") {
		return nil, fmt.Errorf("%q: target is empty or contains «/»: %w", target, ErrInvalidTargetName)
	}

	config := config.NewClient()

	for _, option := range p.config.ClientOptions {
		option(config)
	}

	for _, option := range options {
		option(config)
	}

	c := &Client{
		id:       gears.UID(),
		target:   p.etcdPrefix + "/" + target + "/",
		resolver: resolver.New(),
		config:   config,
	}

	ctx, cancel := context.WithCancel(context.Background())

	serverInfoTree, err := p.Lookup(ctx, target, true)
	if err != nil {
		return nil, err
	}

	select {
	case serverInfoTree := <-serverInfoTree:
		c.updateState(true, serverInfoTree)
	case <-time.After(4 * time.Second):
		cancel()
	}

	if err = ctx.Err(); err != nil {
		return nil, err
	}

	config.GRPCOptions = append(config.GRPCOptions,
		grpc.WithResolvers(c.resolver),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig":[{"`+balancer.Name+`":{}}]}`),
	)

	if p.config.OpenTelemetry != nil {
		statsHandler, err := p.openTelemetry(c.id, nil, "")
		if err != nil {
			return nil, err
		}

		config.GRPCOptions = append(config.GRPCOptions, grpc.WithStatsHandler(statsHandler))
	}

	c.client, err = grpc.NewClient(c.resolver.Scheme()+":"+target, config.GRPCOptions...)
	if err != nil {
		return nil, err
	}

	go func() {
		for serverInfoTree := range serverInfoTree {
			c.updateState(false, serverInfoTree)
		}
	}()

	return c, nil
}

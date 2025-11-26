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

package balancer

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/endpointsharding"
	"google.golang.org/grpc/balancer/pickfirst"
)

func (builder) Build(cc balancer.ClientConn, opts balancer.BuildOptions) balancer.Balancer {
	b := &rpcBalancer{
		ClientConn: cc,
	}

	childBuilder := balancer.Get(pickfirst.Name).Build
	b.Balancer = endpointsharding.NewBalancer(b, opts, childBuilder, endpointsharding.Options{})

	return b
}

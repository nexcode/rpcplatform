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
	"github.com/nexcode/rpcplatform/internal/grpcattrs"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/pickfirst"
)

func (b *rpcBalancer) UpdateClientConnState(ccs balancer.ClientConnState) error {
	b.config = grpcattrs.GetClientConfig(ccs.ResolverState.Attributes)

	return b.Balancer.UpdateClientConnState(balancer.ClientConnState{
		ResolverState: pickfirst.EnableHealthListener(ccs.ResolverState),
	})
}

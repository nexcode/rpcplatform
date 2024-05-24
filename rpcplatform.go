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
	"github.com/nexcode/rpcplatform/internal/config"
	"github.com/nexcode/rpcplatform/internal/gears"
	"github.com/nexcode/rpcplatform/options"
	etcd "go.etcd.io/etcd/client/v3"
)

// New creates an RPCPlatform object for further creation of clients and servers.
// All methods of this object are thread safe. You can create this object once
// and use it in different places in your program.
func New(etcdPrefix string, etcdClient *etcd.Client, opts ...options.Option) (*RPCPlatform, error) {
	if etcdPrefix != "" {
		etcdPrefix = gears.FixPath(etcdPrefix)
	}

	return &RPCPlatform{
		config: options.Make(etcdClient, etcdPrefix, opts),
	}, nil
}

type RPCPlatform struct {
	config *config.Config
}

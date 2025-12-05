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
	"fmt"
	"strings"

	"github.com/nexcode/rpcplatform/internal/config"
	etcd "go.etcd.io/etcd/client/v3"
)

// New creates an RPCPlatform object for further creation of clients and servers.
// All methods of this object are thread safe. You can create this object once
// and use it in different places in your program.
func New(etcdPrefix string, etcdClient *etcd.Client, options ...PlatformOption) (*RPCPlatform, error) {
	if strings.Contains(etcdPrefix, "//") {
		return nil, fmt.Errorf("%q: prefix contains «//»: %w", etcdPrefix, ErrInvalidEtcdPrefix)
	}

	if etcdPrefix != "" {
		if etcdPrefix[0] != '/' {
			etcdPrefix = "/" + etcdPrefix
		}

		if etcdPrefix[len(etcdPrefix)-1] == '/' {
			etcdPrefix = etcdPrefix[:len(etcdPrefix)-1]
		}
	}

	config := config.NewPlatform()
	for _, option := range options {
		option(config)
	}

	rpcp := &RPCPlatform{
		etcdPrefix: etcdPrefix,
		etcdClient: etcdClient,
		config:     config,
	}

	return rpcp, nil
}

type RPCPlatform struct {
	etcdPrefix string
	etcdClient *etcd.Client
	config     *config.Platform
}

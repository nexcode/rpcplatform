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

package options

import (
	"github.com/nexcode/rpcplatform/internal/config"
	etcd "go.etcd.io/etcd/client/v3"
)

func Make(etcdClient *etcd.Client, etcdPrefix string, options []Option) *config.Config {
	config := &config.Config{
		EtcdClient: etcdClient,
		EtcdPrefix: etcdPrefix,
	}

	for _, option := range options {
		option.apply(config)
	}

	return config
}

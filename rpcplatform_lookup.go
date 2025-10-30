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

package rpcplatform

import (
	"context"
	"strings"

	"github.com/nexcode/rpcplatform/internal/gears"
	etcd "go.etcd.io/etcd/client/v3"
)

// Lookup returns information about available servers by name. If the watch is set to true, a new portion of data
// will be provided with each change. Otherwise the channel will be closed immediately after the first data is written.
// The keys of the returned map are server IDs.
func (p *RPCPlatform) Lookup(ctx context.Context, name string, watch bool) (<-chan map[string]*ServerInfo, error) {
	target := p.etcdPrefix + gears.FixPath(name) + "/"

	resp, err := p.etcdClient.Get(ctx, target, etcd.WithPrefix())
	if err != nil {
		return nil, err
	}

	serverInfoFlat := make(map[string]string, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		trimKey := strings.TrimPrefix(string(kv.Key), target)
		serverInfoFlat[trimKey] = string(kv.Value)
	}

	serverInfoTree := make(chan map[string]*ServerInfo, 1)
	serverInfoTree <- makeServerInfo(serverInfoFlat)

	if !watch {
		close(serverInfoTree)
		return serverInfoTree, nil
	}

	go func() {
		watchChan := p.etcdClient.Watch(ctx, target,
			etcd.WithPrefix(), etcd.WithRev(resp.Header.Revision+1),
		)

		for data := range watchChan {
			for _, event := range data.Events {
				trimKey := strings.TrimPrefix(string(event.Kv.Key), target)

				switch event.Type {
				case etcd.EventTypeDelete:
					delete(serverInfoFlat, trimKey)
				case etcd.EventTypePut:
					serverInfoFlat[trimKey] = string(event.Kv.Value)
				}
			}

			serverInfoTree <- makeServerInfo(serverInfoFlat)
		}

		close(serverInfoTree)
	}()

	return serverInfoTree, nil
}

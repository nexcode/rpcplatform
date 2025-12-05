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
	"fmt"
	"strings"

	etcd "go.etcd.io/etcd/client/v3"
)

// Lookup returns information about available servers with the given name.
// If watch is true, the returned channel sends updates whenever servers change.
// If watch is false, the channel closes after the first update.
// The returned map keys are server IDs.
func (p *RPCPlatform) Lookup(ctx context.Context, target string, watch bool) (<-chan map[string]*ServerInfo, error) {
	if target == "" || strings.Contains(target, "/") {
		return nil, fmt.Errorf("%q: target is empty or contains «/»: %w", target, ErrInvalidTargetName)
	}

	target = p.etcdPrefix + "/" + target + "/"

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

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
	etcd "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

func (c *Client) stateWatcher() error {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	resp, err := c.etcd.Get(ctx, c.target, etcd.WithPrefix())
	cancel()

	if err != nil {
		return err
	}

	serverInfo := make(map[string]string, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		trimKey := strings.TrimPrefix(string(kv.Key), c.target)
		serverInfo[trimKey] = string(kv.Value)
	}

	c.updateState(serverInfo)

	watchChan := c.etcd.Watch(context.Background(), c.target,
		etcd.WithPrefix(), etcd.WithRev(resp.Header.Revision+1),
	)

	go func() {
		for data := range watchChan {
			for _, event := range data.Events {
				trimKey := strings.TrimPrefix(string(event.Kv.Key), c.target)

				switch event.Type {
				case etcd.EventTypeDelete:
					delete(serverInfo, trimKey)
				case etcd.EventTypePut:
					serverInfo[trimKey] = string(event.Kv.Value)
				}
			}

			c.updateState(serverInfo)
		}
	}()

	return nil
}

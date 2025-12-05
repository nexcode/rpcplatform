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
	"time"

	"github.com/nexcode/rpcplatform/internal/attributes"
	etcd "go.etcd.io/etcd/client/v3"
)

// Serve starts the gRPC server and blocks until it exits or an error occurs.
func (s *Server) Serve() error {
	global, cancel := context.WithCancel(context.Background())

	defer func() {
		cancel()
	}()

	go func() {
		path := s.name + "/" + s.id
		attributes := attributes.Values(s.config.Attributes)

		for {
			if global.Err() != nil {
				break
			}

			ctx, cancel := context.WithTimeout(global, 4*time.Second)
			lease, err := s.etcd.Grant(ctx, 4)
			cancel()

			if err != nil {
				fmt.Println(err)
				continue
			}

			addr := s.config.PublicAddr
			if addr == "" {
				addr = s.listener.Addr().String()
			}

			ops := make([]etcd.Op, 0, len(attributes)/2+1)
			ops = append(ops, etcd.OpPut(path, addr, etcd.WithLease(lease.ID)))

			for i := 0; i < len(attributes); i += 2 {
				ops = append(ops, etcd.OpPut(path+"/"+attributes[i], attributes[i+1], etcd.WithLease(lease.ID)))
			}

			ctx, cancel = context.WithTimeout(global, 4*time.Second)
			resp, err := s.etcd.Txn(ctx).Then(ops...).Commit()
			cancel()

			if err != nil {
				fmt.Println(err)
				continue
			}

			if !resp.Succeeded {
				continue
			}

			keepAlive, err := s.etcd.KeepAlive(global, lease.ID)
			if err != nil {
				fmt.Println(err)
				continue
			}

			for {
				if _, ok := <-keepAlive; !ok {
					break
				}
			}
		}
	}()

	return s.server.Serve(s.listener)
}

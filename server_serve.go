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
	"github.com/nexcode/rpcplatform/internal/gears"
	etcd "go.etcd.io/etcd/client/v3"
	"time"
)

// Serve starts the gRPC Server and return error if it occurs.
func (s *Server) Serve() error {
	global, cancel := context.WithCancel(context.Background())
	serving := true

	defer func() {
		serving = false
		cancel()
	}()

	go func() {
		path := s.name + "/" + gears.UID()

		for {
			if !serving {
				break
			}

			ctx, cancel := context.WithTimeout(global, 4*time.Second)
			lease, err := s.etcd.Grant(ctx, 4)
			cancel()

			if err != nil {
				fmt.Println(err)
				continue
			}

			ctx, cancel = context.WithTimeout(global, 4*time.Second)
			_, err = s.etcd.Put(ctx, path, s.listener.Addr().String(), etcd.WithLease(lease.ID))
			cancel()

			if err != nil {
				fmt.Println(err)
				continue
			}

			if s.attributes != nil {
				for key, value := range s.attributes.m {
					ctx, cancel = context.WithTimeout(global, 4*time.Second)
					_, err = s.etcd.Put(ctx, path+"/"+key, value, etcd.WithLease(lease.ID))
					cancel()

					if err != nil {
						fmt.Println(err)
					}
				}
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

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
	"log"
	"time"

	"github.com/nexcode/rpcplatform/internal/attributes"
	"github.com/nexcode/rpcplatform/internal/gears"
	etcd "go.etcd.io/etcd/client/v3"
)

// Serve starts the gRPC server and blocks until it exits or an error occurs.
func (s *Server) Serve(ctx context.Context) error {
	path := s.name + "/" + s.id
	attributes := attributes.Values(s.config.Attributes)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		defer func() {
			timer := time.AfterFunc(s.config.StopTimeout, func() {
				s.Server().Stop()
			})

			if s.config.StopTimeout > 0 {
				s.Server().GracefulStop()
				timer.Stop()
			}
		}()

		for {
			if ctx.Err() != nil {
				return
			}

			ctxTimeout, cancelTimeout := gears.ContextTimeout(ctx, s.config.EtcdClientTimeout)
			lease, err := s.etcd.Grant(ctxTimeout, int64(s.config.EtcdLeaseTimeout.Seconds()))
			cancelTimeout()

			if err != nil {
				log.Println(err)
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

			ctxTimeout, cancelTimeout = gears.ContextTimeout(ctx, s.config.EtcdClientTimeout)
			resp, err := s.etcd.Txn(ctxTimeout).Then(ops...).Commit()
			cancelTimeout()

			if err != nil {
				log.Println(err)
				continue
			}

			if !resp.Succeeded {
				continue
			}

			keepAlive, err := s.etcd.KeepAlive(ctx, lease.ID)
			if err != nil {
				log.Println(err)
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

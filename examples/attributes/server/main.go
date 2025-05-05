/*
 * Copyright 2024 RPCPlatform Authors
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

package main

import (
	"context"
	"fmt"

	"github.com/nexcode/rpcplatform"
	"github.com/nexcode/rpcplatform/examples/quickstart/proto"
	etcd "go.etcd.io/etcd/client/v3"
)

type sumServer struct {
	proto.UnimplementedSumServer
}

func (s *sumServer) Sum(_ context.Context, request *proto.SumRequest) (*proto.SumResponse, error) {
	a := request.GetA()
	b := request.GetB()
	sum := a + b

	fmt.Println("request:", a, "+", b)

	return &proto.SumResponse{
		Sum: &sum,
	}, nil
}

func main() {
	etcdClient, err := etcd.New(etcd.Config{
		Endpoints: []string{"localhost:2379"},
	})

	if err != nil {
		panic(err)
	}

	rpcp, err := rpcplatform.New("rpcplatform", etcdClient)
	if err != nil {
		panic(err)
	}

	attributes := rpcplatform.Attributes().Server()
	attributes.SetBalancerWeight(4)

	server, err := rpcp.NewServer("server", "localhost:", attributes)
	if err != nil {
		panic(err)
	}

	proto.RegisterSumServer(server.Server(), &sumServer{})

	if err = server.Serve(); err != nil {
		panic(err)
	}
}

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
	"math/rand"
	"time"

	"github.com/nexcode/rpcplatform"
	"github.com/nexcode/rpcplatform/examples/quickstart/proto"
	"github.com/nexcode/rpcplatform/options"
	etcd "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	etcdClient, err := etcd.New(etcd.Config{
		Endpoints: []string{"localhost:2379"},
	})

	if err != nil {
		panic(err)
	}

	rpcp, err := rpcplatform.New("rpcplatform", etcdClient,
		options.Platform.ClientOptions(
			options.Client.GRPCOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
		),
	)

	if err != nil {
		panic(err)
	}

	client, err := rpcp.NewClient("server", options.Client.MaxActiveServers(2))
	if err != nil {
		panic(err)
	}

	sumClient := proto.NewSumClient(client.Client())

	for {
		time.Sleep(time.Second)

		a := int64(rand.Intn(10))
		b := int64(rand.Intn(10))

		resp, err := sumClient.Sum(context.Background(), &proto.SumRequest{
			A: &a,
			B: &b,
		})

		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(a, "+", b, "=", resp.GetSum())
	}
}

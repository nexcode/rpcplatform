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

package main

import (
	"context"
	"fmt"
	"github.com/nexcode/rpcplatform"
	"github.com/nexcode/rpcplatform/examples/quickstart/proto"
	"github.com/nexcode/rpcplatform/options"
	etcd "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
)

type sumServer struct {
	proto.UnimplementedSumServer
}

func (s *sumServer) Sum(_ context.Context, in *proto.SumRequest) (*proto.SumResponse, error) {
	a := in.GetA()
	b := in.GetB()

	fmt.Println("request:", a, "+", b)

	return &proto.SumResponse{
		Sum: a + b,
	}, nil
}

func main() {
	etcdClient, err := etcd.New(etcd.Config{
		Endpoints: []string{"localhost:2379"},
	})

	if err != nil {
		panic(err)
	}

	jaegerExporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		panic(err)
	}

	zipkinExporter, err := zipkin.New("http://localhost:9411/api/v2/spans")
	if err != nil {
		panic(err)
	}

	rpcp, err := rpcplatform.New("rpcplatform", etcdClient,
		options.OpenTelemetry("server", 1, jaegerExporter, zipkinExporter),
	)

	if err != nil {
		panic(err)
	}

	server, err := rpcp.NewServer("server", "localhost:")
	if err != nil {
		panic(err)
	}

	proto.RegisterSumServer(server.Server(), &sumServer{})

	if err = server.Serve(); err != nil {
		panic(err)
	}
}

# RPCPlatform

[![Build](https://github.com/nexcode/rpcplatform/actions/workflows/build.yml/badge.svg)](https://github.com/nexcode/rpcplatform/actions/workflows/build.yml)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/nexcode/rpcplatform)](https://pkg.go.dev/github.com/nexcode/rpcplatform)
[![GoReportCard](https://goreportcard.com/badge/github.com/nexcode/rpcplatform)](https://goreportcard.com/report/github.com/nexcode/rpcplatform)
[![CodeCov](https://codecov.io/gh/nexcode/rpcplatform/graph/badge.svg)](https://codecov.io/gh/nexcode/rpcplatform)

An **easy-to-use** platform for building microservices without complex infrastructure dependencies.
Only [etcd](https://etcd.io) is required. Out of the box, you get service discovery, distributed tracing, and other useful features.
[gRPC](https://grpc.io) is used for communication between services.

## etcd is required

If you don't have etcd in your infrastructure, you can run it via
[Docker](https://etcd.io/docs/v3.6/op-guide/container/) for testing:

```shell
docker run -d --name etcd \
	-p 2379:2379 -p 2380:2380 \
	-e ETCD_NAME=etcd -e ETCD_INITIAL_CLUSTER=etcd=http://127.0.0.1:2380 \
	-e ETCD_INITIAL_ADVERTISE_PEER_URLS=http://127.0.0.1:2380 -e ETCD_LISTEN_PEER_URLS=http://0.0.0.0:2380 \
	-e ETCD_ADVERTISE_CLIENT_URLS=http://127.0.0.1:2379 -e ETCD_LISTEN_CLIENT_URLS=http://0.0.0.0:2379 \
	gcr.io/etcd-development/etcd:v3.6.5
```

Of course, you can use Docker in production or install etcd using your favorite package manager.
Just remember that the example above is for testing purposes!

## How does it work?

All you need to do is assign a name to your server. When it starts, it automatically selects an available port and listens on it (unless you specify otherwise).
All clients will connect to this server by its name. If there are multiple server instances with the same name, the load is automatically distributed among them.

> The following code examples use a pre-built [proto](examples/quickstart/proto).

First, let's create a new `rpcplatform` instance and a new server named `myServerName`, register the implementation of our `Sum` service, and run it on localhost:

```go
package main

import (
	"context"

	"github.com/nexcode/rpcplatform"
	"github.com/nexcode/rpcplatform/examples/quickstart/proto"
	etcd "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type sumServer struct {
	proto.UnimplementedSumServer
}

func (s *sumServer) Sum(_ context.Context, request *proto.SumRequest) (*proto.SumResponse, error) {
	return proto.SumResponse_builder{
		Sum: new(request.GetA() + request.GetB()),
	}.Build(), nil
}

func main() {
	etcdClient, err := etcd.New(etcd.Config{
		Endpoints: []string{"localhost:2379"},
	})

	if err != nil {
		panic(err)
	}

	rpcp, err := rpcplatform.New("rpcplatform", etcdClient,
		rpcplatform.PlatformOptions.ClientOptions(
			rpcplatform.ClientOptions.GRPCOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
		),
	)

	if err != nil {
		panic(err)
	}

	server, err := rpcp.NewServer("myServerName", "localhost:")
	if err != nil {
		panic(err)
	}

	proto.RegisterSumServer(server.Server(), &sumServer{})

	if err = server.Serve(); err != nil {
		panic(err)
	}
}
```

For the client, we also create a new `rpcplatform` instance and a new client named `myServerName`:

```go
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/nexcode/rpcplatform"
	"github.com/nexcode/rpcplatform/examples/quickstart/proto"
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
		rpcplatform.PlatformOptions.ClientOptions(
			rpcplatform.ClientOptions.GRPCOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
		),
	)

	if err != nil {
		panic(err)
	}

	client, err := rpcp.NewClient("myServerName")
	if err != nil {
		panic(err)
	}

	sumClient := proto.NewSumClient(client.Client())

	for {
		time.Sleep(time.Second)

		sumRequest := proto.SumRequest_builder{
			A: new(int64(rand.Intn(10))),
			B: new(int64(rand.Intn(10))),
		}.Build()

		sumResponse, err := sumClient.Sum(context.Background(), sumRequest)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(sumRequest.GetA(), "+", sumRequest.GetB(), "=", sumResponse.GetSum())
	}
}
```

That's all you need: add or remove server instances dynamically and create clients at any time — `rpcplatform` automatically handles service discovery and load balancing.

### OpenTelemetry

To visualize our **service graph** and get **telemetry for all gRPC methods**, we need to run containers with telemetry services and enable telemetry in `rpcplatform`.

Let's run containers with Zipkin and Jaeger:

```shell
docker run -d --name zipkin -p 9411:9411 openzipkin/zipkin
docker run -d --name jaeger -p 16686:16686 -p 4317:4317 jaegertracing/all-in-one
```

Now let's create the necessary collectors and add the `OpenTelemetry` option to the `rpcplatform` instance:

```go
otlpExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint("localhost:4317"), otlptracegrpc.WithInsecure())
if err != nil {
	panic(err)
}

zipkinExporter, err := zipkin.New("http://localhost:9411/api/v2/spans")
if err != nil {
	panic(err)
}

rpcp, err := rpcplatform.New("rpcplatform", etcdClient,
	rpcplatform.PlatformOptions.OpenTelemetry("myServiceName", 1, otlpExporter, zipkinExporter),
	// other options...
)

if err != nil {
	panic(err)
}
```

The tracing dashboards are available at:

| Zipkin (`http://localhost:9411`)             | Jaeger (`http://localhost:16686`)            |
| :------------------------------------------: | :------------------------------------------: |
| ![Zipkin](examples/opentelemetry/zipkin.png) | ![Jaeger](examples/opentelemetry/jaeger.png) |

## Usage examples

- [QuickStart](examples/quickstart): contains the simplest example without additional features
- [OpenTelemetry](examples/opentelemetry): example integrating distributed tracing systems
- [Attributes](examples/attributes): example using additional settings for client and server

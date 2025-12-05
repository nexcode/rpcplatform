# RPCPlatform

[![Build](https://github.com/nexcode/rpcplatform/actions/workflows/build.yml/badge.svg)](https://github.com/nexcode/rpcplatform/actions/workflows/build.yml)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/nexcode/rpcplatform)](https://pkg.go.dev/github.com/nexcode/rpcplatform)
[![GoReportCard](https://goreportcard.com/badge/github.com/nexcode/rpcplatform)](https://goreportcard.com/report/github.com/nexcode/rpcplatform)
[![CodeCov](https://codecov.io/gh/nexcode/rpcplatform/graph/badge.svg)](https://codecov.io/gh/nexcode/rpcplatform)

An `easy-to-use` platform for creating microservices without complex infrastructure dependencies.
Only [etcd](https://etcd.io) is required. Out of the box you get service discovery, tracing between
services and other useful features. [gRPC](https://grpc.io) is used for communication between services.

## etcd is required

If there is no etcd in your infrastructure, you can run it via
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

All you need to do is assign a name to your server. When it starts, it will automatically select a free port and listen on it (unless you specify otherwise).
All clients will connect to this server by its name. If there are multiple servers with the same name, load balancing is distributed among them.

> In the following code examples, error handling is omitted to improve readability. A pre-built [proto](examples/quickstart/proto) will also be used.

First, let's create a new `rpcplatform` instance:

```go
rpcp, err := rpcplatform.New("rpcplatform", etcdClient,
	rpcplatform.PlatformOptions.ClientOptions(
		rpcplatform.ClientOptions.GRPCOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
	),
)
```

Now let's create a new server named `myServerName`, register the implementation of our `Sum` service and run it on localhost (`sumServer` implementation is omitted):

```go
server, err := rpcp.NewServer("myServerName", "localhost:")
proto.RegisterSumServer(server.Server(), &sumServer{})
err = server.Serve()
```

And finally, we create a client that connects to our `myServerName` (`sumClient` usage is omitted):

```go
client, err := rpcp.NewClient("myServerName")
sumClient := proto.NewSumClient(client.Client())
```

This setup is fully functional: you can add or remove copies of your server and create new clients dynamically.
But to see our **service graph** and get **telemetry for all gRPC methods**, we need to run containers with telemetry services and enable telemetry in `rpcplatform`.

Let's run containers with Zipkin and Jaeger:

```shell
docker run -d --name zipkin -p 9411:9411 openzipkin/zipkin
docker run -d --name jaeger -p 16686:16686 -p 4317:4317 jaegertracing/all-in-one
```

Now let's create the necessary collectors and add the `OpenTelemetry` option to the `rpcplatform` instance:

```go
otlpExporter, err := otlptracegrpc.New(context.Background(),
	otlptracegrpc.WithEndpoint("localhost:4317"), otlptracegrpc.WithInsecure(),
)

zipkinExporter, err := zipkin.New("http://localhost:9411/api/v2/spans")

rpcp, err := rpcplatform.New("rpcplatform", etcdClient,
	rpcplatform.PlatformOptions.ClientOptions(
		rpcplatform.ClientOptions.GRPCOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
	),
	rpcplatform.PlatformOptions.OpenTelemetry("myServiceName", 1, otlpExporter, zipkinExporter),
)
```

The tracing dashboards are available at:

| Zipkin (`http://localhost:9411`)             | Jaeger (`http://localhost:16686`)            |
| :------------------------------------------: | :------------------------------------------: |
| ![Zipkin](examples/opentelemetry/zipkin.png) | ![Jaeger](examples/opentelemetry/jaeger.png) |

## Usage examples (with source code)

- [QuickStart](examples/quickstart): contains the simplest example without additional features
- [OpenTelemetry](examples/opentelemetry): example integrating distributed tracing systems
- [Attributes](examples/attributes): example using additional settings for client and server

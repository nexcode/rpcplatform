# RPCPlatform

[![Build](https://github.com/nexcode/rpcplatform/actions/workflows/build.yml/badge.svg)](https://github.com/nexcode/rpcplatform/actions/workflows/build.yml)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/nexcode/rpcplatform)](https://pkg.go.dev/github.com/nexcode/rpcplatform)
[![GoReportCard](https://goreportcard.com/badge/github.com/nexcode/rpcplatform)](https://goreportcard.com/report/github.com/nexcode/rpcplatform)

An `easy-to-use` platform for creating microservices without complex infrastructure solutions.
Only [etcd](https://etcd.io) required. Out of the box you get a service discovery, tracing between
services and other useful things. [gRPC](https://grpc.io) is used for communication between services.

## etcd required

If there is no etcd in your infrastructure, you can install it from a
[docker container](https://etcd.io/docs/v3.6/op-guide/container/) for tests:

```shell
docker run -d --name etcd \
	-p 2379:2379 -p 2380:2380 \
	gcr.io/etcd-development/etcd:v3.6.5 /usr/local/bin/etcd \
	--name etcd --initial-cluster etcd=http://127.0.0.1:2380 \
	--initial-advertise-peer-urls http://127.0.0.1:2380 --listen-peer-urls http://0.0.0.0:2380 \
	--advertise-client-urls http://127.0.0.1:2379 --listen-client-urls http://0.0.0.0:2379
```

Of course, you can use docker in production or install etcd using your favorite package manager.
Just remember that the example above is for testing purposes!

## How does it work?

All you need to do is give your server a name. When it starts, it will automatically select a free port and run on it (unless you specify otherwise).
All clients will connect to this server by its name. If there are multiple servers with the same name, load balancing will be performed between them.

> In the following code examples, error handling will be removed to improve readability. A pre-built [proto](examples/quickstart/proto) will also be used.

First, let's create a new `rpcplatform` instance:

```go
rpcp, err := rpcplatform.New("rpcplatform", etcdClient,
	rpcplatform.PlatformOptions.ClientOptions(
		rpcplatform.ClientOptions.GRPCOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
	),
)
````

Now let's create a new server named `myServerName`, give it the implementation of our `Sum` service and run it on the localhost (`sumServer` implementation will be omitted):

```go
server, err := rpcp.NewServer("myServerName", "localhost:")
proto.RegisterSumServer(server.Server(), &sumServer{})
err = server.Serve()
````

And finally, we create a client that connects to our `myServerName` (`sumClient` usage will be omitted):

```go
client, err := rpcp.NewClient("myServerName")
sumClient := proto.NewSumClient(client.Client())
````

This is already enough for everything to work, we can add or remove copies of our server and add new clients — everything will work!
But to see our **service graph** and get **telemetry for all gRPC methods**, we need to run containers with telemetry services and enable telemetry in `rpcplatform`.

Let's run containers with Zipkin and Jaeger:

```shell
docker run -d --name zipkin -p 9411:9411 openzipkin/zipkin
docker run -d --name jaeger -p 16686:16686 -p 4317:4317 jaegertracing/all-in-one
```

Now let's create the necessary collectors and add `OpenTelemetry` option to the `rpcplatform` instance:

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
````

The tracing system's web interface is available in a browser:

| Zipkin (`http://localhost:9411`)             | Jaeger (`http://localhost:16686`)            |
| :------------------------------------------: | :------------------------------------------: |
| ![Zipkin](examples/opentelemetry/zipkin.png) | ![Jaeger](examples/opentelemetry/jaeger.png) |

## Usage examples (with source code)

- [QuickStart](examples/quickstart): contains the simplest example without additional features
- [OpenTelemetry](examples/opentelemetry): example with connecting distributed tracing systems
- [Attributes](examples/attributes): example using additional settings for client and server

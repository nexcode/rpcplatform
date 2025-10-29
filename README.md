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
docker run -d \
  -p 2379:2379 \
  -p 2380:2380 \
  --name etcd gcr.io/etcd-development/etcd:v3.6.5 \
  /usr/local/bin/etcd \
  --name etcd \
  --initial-advertise-peer-urls http://127.0.0.1:2380 --listen-peer-urls http://0.0.0.0:2380 \
  --advertise-client-urls http://127.0.0.1:2379 --listen-client-urls http://0.0.0.0:2379 \
  --initial-cluster etcd=http://127.0.0.1:2380
```

Of course, you can use docker in production or install etcd using your favorite package manager.
Just remember that the example above is for testing purposes!

## Usage examples

- [QuickStart](examples/quickstart): contains the simplest example without additional features
- [OpenTelemetry](examples/opentelemetry): example with connecting distributed tracing systems
- [Attributes](examples/attributes): example using additional settings for client and server

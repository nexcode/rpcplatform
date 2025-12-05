## Introduction to gRPC usage

This example contains a pre-compiled [proto](proto) package.
Detailed instructions for rebuilding (if needed) are [attached](proto/README.md).

## Launching this demo

```shell
cd examples/quickstart/server
go run .
```

```shell
cd examples/quickstart/server
go run .
```

```shell
cd examples/quickstart/client
go run .
```

You should now have two servers and one client running!

## Does it work?

You will see the corresponding logs in your console:

```shell
request: 4 + 9
request: 1 + 3
```

```shell
request: 3 + 0
```

```shell
4 + 9 = 13
3 + 0 = 3
1 + 3 = 4
```

You can start additional servers, stop existing ones, and add or remove clients.
The system will continue to function correctly!

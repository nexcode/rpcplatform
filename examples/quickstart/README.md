## Introduction to gRPC usage

This example contains an already compiled [proto](proto) package.
Detailed instructions for rebuilding (if needed) are [attached](proto/README.md).

## Launching this demo

```shell
cd examples/quickstart/server
go run main.go
```

```shell
cd examples/quickstart/server
go run .
```

```shell
cd examples/quickstart/client
go run .
```

Now two servers and one client are running!  

## It works?

Your console will have the corresponding logs:

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

You can start new servers, stop running ones and add or remove clients...
Everything will work!

## gRPC building instructions

The first thing you need is [Protocol Buffer Compiler Installation](https://grpc.io/docs/protoc-installation/).  
Then install go plugins for the protocol compiler:

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Now you can rebuild this with the following command:

```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative sum.proto
```

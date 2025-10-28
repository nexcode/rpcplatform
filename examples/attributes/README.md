## Adding attributes

This example is very similar to [QuickStart](../quickstart), but now we will use additional attributes for server.

For the client, we will set the maximum number of active servers (all other servers above this value will be backup servers):

```go
client, err := rpcp.NewClient("server", options.Client.MaxActiveServers(2))
if err != nil {
	panic(err)
}
```

And for the server we will set the balancing weight attribute:

```go
attributes := attributes.New()
attributes.BalancerWeight = 4

server, err := rpcp.NewServer("server", "localhost:", options.Server.Attributes(attributes))
if err != nil {
	panic(err)
}
```

## Launching this demo

```shell
cd examples/opentelemetry/server
go run .
```

```shell
cd examples/opentelemetry/server
go run .
```

```shell
cd examples/opentelemetry/client
go run .
```

## Additional comments

Currently there are two servers and one client running. If you launch another server, then one of the three servers will become a backup server, and the client will continue to interact with only two servers.

If we want active servers to be selected using priority, we can use the `BalancerPriority` option for server attributes.

Each server receives requests based on its weight, so by changing the value of `BalancerWeight` attribute we will distribute load the way we need.

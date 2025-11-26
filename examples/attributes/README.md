## Adding attributes

This example is very similar to [QuickStart](../quickstart), but now we will use additional attributes for the server and will receive them on every update in the channel.

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

Currently there are two servers and one client running. If you launch another server, then one of the three servers will become a backup server, and the client will continue to interact with only two servers. If we want active servers to be selected using priority, we can use the `BalancerPriority` option for server attributes. Each server receives requests based on its weight, so by changing the value of `BalancerWeight` attribute we will distribute load the way we need.

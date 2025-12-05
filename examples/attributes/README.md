## Adding attributes

This example is similar to the [QuickStart](../quickstart) example, but it uses additional server attributes that are sent with every channel update.

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

## Additional notes

Currently, two servers and one client are running. If another server is launched, one becomes a backup server, and the client continues interacting with only two servers. To select active servers by priority, use the `BalancerPriority` server attribute. Each server receives requests based on its weight, so changing the `BalancerWeight` attribute value distributes load as needed.
